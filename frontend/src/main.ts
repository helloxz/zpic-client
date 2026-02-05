import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './style.css'


// 创建 Vue 应用实例
const app = createApp(App)
app.use(createPinia())
// 使用路由
app.use(router)

// 挂载应用
app.mount('#app')
