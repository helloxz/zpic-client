<script setup lang="ts">
import { h, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NIcon } from 'naive-ui'
import {
  ScanOutline,
  LinkOutline,
  SettingsOutline,
  InformationCircleOutline,
  ChevronBackOutline,
  ChevronForwardOutline,
  DocumentTextOutline
} from '@vicons/ionicons5'

// 路由 hooks
const router = useRouter()
const route = useRoute()

// 控制菜单折叠状态
const collapsed = ref<boolean>(false)

// 渲染图标组件
const renderIcon = (icon: any) => {
  return () => h(NIcon, null, { default: () => h(icon) })
}

// 导航菜单选项
const menuOptions = [
  {
    label: '扫描上传',
    key: 'scan-upload',
    icon: renderIcon(ScanOutline)
  },
  {
    label: 'URL上传',
    key: 'url-upload',
    icon: renderIcon(LinkOutline)
  },
  {
    label: '设置',
    key: 'settings',
    icon: renderIcon(SettingsOutline)
  },
  {
    label: '日志',
    key: 'logs',
    icon: renderIcon(DocumentTextOutline)
  },
  {
    label: '关于',
    key: 'about',
    icon: renderIcon(InformationCircleOutline)
  }
]

// 处理菜单点击
function handleMenuUpdate(value: string) {
  router.push(`/${value}`)
}
</script>

<template>
  <div class="side-nav" :class="{ 'side-nav--collapsed': collapsed }">
    <!-- 应用Logo和标题区域 -->
    <div class="nav-header">
      <div class="logo">
        <img src="../assets/images/image_new.png" alt="Logo" />
      </div>
      <transition name="fade">
        <h2 v-if="!collapsed" class="app-title">图床客户端</h2>
      </transition>
    </div>

    <!-- 导航菜单 -->
    <n-menu
      :options="menuOptions"
      :value="route.name"
      @update:value="handleMenuUpdate"
      :collapsed-width="64"
      :collapsed-icon-size="22"
      :collapsed="collapsed"
    />

    <!-- 折叠按钮 -->
    <div class="collapse-btn" @click="collapsed = !collapsed">
      <n-icon size="20">
        <ChevronForwardOutline v-if="collapsed" />
        <ChevronBackOutline v-else />
      </n-icon>
    </div>
  </div>
</template>

<style scoped>
.side-nav {
  width: 240px;
  height: 100vh;
  background: #ffffff;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #e8e8e8;
  transition: width 0.3s ease;
  position: relative;
}

.side-nav--collapsed {
  width: 64px;
}

.nav-header {
  padding: 20px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid #e8e8e8;
  min-height: 72px;
}

.logo{
  display: flex;
  align-items: center;
  justify-content: center;
}
.logo img {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  flex-shrink: 0;
}

.app-title {
  color: #1a1a1a;
  font-size: 16px;
  font-weight: 600;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
}

.collapse-btn {
  position: absolute;
  bottom: 16px;
  right: -12px;
  width: 24px;
  height: 24px;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #666;
  transition: all 0.2s ease;
  z-index: 10;
}

.collapse-btn:hover {
  background: #f5f5f5;
  color: #2080f0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
