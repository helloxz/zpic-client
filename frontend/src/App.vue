<script setup lang="ts">
import { ref,onMounted } from 'vue'
import SideNav from './components/SideNav.vue'
import { useBaseStore } from './stores/base';

const baseStore = useBaseStore();
// 控制侧边栏是否折叠
const sideNavCollapsed = ref<boolean>(false)

onMounted(() => {
  baseStore.fetchAlbumList();
});
</script>

<template>
  <div class="app-container">
    <!-- 左侧导航栏 -->
    <SideNav
      class="side-nav"
      :class="{ 'side-nav--collapsed': sideNavCollapsed }"
    />

    <!-- 右侧主内容区域 -->
    <main class="main-content" :class="{ 'main-content--expanded': sideNavCollapsed }">
      <!-- 路由视图：根据URL显示对应的页面组件 -->
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <n-message-provider>
            <n-dialog-provider>
              <component :is="Component" />
            </n-dialog-provider>
          </n-message-provider>
        </transition>
      </router-view>
    </main>
  </div>
</template>

<style>
/* 全局样式重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol',
    'Noto Color Emoji';
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: #f5f7fa;
  color: #333;
  line-height: 1.5;
}

#app {
  height: 100vh;
  overflow: hidden;
}
</style>

<style scoped>
/* 应用容器 - Flex两栏布局 */
.app-container {
  display: flex;
  height: 100vh;
  width: 100%;
}

/* 左侧导航栏样式 */
.side-nav {
  flex-shrink: 0;
  transition: width 0.3s ease;
}

/* 右侧主内容区域 */
.main-content {
  flex: 1;
  min-width: 0;
  background: #f5f7fa;
  overflow-y: auto;
  transition: margin-left 0.3s ease;
}

/* 页面过渡动画 - 淡入淡出 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
