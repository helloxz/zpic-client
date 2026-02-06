<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, h, reactive, computed } from 'vue'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import {
  NCard,
  NButton,
  NIcon,
  NDataTable,
  NTag,
  NTooltip,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NDropdown,
  useMessage,
  NSelect,
  NInputNumber
} from 'naive-ui'
import { useBaseStore } from '../stores/base'
import {
  AddOutline,
  TrashOutline,
  SyncOutline,
  CopyOutline,
  DownloadOutline
} from '@vicons/ionicons5'
import { GetUrlsList, AddUrls, UpdateUrlsStatus, DeleteUrlsByIds, ExportUrlsToCsv } from '../../wailsjs/go/core/AppCore'

interface UrlItem {
  id: number
  origin_url: string
  filename: string
  url: string
  image_width: number
  image_height: number
  created_at: string
  updated_at: string
  status: number
}

const message = useMessage()
const baseStore = useBaseStore()
const selectedAlbumId = ref(0)
const albumOptions = computed(() =>
  baseStore.albumList.map((item) => ({
    label: item.name,
    value: item.id
  }))
)

// 移除单独的 page/pageSize refs，改用 DataTable 内置 pagination（remote）
// const page = ref(1)
// const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)
const checkedRowKeys = ref<number[]>([])

const addUrlsText = ref('')
const submitting = ref(false)
const showAddModal = ref(false)
const showDeleteConfirm = ref(false)
const deleting = ref(false)


const renderEllipsisWithTooltip = (value?: string) => {
  const text = (value ?? '').trim()
  const displayText = text || '-'

  return h(
    NTooltip,
    { trigger: 'hover' },
    {
      trigger: () =>
        h(
          'div',
          {
            style: {
              maxWidth: '120px',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap'
            }
          },
          displayText
        ),
      default: () => displayText
    }
  )
}

function createColumns(): DataTableColumns<UrlItem> {
  return [
    {
      type: 'selection',
      width: 70,
      fixed: 'left'
    },
    {
      title: 'ID',
      key: 'id',
      width: 70,
      fixed: 'left'
    },
    {
      title: '原始URL',
      key: 'origin_url',
      width: 120,
      render(row) {
        return renderEllipsisWithTooltip(row.origin_url)
      }
    },
    {
      title: '文件名',
      key: 'filename',
      width: 120,
      render(row) {
        return renderEllipsisWithTooltip(row.filename)
      }
    },
    {
      title: 'URL',
      key: 'url',
      width: 120,
      render(row) {
        return renderEllipsisWithTooltip(row.url)
      }
    },
    {
      title: '宽',
      key: 'image_width',
      width: 70
    },
    {
      title: '高',
      key: 'image_height',
      width: 70
    },
    {
      title: '添加时间',
      key: 'created_at',
      width: 120,
      render(row) {
        return baseStore.formatDateTimeToSecond(row.created_at)
      }
    },
    {
      title: '更新时间',
      key: 'updated_at',
      width: 120,
      render(row) {
        return baseStore.formatDateTimeToSecond(row.updated_at)
      }
    },
    {
      title: '状态',
      key: 'status',
      width: 70,
      fixed: 'right',
      render(row) {
        const statusMap: Record<number, { type: 'default' | 'success' | 'error' | 'warning'; text: string }> = {
          0: { type: 'default', text: '未开始' },
          1: { type: 'warning', text: '上传中' },
          2: { type: 'success', text: '已完成' },
          3: { type: 'error', text: '上传失败' }
        }
        const config = statusMap[row.status] || { type: 'default', text: '未知' }
        return h(NTag, { type: config.type, size: 'small' }, () => config.text)
      }
    },
    {
      title: '操作',
      key: 'actions',
      width: 60,
      fixed: 'right',
      render(row) {
        return h(NTooltip, { trigger: 'hover' }, {
          trigger: () =>
            h(
              NButton,
              {
                size: 'small',
                quaternary: true,
                onClick: () => copyUrl(row.url)
              },
              { icon: () => h(NIcon, null, () => h(CopyOutline)) }
            ),
          default: () => '复制URL'
        })
      }
    }
  ]
}

const columns = createColumns()

const pagination = reactive<PaginationProps>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (p) => {
    pagination.page = p
    fetchData()
  },
  onUpdatePageSize: (ps) => {
    pagination.pageSize = ps
    pagination.page = 1
    fetchData()
  }
})

const tableData = ref<UrlItem[]>([])

const copyUrl = (url: string) => {
  navigator.clipboard.writeText(url).then(() => {
    message.success('已复制到剪贴板')
  }).catch(() => {
    message.error('复制失败')
  })
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await GetUrlsList({ page: pagination.page as number, limit: pagination.pageSize as number })
    if (res.status && res.data) {
      const data = res.data as { items: UrlItem[]; total: number }
      tableData.value = data.items || []
      total.value = data.total || 0
      pagination.itemCount = total.value
    }
  } catch (err) {
    console.error('获取列表失败:', err)
  } finally {
    loading.value = false
  }
}

const validateUrls = (value: string) => {
  const items = value
    .split('\n')
    .map((item) => item.trim())
    .filter(Boolean)

  if (items.length === 0) {
    return { valid: false, message: '请输入至少一个URL' }
  }

  if (items.length > 100) {
    return { valid: false, message: '最多支持100个URL' }
  }

  const invalidItem = items.find((item) => {
    try {
      const url = new URL(item)
      return !['http:', 'https:'].includes(url.protocol)
    } catch {
      return true
    }
  })

  if (invalidItem) {
    return { valid: false, message: `无效的URL: ${invalidItem}` }
  }

  return { valid: true }
}

const handleAddSubmit = async () => {
  const validation = validateUrls(addUrlsText.value)
  if (!validation.valid) {
    message.warning(validation.message || 'URL格式不正确')
    return
  }

  submitting.value = true
  try {
    const res = await AddUrls({ album_id: selectedAlbumId.value, urls: addUrlsText.value })
    if (res.status) {
      message.success(res.msg || '添加成功')
      showAddModal.value = false
      addUrlsText.value = ''
      fetchData()
    } else {
      message.error(res.msg || '添加失败')
    }
  } catch (err) {
    console.error('添加失败:', err)
    message.error('添加失败，请重试')
  } finally {
    submitting.value = false
  }
}

const showAddDialog = () => {
  addUrlsText.value = ''
  showAddModal.value = true
}

// 导出相关逻辑
const exportLimit = ref(1000)
const exporting = ref(false)

const handleExport = async () => {
  try {
    exporting.value = true
    // 调用后端导出方法，后端会处理 SaveFileDialog
    const res = await ExportUrlsToCsv({ limit: exportLimit.value })

    if (res.status) {
      message.success(res.msg || '导出成功')
    } else {
      // 如果是取消导出，使用 info 提示而不是 error
      if (res.msg === '已取消导出') {
        message.info(res.msg)
      } else {
        message.error(res.msg || '导出失败')
      }
    }
  } catch (err) {
    console.error('导出异常:', err)
    message.error('导出过程发生错误')
  } finally {
    exporting.value = false
  }
}

const changeStatus = () => {
  console.log('change status')
}

const deleteSelected = () => {
  if (!checkedRowKeys.value.length) {
    message.warning('请先选择要删除的行')
    return
  }
  showDeleteConfirm.value = true
}

const handleDeleteConfirm = async () => {
  if (!checkedRowKeys.value.length) {
    showDeleteConfirm.value = false
    return
  }

  deleting.value = true
  try {
    const res = await DeleteUrlsByIds({ ids: checkedRowKeys.value })
    if (res.status) {
      message.success(res.msg || '删除成功')
      checkedRowKeys.value = []
      showDeleteConfirm.value = false
      fetchData()
    } else {
      message.error(res.msg || '删除失败')
    }
  } catch (err) {
    console.error('删除失败:', err)
    message.error('删除失败，请重试')
  } finally {
    deleting.value = false
  }
}

const statusOptions = [
  { label: '待上传', key: 0 },
  { label: '已完成', key: 2 }
]

const updateSelectedStatus = async (status: number) => {
  if (!checkedRowKeys.value.length) {
    message.warning('请先选择要操作的行')
    return
  }

  loading.value = true
  try {
    const res = await UpdateUrlsStatus({ ids: checkedRowKeys.value, status })
    if (res.status) {
      message.success(res.msg || '状态更新成功')
      checkedRowKeys.value = []
      fetchData()
    } else {
      message.error(res.msg || '状态更新失败')
    }
  } catch (err) {
    console.error('更新状态失败:', err)
    message.error('更新状态失败，请重试')
  } finally {
    loading.value = false
  }
}

const refreshTimer = ref<number | null>(null)

onMounted(() => {
  baseStore.fetchAlbumList()
  fetchData()
  if (refreshTimer.value !== null) return
  refreshTimer.value = window.setInterval(() => {
    fetchData()
  }, 20000)
})

onBeforeUnmount(() => {
  if (refreshTimer.value !== null) {
    window.clearInterval(refreshTimer.value)
    refreshTimer.value = null
  }
})
</script>

<template>
  <div class="url-upload">
    <NCard class="content-card">
      <div class="toolbar">
        <div class="toolbar-actions">
          <NButton @click="showAddDialog">
            <template #icon><NIcon><AddOutline /></NIcon></template>
            添加URL
          </NButton>

          <NDropdown trigger="hover" :options="statusOptions" @select="updateSelectedStatus">
            <NButton>
              <template #icon><NIcon><SyncOutline /></NIcon></template>
              更改状态
            </NButton>
          </NDropdown>

          <NButton type="error" @click="deleteSelected">
            <template #icon><NIcon><TrashOutline /></NIcon></template>
            删除选中
          </NButton>

          <!-- 导出区域 -->
          <div class="export-area">
            <span class="export-label">导出最近:</span>
            <NInputNumber 
              v-model:value="exportLimit" 
              :min="1" 
              :max="10000" 
              style="width: 100px" 
              size="medium"
              :show-button="false"
            />
            <span class="export-label">条</span>
            <NButton @click="handleExport" :loading="exporting">
              <template #icon><NIcon><DownloadOutline /></NIcon></template>
              导出表格
            </NButton>
          </div>
        </div>
      </div>

      <NDataTable
        :columns="columns"
        :data="tableData"
        :pagination="pagination"
        :max-height="518"
        :min-height="518"
        :bordered="false"
        :loading="loading"
        :row-key="(row: UrlItem) => row.id"
        v-model:checked-row-keys="checkedRowKeys"
        class="url-table"
        :scroll-x="1800"
        remote
      />
    </NCard>

    <!-- 外部 NPagination 删除（改用表格内置分页） -->
    <!-- ...existing code... -->

    <!-- 新增：本地 Modal，不依赖 useDialog -->
    <NModal v-model:show="showAddModal" preset="card" title="添加URL" :mask-closable="false" style="width: 600px;">
      <NForm>
        <NFormItem label="URL列表">
          <NInput
            v-model:value="addUrlsText"
            type="textarea"
            placeholder="请输入图片URL，每行一个，一次最多100个"
            :rows="10"
          />
        </NFormItem>
      </NForm>

      <template #footer>
        <div class="modal-footer">
          <NSelect
            v-model:value="selectedAlbumId"
            :options="albumOptions"
            placeholder="选择相册"
            style="min-width: 160px;"
          />
          <NButton type="primary" :loading="submitting" @click="handleAddSubmit">提交</NButton>
          <NButton @click="showAddModal = false" :disabled="submitting">取消</NButton>
        </div>
      </template>
    </NModal>

    <NModal
      v-model:show="showDeleteConfirm"
      preset="card"
      title="确认删除"
      :mask-closable="false"
      style="width: 440px;"
    >
      <div>仅删除本地任务，不影响已上传成功的云端数据，是否继续？</div>
      <template #footer>
        <div class="modal-footer">
          <NButton @click="showDeleteConfirm = false" :disabled="deleting">取消</NButton>
          <NButton type="error" :loading="deleting" @click="handleDeleteConfirm">确认删除</NButton>
        </div>
      </template>
    </NModal>
  </div>
</template>

<style scoped>
.url-upload {
  padding: 18px;
  height: 100%;
  display: flex;
  flex-direction: column;
  /* margin-bottom: 24px; */
}

.content-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  border-radius: 12px;
  margin-bottom: 16px;
}

.toolbar {
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 16px;
}

.toolbar-actions {
  display: flex;
  gap: 8px;
  align-items: center; /* 确保垂直居中 */
  flex-wrap: wrap;
}

.export-area {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: 12px;
  border-left: 1px solid #eee;
  padding-left: 12px;
}

.export-label {
  font-size: 14px;
  color: #666;
}

.url-table {
  flex: 1;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
