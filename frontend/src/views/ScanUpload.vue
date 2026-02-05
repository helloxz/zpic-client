<script setup lang="ts">
import { ref } from 'vue'
import { NIcon, NButton, NTag } from 'naive-ui'
import {
  CloudUploadOutline,
  TrashBinOutline
} from '@vicons/ionicons5'

// 文件项接口定义
interface FileItem {
  id: number
  name: string
  size: number
  type: string
  url: string
  status: 'pending' | 'success' | 'error'
}

// 拖拽区域引用
const dragAreaRef = ref<HTMLElement | null>(null)
// 已选择文件列表
const fileList = ref<FileItem[]>([])
// 是否正在拖拽
const isDragging = ref<boolean>(false)
// 是否正在上传
const uploading = ref<boolean>(false)

/**
 * 处理拖拽悬停事件
 * @param {DragEvent} e - 拖拽事件对象
 */
const handleDragOver = (e: DragEvent) => {
  e.preventDefault()
  isDragging.value = true
}

/**
 * 处理拖拽离开事件
 */
const handleDragLeave = () => {
  isDragging.value = false
}

/**
 * 处理拖拽放下事件
 * @param {DragEvent} e - 拖拽事件对象
 */
const handleDrop = (e: DragEvent) => {
  e.preventDefault()
  isDragging.value = false
  const files = e.dataTransfer?.files
  if (files) {
    handleFiles(files)
  }
}

/**
 * 处理文件选择
 * @param {FileList} files - 文件列表
 */
const handleFiles = (files: FileList) => {
  for (let i = 0; i < files.length; i++) {
    const file = files[i]
    // 只处理图片文件
    if (file.type.startsWith('image/')) {
      fileList.value.push({
        id: Date.now() + i,
        name: file.name,
        size: file.size,
        type: file.type,
        url: URL.createObjectURL(file),
        status: 'pending'
      })
    }
  }
}

/**
 * 处理文件输入框变化
 * @param {Event} e - 事件对象
 */
const handleFileSelect = (e: Event) => {
  const target = e.target as HTMLInputElement
  const files = target.files
  if (files) {
    handleFiles(files)
  }
}

/**
 * 上传所有待上传的文件
 */
const uploadFiles = () => {
  uploading.value = true
  // 模拟上传过程
  setTimeout(() => {
    uploading.value = false
    fileList.value.forEach(file => {
      file.status = 'success'
    })
  }, 2000)
}

/**
 * 清空所有已选择的文件
 */
const clearFiles = () => {
  fileList.value = []
}

/**
 * 格式化文件大小
 * @param {number} bytes - 字节大小
 * @returns {string} 格式化后的大小字符串
 */
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 移除单个文件
 * @param {number} id - 文件ID
 */
const removeFile = (id: number) => {
  const index = fileList.value.findIndex(f => f.id === id)
  if (index > -1) {
    fileList.value.splice(index, 1)
  }
}
</script>

<template>
  <div class="scan-upload">
    <!-- 页面标题区域 -->
    <div class="page-header">
      <h1>扫描上传</h1>
      <p class="subtitle">选择或拖拽图片到下方区域进行上传</p>
    </div>

    <!-- 拖拽上传区域 -->
    <div
      ref="dragAreaRef"
      class="drag-area"
      :class="{ 'drag-active': isDragging }"
      @dragover="handleDragOver"
      @dragleave="handleDragLeave"
      @drop="handleDrop"
    >
      <div class="drag-content">
        <div class="upload-icon">
          <n-icon :size="64" color="#2080f0">
            <CloudUploadOutline />
          </n-icon>
        </div>
        <p class="drag-text">拖拽图片到此处，或</p>
        <n-button type="primary" @click="($refs.fileInput as HTMLInputElement).click()">
          选择文件
        </n-button>
        <!-- 隐藏的文件输入框 -->
        <input
          ref="fileInput"
          type="file"
          multiple
          accept="image/*"
          style="display: none"
          @change="handleFileSelect"
        />
        <p class="drag-hint">支持 PNG、JPG、GIF、WebP 格式</p>
      </div>
    </div>

    <!-- 文件列表区域 -->
    <div class="file-list" v-if="fileList.length > 0">
      <div class="list-header">
        <span class="list-title">已选择 {{ fileList.length }} 个文件</span>
        <div class="list-actions">
          <n-button text @click="clearFiles">清空</n-button>
          <n-button type="primary" @click="uploadFiles" :loading="uploading">
            上传全部
          </n-button>
        </div>
      </div>

      <div class="file-items">
        <div
          v-for="file in fileList"
          :key="file.id"
          class="file-item"
        >
          <!-- 文件预览图 -->
          <div class="file-preview">
            <img :src="file.url" :alt="file.name" />
          </div>
          <!-- 文件信息 -->
          <div class="file-info">
            <span class="file-name">{{ file.name }}</span>
            <span class="file-size">{{ formatFileSize(file.size) }}</span>
          </div>
          <!-- 上传状态标签 -->
          <div class="file-status">
            <n-tag v-if="file.status === 'pending'" type="default">等待上传</n-tag>
            <n-tag v-else-if="file.status === 'success'" type="success">上传成功</n-tag>
            <n-tag v-else-if="file.status === 'error'" type="error">上传失败</n-tag>
          </div>
          <!-- 操作按钮 -->
          <div class="file-actions">
            <n-button text type="error" @click="removeFile(file.id)">
              <template #icon>
                <n-icon>
                  <TrashBinOutline />
                </n-icon>
              </template>
            </n-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 扫描上传页面容器 */
.scan-upload {
  padding: 32px;
  max-width: 900px;
  margin: 0 auto;
}

/* 页面标题区域 */
.page-header {
  margin-bottom: 32px;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 8px 0;
}

.subtitle {
  color: #666;
  font-size: 14px;
  margin: 0;
}

/* 拖拽区域样式 */
.drag-area {
  border: 2px dashed #d9d9d9;
  border-radius: 12px;
  padding: 48px;
  text-align: center;
  transition: all 0.3s ease;
  background: #fafafa;
}

.drag-area.drag-active {
  border-color: #2080f0;
  background: rgba(32, 128, 240, 0.05);
}

.drag-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.upload-icon {
  margin-bottom: 8px;
}

.drag-text {
  color: #666;
  font-size: 16px;
  margin: 0;
}

.drag-hint {
  color: #999;
  font-size: 12px;
  margin: 0;
}

/* 文件列表区域 */
.file-list {
  margin-top: 32px;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.list-title {
  font-weight: 500;
  color: #333;
}

.list-actions {
  display: flex;
  gap: 12px;
}

.file-items {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* 单个文件项样式 */
.file-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  transition: box-shadow 0.2s ease;
}

.file-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.file-preview {
  width: 56px;
  height: 56px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.file-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.file-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.file-name {
  font-weight: 500;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-size {
  font-size: 12px;
  color: #999;
}

.file-status {
  flex-shrink: 0;
}

.file-actions {
  flex-shrink: 0;
}
</style>
