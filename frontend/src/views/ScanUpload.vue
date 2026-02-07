<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, h, reactive } from 'vue'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import {
  NCard,
  NButton,
  NIcon,
  NDataTable,
  NTag,
  NTooltip,
  NPopconfirm,
  useMessage
} from 'naive-ui'
import {
  ScanOutline,
  TrashBinOutline,
  RefreshOutline
} from '@vicons/ionicons5'
import { GetScanList, AddScanTask, DeleteTasks, SelectScanDirectory } from '../../wailsjs/go/core/AppCore'

interface TaskItem {
  id: number
  path: string
  success_num: number
  failed_num: number
  total_num: number
  status: number
  created_at: string
  updated_at: string
}

const message = useMessage()
const loading = ref(false)
const checkedRowKeys = ref<number[]>([])
const refreshTimer = ref<number | null>(null)

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

const tableData = ref<TaskItem[]>([])

const formatDateTimeToSecond = (value: string) => {
  const text = (value || '').trim()
  if (!text) return '-'

  if (!text.includes('T') && text.includes(' ')) return text

  const d = new Date(text)
  if (Number.isNaN(d.getTime())) return text

  const pad2 = (n: number) => String(n).padStart(2, '0')
  const yyyy = d.getFullYear()
  const mm = pad2(d.getMonth() + 1)
  const dd = pad2(d.getDate())
  const hh = pad2(d.getHours())
  const mi = pad2(d.getMinutes())
  const ss = pad2(d.getSeconds())
  return `${yyyy}-${mm}-${dd} ${hh}:${mi}:${ss}`
}

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
              maxWidth: '200px',
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

function createColumns(): DataTableColumns<TaskItem> {
  return [
    {
      type: 'selection',
      width: 70,
      fixed: 'left'
    },
    {
      title: 'ID',
      key: 'id',
      width: 80,
      fixed: 'left'
    },
    {
      title: '扫描路径',
      key: 'path',
      width: 200,
      render(row) {
        return renderEllipsisWithTooltip(row.path)
      }
    },
    {
      title: '成功',
      key: 'success_num',
      width: 80
    },
    {
      title: '失败',
      key: 'failed_num',
      width: 80
    },
    {
      title: '总数',
      key: 'total_num',
      width: 80
    },
    {
      title: '创建时间',
      key: 'created_at',
      width: 120,
      render(row) {
        return formatDateTimeToSecond(row.created_at)
      }
    },
    {
      title: '更新时间',
      key: 'updated_at',
      width: 120,
      render(row) {
        return formatDateTimeToSecond(row.updated_at)
      }
    },
    {
      title: '状态',
      key: 'status',
      width: 80,
      fixed: 'right',
      render(row) {
        const statusMap: Record<number, { type: 'default' | 'success' | 'error' | 'warning'; text: string }> = {
          0: { type: 'default', text: '待扫描' },
          1: { type: 'warning', text: '等待上传' },
          2: { type: 'warning', text: '上传中' },
          3: { type: 'success', text: '上传完成' }
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
                type: 'warning',
                onClick: () => handleRetryFailed(row)
              },
              { icon: () => h(NIcon, null, () => h(RefreshOutline)) }
            ),
          default: () => '失败重试'
        })
      }
    }
  ]
}

const columns = createColumns()

const fetchData = async () => {
  try {
    const res = await GetScanList({ page: pagination.page as number, limit: pagination.pageSize as number })
    if (res.status && res.data) {
      const data = res.data as { items: TaskItem[]; total: number }
      tableData.value = data.items || []
      pagination.itemCount = data.total || 0
    }
  } catch (err) {
    console.error('获取任务列表失败:', err)
  }
}

const handleScanUpload = async () => {
  try {
    const path = await SelectScanDirectory()
    if (path) {
      const res = await AddScanTask({ path })
      if (res.status) {
        message.success('任务创建成功')
        fetchData()
      } else {
        message.error(res.msg || '创建失败')
      }
    }
  } catch (err) {
    console.error('选择目录失败:', err)
    message.error('选择目录失败')
  }
}

const handleDeleteConfirm = async () => {
  if (!checkedRowKeys.value.length) return

  loading.value = true
  try {
    const res = await DeleteTasks({ ids: checkedRowKeys.value })
    if (res.status) {
      const data = res.data as { deleted_count: number; skipped_count: number; message: string }
      message.success(`删除成功：${data?.deleted_count || 0}个任务已删除`)
      checkedRowKeys.value = []
      fetchData()
    } else {
      message.error(res.msg || '删除失败')
    }
  } catch (err) {
    console.error('删除失败:', err)
    message.error(err instanceof Error ? err.message : '删除失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

const handleRetryFailed = (row: TaskItem) => {
  console.log('失败重试', row.id)
}

onMounted(() => {
  fetchData()
  if (refreshTimer.value !== null) return
  refreshTimer.value = window.setInterval(() => {
    fetchData()
  }, 30000)
})

onBeforeUnmount(() => {
  if (refreshTimer.value !== null) {
    window.clearInterval(refreshTimer.value)
    refreshTimer.value = null
  }
})
</script>

<template>
  <div class="scan-upload">
    <NCard class="content-card">
      <div class="toolbar">
        <div class="toolbar-actions">
          <NButton type="primary" @click="handleScanUpload">
            <template #icon><NIcon><ScanOutline /></NIcon></template>
            扫描上传
          </NButton>

          <NPopconfirm
            @positive-click="handleDeleteConfirm"
            positive-text="确认删除"
            negative-text="取消"
          >
            <template #trigger>
              <NButton
                type="error"
                :disabled="!checkedRowKeys.length"
                :loading="loading"
              >
                <template #icon><NIcon><TrashBinOutline /></NIcon></template>
                删除选中
              </NButton>
            </template>
            该操作将删除本地上传任务，不影响已上传的云端数据。
          </NPopconfirm>
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
        :row-key="(row: TaskItem) => row.id"
        v-model:checked-row-keys="checkedRowKeys"
        class="scan-table"
        :scroll-x="1400"
        remote
      />
    </NCard>
  </div>
</template>

<style scoped>
.scan-upload {
  padding: 18px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.content-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  border-radius: 12px;
}

.toolbar {
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 16px;
}

.toolbar-actions {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
}

.scan-table {
  flex: 1;
}
</style>
