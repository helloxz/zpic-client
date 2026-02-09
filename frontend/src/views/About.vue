<script setup lang="ts">
import { ref, h,onMounted } from 'vue'
import { NCard, NIcon, NTag, useNotification, NDialogProvider, useDialog } from 'naive-ui'
import type { Component } from 'vue'
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime"
import { GetAppInfo } from '../../wailsjs/go/core/AppCore'
import axios from 'axios'
import {
  CloudUploadOutline,
  LinkOutline,
  SettingsOutline,
  LockClosedOutline,
  CheckmarkCircleOutline,
  LogoGithub,
  OpenOutline,
  GlobeOutline,
  RefreshOutline
} from '@vicons/ionicons5'

const notification = useNotification()
const dialog = useDialog()


const appInfos = ref({
  version:'',
  os:''
})

// 更新信息接口定义
interface UpdateInfo {
  latest: string
  date: string
  downloadUrl: string
}

const getAppInfos = async () => {
  const result = await GetAppInfo()
  if(result.status){
    appInfos.value = result.data
    // console.log("获取应用信息成功:", appInfos.value)
  }
}

// 更新信息
const updateInfo = ref<UpdateInfo>({
  latest: '1.0.0',
  date: '2026-02-07',
  downloadUrl: '#'
})

// 功能特性数据接口定义
interface Feature {
  icon: Component
  title: string
  description: string
}

// 功能特性数据
const features: Feature[] = [
  {
    icon: CloudUploadOutline,
    title: '扫描上传',
    description: '扫描指定文件夹，一次性上传数千张图片。'
  },
  {
    icon: LinkOutline,
    title: '链接上传',
    description: '支持粘贴图片链接，批量帮您转换新的图片链接。'
  },
  {
    icon: SettingsOutline,
    title: '灵活配置',
    description: '可自定义上传参数、上传动作。'
  },
  {
    icon: LockClosedOutline,
    title: '本地存储',
    description: 'API Token 仅存储在本地，保障账户安全'
  }
]

// 相关链接数据接口定义
interface Link {
  name: string
  url: string
  icon: Component
}

// 相关链接数据
const links: Link[] = [
  {
    name: '官网',
    url: 'https://dwz.ovh/ovbp',
    icon: GlobeOutline
  },
  {
    name: 'GitHub',
    url: 'https://github.com/helloxz/zpic-client',
    icon: LogoGithub
  },
  {
    name: '问题反馈',
    url: 'https://github.com/helloxz/zpic-client/issues',
    icon: OpenOutline
  },
  {
    name: '使用文档',
    url: 'https://dwz.ovh/myze',
    icon: OpenOutline
  }
]

const checkingUpdate = ref(false)

const compareVersions = (v1: string, v2: string): number => {
  const parts1 = v1.split('.').map(Number)
  const parts2 = v2.split('.').map(Number)
  for (let i = 0; i < Math.max(parts1.length, parts2.length); i++) {
    const p1 = parts1[i] || 0
    const p2 = parts2[i] || 0
    if (p1 !== p2) {
      return p1 - p2
    }
  }
  return 0
}

const checkUpdate = async () => {
  if (checkingUpdate.value) return
  checkingUpdate.value = true

  try {
    const formData = new URLSearchParams()
    formData.append('path', '/UniBin/zpic-client')
    formData.append('type', 'public')

    const response = await axios.post('https://soft.xiaoz.org/api/filelist', formData, {
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
    })

    if (response.data.code === 200 && response.data.data && response.data.data.length > 0) {
      const versions = response.data.data.map((item: any) => item.name)
      const remoteMaxVersion = versions.reduce((max: string, curr: string) => {
        return compareVersions(curr, max) > 0 ? curr : max
      }, '0.0.0')

      const localVersion = appInfos.value.version

      if (compareVersions(remoteMaxVersion, localVersion) <= 0) {
        notification.success({
          content: '当前版本已是最新',
          duration: 3000
        })
      } else {
        dialog.info({
          title: '发现新版本',
          content: `检测到 ${remoteMaxVersion} 版本`,
          positiveText: '前往下载',
          negativeText: '取消',
          onPositiveClick: () => {
            BrowserOpenURL(`http://soft.xiaoz.org/#/UniBin/zpic-client/${remoteMaxVersion}`)
          }
        })
      }
    }
  } catch (error) {
    notification.error({
      content: '检查更新失败，请稍后重试',
      duration: 3000
    })
  } finally {
    checkingUpdate.value = false
  }
}

onMounted(() => {
  getAppInfos()
})
</script>

<template>
  <div class="about">
    <!-- 页面标题区域 -->
    <div class="page-header">
      <h1>关于 Zpic Client</h1>
      <p class="subtitle">一款简洁高效的图片上传客户端，可用于 <a @click="BrowserOpenURL('https://www.imgurl.org/')" href="javascript:;">ImgURL</a> 和 <a @click="BrowserOpenURL('https://go.piclink.cc/')" href="javascript:;">图链</a> ，适合大批量图片上传需求。</p>
    </div>

    <!-- 应用信息卡片 -->
    <n-card class="about-card">
      <div class="app-info">
        <!-- 应用Logo -->
        <div class="app-logo">
          <img src="../assets/images/appicon.png" alt="ZPIC Logo" />
        </div>
        <!-- 应用详细信息 -->
        <div class="app-details">
          <h2 class="app-name">Zpic 图床客户端</h2>
          <p class="app-version">
            版本 <n-tag size="small" type="info">{{ appInfos.version }}</n-tag>
            <n-button size="small" quaternary circle :loading="checkingUpdate" @click="checkUpdate">
              <template #icon>
                <n-icon><RefreshOutline /></n-icon>
              </template>
            </n-button>
          </p>
          <p class="app-desc">让图片上传更简单、更高效</p>
        </div>
      </div>
    </n-card>

    <!-- 功能特性区域 -->
    <div class="features-section">
      <h3 class="section-title">功能特性</h3>
      <div class="features-grid">
        <div
          v-for="feature in features"
          :key="feature.title"
          class="feature-item"
        >
          <div class="feature-icon">
            <n-icon :size="32" color="#2080f0">
              <component :is="feature.icon" />
            </n-icon>
          </div>
          <h4 class="feature-title">{{ feature.title }}</h4>
          <p class="feature-desc">{{ feature.description }}</p>
        </div>
      </div>
    </div>

    <!-- 更新信息卡片 -->
    <!-- <n-card class="update-card">
      <div class="update-info">
        <div class="update-status">
          <n-icon :size="20" color="#52c41a">
            <CheckmarkCircleOutline />
          </n-icon>
          <span>当前版本已是最新</span>
        </div>
        <p class="update-date">发布日期：{{ updateInfo.date }}</p>
      </div>
    </n-card> -->

    <!-- 相关链接区域 -->
    <div class="links-section">
      <h3 class="section-title">相关链接</h3>
      <div class="links-list">
        <a
          v-for="link in links"
          :key="link.name"
          @click="BrowserOpenURL(link.url)"
          class="link-item"
        >
          <n-icon :size="18">
            <component :is="link.icon" />
          </n-icon>
          <span>{{ link.name }}</span>
        </a>
      </div>
    </div>

    
  </div>
</template>

<style scoped>
/* 关于页面容器 */
.about {
  padding: 18px;
  max-width: 900px;
  margin: 0 auto;
}

/* 页面标题区域 */
.page-header {
  margin-bottom: 18px;
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
/* .subtitle a{
  color: #4098fc;
  text-decoration: underline;
  font-weight: 600;
  cursor: pointer;
}
.subtitle a:hover{
  color: #1d77e5;
} */

/* 应用信息卡片 */
.about-card {
  margin-bottom: 32px;
  border-radius: 12px;
}

.app-info {
  display: flex;
  align-items: center;
  gap: 24px;
}

.app-logo img {
  width: 88px;
  height: 88px;
  border-radius: 16px;
  /* box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1); */
}

.app-details {
  flex: 1;
}

.app-name {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 8px 0;
}

.app-version {
  color: #666;
  font-size: 14px;
  margin: 0 0 8px 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.app-desc {
  color: #999;
  font-size: 14px;
  margin: 0;
}

/* 区域标题 */
.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0 0 16px 0;
}

/* 功能特性区域 */
.features-section {
  margin-bottom: 32px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

/* 功能项样式 */
.feature-item {
  padding: 24px;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 12px;
  text-align: center;
  transition: all 0.3s ease;
}

.feature-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.feature-icon {
  margin-bottom: 16px;
}

.feature-title {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin: 0 0 8px 0;
}

.feature-desc {
  font-size: 13px;
  color: #999;
  margin: 0;
  line-height: 1.6;
}

/* 更新信息卡片 */
.update-card {
  margin-bottom: 32px;
  border-radius: 12px;
}

.update-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.update-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  color: #52c41a;
}

.update-date {
  color: #999;
  font-size: 13px;
  margin: 0;
}

/* 相关链接区域 */
.links-section {
  margin-bottom: 32px;
}

.links-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.links-list a{
  cursor: pointer;
}

/* 链接项样式 */
.link-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: #f5f5f5;
  border-radius: 8px;
  color: #333;
  text-decoration: none;
  transition: all 0.2s ease;
}

.link-item:hover {
  background: #e8e8e8;
  color: #2080f0;
}

/* 版权信息 */
.copyright {
  text-align: center;
  padding-top: 24px;
  border-top: 1px solid #f0f0f0;
}

.copyright p {
  margin: 0;
  color: #999;
  font-size: 13px;
}

.tech-stack {
  margin-top: 8px !important;
  color: #d9d9d9 !important;
}
</style>
