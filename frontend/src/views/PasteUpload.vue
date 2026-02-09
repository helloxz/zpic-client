<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, h, computed } from 'vue'
import { useMessage, NCard, NButton, NIcon, NEmpty, NSpin, NImage, NTooltip, NDropdown, NDrawer, NDrawerContent, NForm, NFormItem, NSelect, NSwitch, useDialog } from 'naive-ui'
import {
  CloudUploadOutline,
  CopyOutline,
  TrashOutline,
  InformationCircleOutline,
  SettingsOutline
} from '@vicons/ionicons5'
import req, { toForm } from '../utils/req'
import { useBaseStore } from '../stores/base'
import { ClipboardSetText,BrowserOpenURL } from "../../wailsjs/runtime/runtime"
import { useI18n } from 'vue-i18n'
const { t, locale } = useI18n()



interface UploadResult {
  imgid: string
  path: string
  url: string
  thumbnail_url: string
  width: number
  height: number
  filename: string
  size: number
  uploaded_at: string
}

interface UploadParams {
  dedup: boolean
  album_id: number
  watermark: boolean
  compress: boolean
}

interface Settings {
  album_id: number
  dedup: boolean
  watermark: boolean
  compress: boolean
  auto_copy: string
}

const message = useMessage()
const dialog = useDialog()
const baseStore = useBaseStore()
const isDragging = ref(false)
const isUploading = ref(false)
const fileInputRef = ref<HTMLInputElement | null>(null)
const uploadHistory = ref<UploadResult[]>([])
const showSettings = ref(false)

const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif', 'image/bmp', 'image/webp']
const ALLOWED_EXTENSIONS = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp']
const MAX_FILE_SIZE = 10 * 1024 * 1024
const MAX_HISTORY = 20
const STORAGE_KEY = 'paste_upload_history'
const SETTINGS_KEY = 'paste_upload_settings'

const defaultSettings: Settings = {
  album_id: 0,
  dedup: true,
  watermark: false,
  compress: false,
  auto_copy: 'url'
}

const settings = ref<Settings>({ ...defaultSettings })

const autoCopyOptions = [
  { label: 'URL', value: 'url' },
  { label: 'Markdown', value: 'markdown' },
  { label: 'HTML', value: 'html' },
  { label: 'BBCode', value: 'bbcode' }
]

const albumOptions = computed(() =>
  baseStore.albumList.map((item) => ({
    label: item.name,
    value: item.id
  }))
)

const loadSettings = () => {
  try {
    const data = localStorage.getItem(SETTINGS_KEY)
    if (data) {
      settings.value = { ...defaultSettings, ...JSON.parse(data) }
    }
  } catch (e) {
    console.error('加载设置失败:', e)
  }
}

const saveSettings = () => {
  try {
    localStorage.setItem(SETTINGS_KEY, JSON.stringify(settings.value))
    message.success('设置已保存')
    showSettings.value = false
  } catch (e) {
    console.error('保存设置失败:', e)
    message.error('保存设置失败')
  }
}

const loadHistory = () => {
  try {
    const data = localStorage.getItem(STORAGE_KEY)
    if (data) {
      uploadHistory.value = JSON.parse(data)
    }
  } catch (e) {
    console.error('加载历史记录失败:', e)
  }
}

const saveHistory = () => {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(uploadHistory.value))
  } catch (e) {
    console.error('保存历史记录失败:', e)
  }
}

const addToHistory = (result: UploadResult) => {
  uploadHistory.value.unshift({
    ...result,
    uploaded_at: new Date().toISOString()
  })
  if (uploadHistory.value.length > MAX_HISTORY) {
    uploadHistory.value = uploadHistory.value.slice(0, MAX_HISTORY)
  }
  saveHistory()
}

const showClearConfirm = () => {
  dialog.warning({
    title: '确认清空',
    content: '此操作不会删除云端图片，仅清空本地上传记录，是否继续？',
    positiveText: '确认清空',
    negativeText: '取消',
    onPositiveClick: () => {
      uploadHistory.value = []
      localStorage.removeItem(STORAGE_KEY)
      message.success('历史记录已清空')
    }
  })
}

const removeItem = (index: number) => {
  uploadHistory.value.splice(index, 1)
  saveHistory()
}

const formatSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
}

const isValidFile = (file: File) => {
  if (file.size > MAX_FILE_SIZE) {
    return false
  }
  if (!ALLOWED_TYPES.includes(file.type)) {
    const ext = '.' + file.name.split('.').pop()?.toLowerCase()
    return ALLOWED_EXTENSIONS.includes(ext)
  }
  return true
}

const getUploadParams = (): UploadParams => {
  return {
    dedup: settings.value.dedup,
    album_id: settings.value.album_id,
    watermark: settings.value.watermark,
    compress: settings.value.compress
  }
}

const uploadFile = async (file: File) => {
  if (!isValidFile(file)) {
    if (file.size > MAX_FILE_SIZE) {
      message.error(`文件大小不能超过 10M，当前文件 ${formatSize(file.size)}`)
    } else {
      message.error('不支持的文件格式，仅支持 jpg/jpeg/png/gif/bmp/webp')
    }
    return
  }

  isUploading.value = true
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      message.error('请先设置 Token')
      return
    }

    const params = getUploadParams()
    const formData = toForm({
      file,
      params: JSON.stringify(params)
    })
    
    const response = await req.post('/api/v3/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })

    if (response.data && response.data.code === 200) {
      addToHistory(response.data.data)
      message.success('上传成功')
      autoCopy(response.data.data.url, response.data.data.filename)
    } else {
      message.error(t(response.data?.msg) || '上传失败')
    }
  } catch (error: any) {
    console.error('上传失败:', error)
    message.error(error.response?.data?.msg || '上传失败，请重试')
  } finally {
    isUploading.value = false
  }
}

const autoCopy = (url: string, filename: string) => {
  const type = settings.value.auto_copy
  let text = ''
  switch (type) {
    case 'url':
      text = url
      break
    case 'html':
      text = `<img title="${filename}" src="${url}" alt="" />`
      break
    case 'markdown':
      text = `![${filename}](${url})`
      break
    case 'bbcode':
      text = `[img]${url}[/img]`
      break
    default:
      return
  }
  ClipboardSetText(text)
}

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    uploadFile(target.files[0])
  }
  target.value = ''
}

const handleDrop = (event: DragEvent) => {
  event.preventDefault()
  isDragging.value = false
  if (event.dataTransfer?.files && event.dataTransfer.files.length > 0) {
    uploadFile(event.dataTransfer.files[0])
  }
}

const handleDragOver = (event: DragEvent) => {
  event.preventDefault()
  isDragging.value = true
}

const handleDragLeave = (event: DragEvent) => {
  event.preventDefault()
  isDragging.value = false
}

const handlePaste = async (event: ClipboardEvent) => {
  const items = event.clipboardData?.items
  if (!items) return

  for (const item of items) {
    if (item.type.startsWith('image/')) {
      const file = item.getAsFile()
      if (file) {
        uploadFile(file)
        return
      }
    }
  }
}

const copyOptions = [
  {
    label: '复制 URL',
    key: 'url'
  },
  {
    label: '复制 HTML',
    key: 'html'
  },
  {
    label: '复制 Markdown',
    key: 'markdown'
  },
  {
    label: '复制 BBCode',
    key: 'bbcode'
  }
]

const handleCopy = (key: string, url: string, thumbnailUrl: string) => {
  let text = ''
  switch (key) {
    case 'url':
      text = url
      break
    case 'html':
      text = `<img src="${url}" alt="" />`
      break
    case 'markdown':
      text = `![${thumbnailUrl}](${url})`
      break
    case 'bbcode':
      text = `[img]${url}[/img]`
      break
  }
  navigator.clipboard.writeText(text).then(() => {
    message.success('已复制到剪贴板')
  }).catch(() => {
    message.error('复制失败')
  })
}

const openFileDialog = () => {
  fileInputRef.value?.click()
}

onMounted(() => {
  baseStore.fetchAlbumList()
  loadSettings()
  loadHistory()
  document.addEventListener('paste', handlePaste)
})

onBeforeUnmount(() => {
  document.removeEventListener('paste', handlePaste)
})
</script>

<template>
  <div class="paste-upload">
    <NCard class="upload-card" :class="{ 'uploading': isUploading }">
      <div
        class="upload-area"
        :class="{ 'dragging': isDragging, 'uploading': isUploading }"
        @click="openFileDialog"
        @drop="handleDrop"
        @dragover="handleDragOver"
        @dragleave="handleDragLeave"
      >
        <input
          ref="fileInputRef"
          type="file"
          accept=".jpg,.jpeg,.png,.gif,.bmp,.webp,image/jpeg,image/png,image/gif,image/bmp,image/webp"
          @change="handleFileSelect"
          class="file-input"
        />
        
        <div v-if="isUploading" class="upload-loading">
          <NSpin size="large" />
          <span class="loading-text">上传中...</span>
        </div>
        
        <div v-else class="upload-placeholder">
          <NIcon size="48" color="#2080f0">
            <CloudUploadOutline />
          </NIcon>
          <p class="upload-title">粘贴上传</p>
          <p class="upload-hint">点击选择文件、拖拽或 Ctrl+V / Command+V 粘贴</p>
          <p class="upload-formats">支持 JPG、JPEG、PNG、GIF、BMP、WebP</p>
        </div>
      </div>
    </NCard>

    <NCard class="history-card">
      <template #header>
        <div class="history-header">
          <div class="history-title-wrapper">
            <span class="history-title">上传记录</span>
            <NTooltip trigger="hover" placement="right">
              <template #trigger>
                <NButton quaternary circle size="small">
                  <template #icon><NIcon><InformationCircleOutline /></NIcon></template>
                </NButton>
              </template>
              仅展示最近20条上传记录，更多记录请前往网页版查看。
            </NTooltip>
          </div>
          <div class="header-actions">
            
            <NButton v-if="uploadHistory.length > 0" quaternary type="error" size="small" @click="showClearConfirm">
              <template #icon><NIcon><TrashOutline /></NIcon></template>
              清空
            </NButton>
            <NButton quaternary type="default" size="small" @click="showSettings = true">
              <template #icon><NIcon><SettingsOutline /></NIcon></template>
              设置
            </NButton>
          </div>
        </div>
      </template>

      <div v-if="uploadHistory.length === 0" class="empty-history">
        <NEmpty description="暂无上传记录" />
      </div>

      <div v-else class="history-list">
        <div
          v-for="(item, index) in uploadHistory"
          :key="index"
          class="history-item"
        >
          <div class="history-thumb">
            <NImage
              :src="item.thumbnail_url"
              :alt="item.filename"
              :width="120"
              :height="90"
              object-fit="cover"
              :preview-disabled="true"
              class="thumb-image"
              @click="BrowserOpenURL(item.url)"
            />
          </div>
          <div class="history-info">
            <p class="history-filename">{{ item.filename }}</p>
            <p class="history-size">{{ formatSize(item.size) }}</p>
            <div class="history-row">
              <span class="history-dimensions">{{ item.width }} × {{ item.height }}</span>
              <div class="btns">
                <NButton quaternary type="error" size="small" @click="removeItem(index)">
                  <template #icon><NIcon><TrashOutline /></NIcon></template>
                </NButton>
                <NDropdown
                  :options="copyOptions"
                  @select="(key) => handleCopy(key, item.url, item.thumbnail_url)"
                  trigger="hover"
                >
                  <NButton quaternary type="primary" size="small">
                    <template #icon><NIcon><CopyOutline /></NIcon></template>
                  </NButton>
                </NDropdown>
              </div>
            </div>
          </div>
        </div>
      </div>
    </NCard>

    <NDrawer v-model:show="showSettings" :width="320" placement="right">
      <NDrawerContent title="上传设置" closable>
        <NForm label-placement="top">
          <NFormItem label="选择相册">
            <NSelect
              v-model:value="settings.album_id"
              :options="albumOptions"
              placeholder="选择相册"
            />
          </NFormItem>
          <NFormItem label="图片去重">
            <NSwitch v-model:value="settings.dedup" />
          </NFormItem>
          <NFormItem label="图片压缩">
            <NSwitch v-model:value="settings.compress" />
          </NFormItem>
          <NFormItem label="文字水印">
            <NSwitch v-model:value="settings.watermark" />
          </NFormItem>
          <NFormItem label="上传成功后自动复制">
            <NSelect
              v-model:value="settings.auto_copy"
              :options="autoCopyOptions"
              placeholder="选择复制格式"
            />
          </NFormItem>
        </NForm>
        <template #footer>
          <NButton type="primary" @click="saveSettings">保存</NButton>
        </template>
      </NDrawerContent>
    </NDrawer>
  </div>
</template>

<style scoped>
.paste-upload {
  padding: 18px;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.upload-card {
  border-radius: 12px;
}

.upload-card.uploading {
  pointer-events: none;
}

.upload-area {
  height: 200px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #fafafa;
  position: relative;
}

.upload-area:hover {
  border-color: #2080f0;
  background: #f0f7ff;
}

.upload-area.dragging {
  border-color: #2080f0;
  background: #e6f4ff;
}

.upload-area.uploading {
  border-style: solid;
  background: rgba(255, 255, 255, 0.9);
}

.file-input {
  display: none;
}

.upload-placeholder {
  text-align: center;
  pointer-events: none;
}

.upload-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 12px 0 8px;
}

.upload-hint {
  font-size: 14px;
  color: #666;
  margin: 0 0 8px;
}

.upload-formats {
  font-size: 12px;
  color: #999;
}

.upload-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.loading-text {
  font-size: 14px;
  color: #2080f0;
}

.history-card {
  flex: 1;
  border-radius: 12px;
  overflow: hidden;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.history-title-wrapper {
  display: flex;
  align-items: center;
  gap: 4px;
}

.history-title {
  font-size: 16px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.empty-history {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
}

.history-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
  max-height: calc(100vh - 385px);
  overflow-y: auto;
  padding: 4px;
  padding-right: 10px;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f9f9f9;
  border-radius: 8px;
  transition: background 0.2s ease;
}

.history-item:hover {
  background: #f0f0f0;
}

.history-thumb {
  flex-shrink: 0;
}

.thumb-image {
  display: flex;
  align-items: center;
  border-radius: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  cursor: pointer;
}

.history-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.history-filename {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.history-size {
  font-size: 12px;
  color: #666;
  margin: 0;
}

.history-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.history-dimensions {
  font-size: 12px;
  color: #999;
}
</style>
