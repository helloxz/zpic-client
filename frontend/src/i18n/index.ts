import { createI18n } from 'vue-i18n'

// 同步加载语言包（默认中文）
import zh from './locales/zh.ts'
import en from './locales/en.ts'

// 获取用户设置的语言
function getUserLanguage() {
  // 1. 优先从 localStorage 获取用户设置的语言
  const userLang = localStorage.getItem('user_language')
  if (userLang && (userLang === 'zh' || userLang === 'en')) {
    return userLang
  }
  
  // 2. 如果没有用户设置，则根据浏览器语言判断
  const browserLang = navigator.language || navigator.language
  if (browserLang && browserLang.startsWith('zh')) {
    return 'zh'
  }
  
  // 3. 其他情况默认英文
  return 'en'
}

const i18n = createI18n({
    legacy: false,
  locale: getUserLanguage(),           // 根据优先级设置语言
  fallbackLocale: 'en',   // 备选语言
  messages: {
    zh,
    en
  },
  // 可选：关闭一些警告（生产环境更干净）
  silentFallbackWarn: true,
  silentTranslationWarn: true
})

export default i18n