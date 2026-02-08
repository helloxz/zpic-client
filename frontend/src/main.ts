import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import i18n from './i18n'
import './style.css'


// 创建 Vue 应用实例
const app = createApp(App)
app.use(createPinia())
app.use(i18n)
// 使用路由
app.use(router)

// 挂载应用
app.mount('#app')
