<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NCard, NLog, NSpin, NEmpty, NButton, NIcon } from 'naive-ui'
import { GetRecentLogs, OpenLogDirectory } from '../../wailsjs/go/core/AppCore'
import { useMessage } from 'naive-ui'
import { FolderOpenOutline } from '@vicons/ionicons5'

const message = useMessage()

// 日志数据
const logs = ref([])
// 加载状态
const loading = ref(true)
// 日志选项配置
const logOptions = {
  showIcon: true,
  lineHeight: '24px',
  fontSize: '13px'
}

/**
 * 获取日志数据
 * 调用后端 GetRecentLogs 方法获取最近的日志信息
 */
const fetchLogs = async () => {
  loading.value = true
  const result = await GetRecentLogs()
  if(result.status){
    let resLogs = result.data
    logs.value = resLogs
    // console.log("获取日志成功:", resLogs)
  }
}

function handleReachBottom() {
  message.info('Reach Bottom')
}

const handleOpenLogDirectory = async () => {
  await OpenLogDirectory()
}

onMounted(() => {
  fetchLogs()
})
</script>

<template>
  <div class="log-page">
    <!-- 页面标题区域 -->
    <div class="page-header">
      <div class="header-row">
        <div class="header-text">
          <h1>日志</h1>
          <p class="subtitle">查看应用程序错误日志</p>
        </div>
        <n-button quaternary @click="handleOpenLogDirectory">
          <template #icon>
            <n-icon :component="FolderOpenOutline" />
          </template>
          打开日志目录
        </n-button>
      </div>
    </div>

    <!-- 日志卡片 - 高度占满整个窗口 -->
    <n-card class="log-card">

      <!-- 日志内容 - 使用 NLog 组件展示 -->
      <div class="log-container">
        <n-log
          :lines="logs"
          :rows="33"
          scrollTo="{position:bottom}"
        />
        <!-- 空状态提示 -->
      </div>
    </n-card>
  </div>
</template>

<style scoped>
/* 日志页面容器 */
.log-page {
  padding: 18px;
  /* height:600px !important; */
  display: flex;
  flex-direction: column;
  /*超出高度显示滚动条*/
  /* overflow: auto; */
}

/* 页面标题区域 */
.page-header {
  margin-bottom: 18px;
  flex-shrink: 0;
}

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
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

/* 日志卡片 - 高度占满剩余空间 */
.log-card {
  /*超出高度显示滚动条*/
  /* overflow: auto; */
  flex: 1;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  /* height:500px; */
  /* overflow: hidden; */
}

/* 日志容器 - 占据卡片内所有可用空间 */
.log-container {
  flex: 1;
  /* overflow: hidden; */
  min-height: 0;
}

/* 覆盖 NLog 组件的默认样式 */
:deep(.n-log) {
  height: 100%;
}

/* 确保 NLog 内部容器可以滚动 */
:deep(.n-log__content) {
  height: 100%;
  overflow-y: auto;
}
</style>
