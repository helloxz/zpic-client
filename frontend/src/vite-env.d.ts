// Shim for Vue Single File Components
// 此文件告诉 TypeScript 如何处理 .vue 文件

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}
