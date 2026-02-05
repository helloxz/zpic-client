<script setup lang="ts">
import { ref, h } from 'vue'
import { NCard, NIcon, NTag } from 'naive-ui'
import type { Component } from 'vue'
import {
  CloudUploadOutline,
  LinkOutline,
  SettingsOutline,
  LockClosedOutline,
  CheckmarkCircleOutline,
  LogoGithub,
  OpenOutline
} from '@vicons/ionicons5'

// 应用版本号
const version = ref<string>('1.0.0')

// 更新信息接口定义
interface UpdateInfo {
  latest: string
  date: string
  downloadUrl: string
}

// 更新信息
const updateInfo = ref<UpdateInfo>({
  latest: '1.0.0',
  date: '2024-01-01',
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
    title: '多种上传方式',
    description: '支持本地文件拖拽、URL导入、剪贴板粘贴等多种上传方式'
  },
  {
    icon: LinkOutline,
    title: '多图床支持',
    description: '兼容 SM.MS、ImgBB、Imgur 等主流图床服务'
  },
  {
    icon: SettingsOutline,
    title: '灵活配置',
    description: '可自定义上传参数、快捷键、文件命名规则等'
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
    name: 'GitHub',
    url: 'https://github.com/example/zpic-client',
    icon: LogoGithub
  },
  {
    name: '问题反馈',
    url: 'https://github.com/example/zpic-client/issues',
    icon: OpenOutline
  },
  {
    name: '使用文档',
    url: 'https://github.com/example/zpic-client/wiki',
    icon: OpenOutline
  }
]
</script>

<template>
  <div class="about">
    <!-- 页面标题区域 -->
    <div class="page-header">
      <h1>关于 ZPIC</h1>
      <p class="subtitle">一款简洁高效的图床上传客户端</p>
    </div>

    <!-- 应用信息卡片 -->
    <n-card class="about-card">
      <div class="app-info">
        <!-- 应用Logo -->
        <div class="app-logo">
          <img src="../assets/images/logo-universal.png" alt="ZPIC Logo" />
        </div>
        <!-- 应用详细信息 -->
        <div class="app-details">
          <h2 class="app-name">ZPIC 图床客户端</h2>
          <p class="app-version">
            版本 {{ version }}
            <n-tag v-if="updateInfo.latest === version" type="success" size="small">
              <template #icon>
                <n-icon>
                  <CheckmarkCircleOutline />
                </n-icon>
              </template>
              最新
            </n-tag>
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
    <n-card class="update-card">
      <div class="update-info">
        <div class="update-status">
          <n-icon :size="20" color="#52c41a">
            <CheckmarkCircleOutline />
          </n-icon>
          <span>当前版本已是最新</span>
        </div>
        <p class="update-date">发布日期：{{ updateInfo.date }}</p>
      </div>
    </n-card>

    <!-- 相关链接区域 -->
    <div class="links-section">
      <h3 class="section-title">相关链接</h3>
      <div class="links-list">
        <a
          v-for="link in links"
          :key="link.name"
          :href="link.url"
          target="_blank"
          class="link-item"
        >
          <n-icon :size="18">
            <component :is="link.icon" />
          </n-icon>
          <span>{{ link.name }}</span>
        </a>
      </div>
    </div>

    <!-- 版权信息 -->
    <div class="copyright">
      <p>© 2024 ZPIC Client. All rights reserved.</p>
      <p class="tech-stack">Built with Wails + Vue 3 + Naive UI</p>
    </div>
  </div>
</template>

<style scoped>
/* 关于页面容器 */
.about {
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
  width: 80px;
  height: 80px;
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
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
