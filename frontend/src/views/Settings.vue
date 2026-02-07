<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { NCard, NForm, NFormItem, NInput, NSelect, NButton, NIcon, NAlert } from 'naive-ui'
import { SaveOutline, RefreshOutline } from '@vicons/ionicons5'
import { UpdateSetting, GetSetting } from '../../wailsjs/go/core/AppCore'
import { core } from '../../wailsjs/go/models'
import { useMessage } from 'naive-ui'
import axios from 'axios'
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime"

const message = useMessage()

interface SettingForm {
  server: string
  baseUrl: string
  token: string
  httpProxy: string
}

const serverOptions = [
  { label: 'ImgURL (www.imgurl.org)', value: 'imgurl' },
  { label: '图链 (go.piclink.cc)', value: 'piclink' },
  { label: '自定义', value: 'other' }
]

const serverBaseUrls: Record<string, string> = {
  imgurl: 'https://www.imgurl.org',
  piclink: 'https://go.piclink.cc'
}

const formRef = ref<any>(null)
const loading = ref(false)
const saving = ref(false)
const validating = ref(false)

const form = ref<SettingForm>({
  server: 'imgurl',
  baseUrl: '',
  token: '',
  httpProxy: ''
})

const rules = {
  baseUrl: [
    {
      validator: (_: unknown, value: string) => {
        if (!value && form.value.server === 'other') {
          return new Error('请输入服务器地址')
        }
        if (value && !/^https?:\/\//i.test(value)) {
          return new Error('必须以 http:// 或 https:// 开头')
        }
        if (value && value.endsWith('/')) {
          return new Error('末尾不能有 /')
        }
        return true
      },
      trigger: 'blur'
    }
  ],
  token: [
    { required: true, message: '请输入Token' }
  ],
  httpProxy: [
    {
      validator: (_: unknown, value: string) => {
        if (value && !/^https?:\/\//i.test(value)) {
          return new Error('必须以 http:// 或 https:// 开头')
        }
        return true
      },
      trigger: 'blur'
    }
  ]
}

const showBaseUrlInput = ref(false)

const initSettings = async () => {
  loading.value = true
  try {
    const data: core.SettingData = await GetSetting()
    if (data.base_url) {
      if (data.base_url === serverBaseUrls.imgurl) {
        form.value.server = 'imgurl'
        form.value.baseUrl = serverBaseUrls.imgurl
        showBaseUrlInput.value = false
      } else if (data.base_url === serverBaseUrls.piclink) {
        form.value.server = 'piclink'
        form.value.baseUrl = serverBaseUrls.piclink
        showBaseUrlInput.value = false
      } else {
        form.value.server = 'other'
        form.value.baseUrl = data.base_url
        showBaseUrlInput.value = true
      }
    } else {
      form.value.server = 'imgurl'
      form.value.baseUrl = serverBaseUrls.imgurl
      showBaseUrlInput.value = false
    }
    form.value.token = data.token || ''
    form.value.httpProxy = data.http_proxy || ''
  } catch (err) {
    console.error('获取设置失败:', err)
  } finally {
    loading.value = false
  }
}

const handleServerChange = (value: string) => {
  if (value === 'other') {
    showBaseUrlInput.value = true
    form.value.baseUrl = ''
  } else {
    showBaseUrlInput.value = false
    form.value.baseUrl = serverBaseUrls[value]
  }
}

watch(() => form.value.server, (newVal) => {
  handleServerChange(newVal)
})


const handleSave = async () => {
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  const baseUrl = form.value.server === 'other'
    ? form.value.baseUrl
    : serverBaseUrls[form.value.server as keyof typeof serverBaseUrls]

  validating.value = true
  try {
    const response = await axios.get(`${baseUrl}/api/v3/album_list`, {
      headers: {
        Authorization: `Bearer ${form.value.token}`
      },
      timeout: 10000
    })

    if (response.data.code !== 200) {
      message.error(response.data.msg || 'Token 无效')
      validating.value = false
      return
    }
  } catch (err: any) {
    const errorMsg = err.response?.data?.msg || 'Token 验证失败，请检查服务器地址和 Token'
    message.error(errorMsg)
    validating.value = false
    return
  }

  saving.value = true
  try {
    const success = await UpdateSetting({
      base_url: baseUrl,
      token: form.value.token,
      http_proxy: form.value.httpProxy
    })
    if (success) {
      // 清理相册的sessionStorage缓存
      sessionStorage.removeItem('albumListCache')
      message.success('设置已保存')
    } else {
      message.error('设置保存失败')
    }
  } catch (err) {
    console.error('保存设置失败:', err)
    message.error('设置保存失败')
  } finally {
    saving.value = false
    validating.value = false
  }
}

const handleReset = () => {
  initSettings()
}

onMounted(() => {
  initSettings()
})
</script>

<template>
  <div class="settings">
    <div class="page-header">
      <h1>设置</h1>
      <p class="subtitle">配置图床服务器和Token信息</p>
    </div>

    <NCard class="setting-card">
      <NForm
        ref="formRef"
        :model="form"
        :rules="rules"
        label-placement="top"
        :loading="loading"
      >
        <NFormItem label="选择服务器" path="server">
          <NSelect
            v-model:value="form.server"
            :options="serverOptions"
            placeholder="请选择服务器"
          />
        </NFormItem>

        <NFormItem v-if="showBaseUrlInput" label="服务器地址" path="baseUrl">
          <NInput
            v-model:value="form.baseUrl"
            placeholder="请输入服务器地址，如：https://your-server.com"
          />
        </NFormItem>

        <NFormItem label="Token" path="token">
          <NInput
            v-model:value="form.token"
            type="password"
            show-password-on="click"
            placeholder="请输入API Token"
          />
        </NFormItem>

        <NFormItem label="HTTP 代理（可选）" path="httpProxy">
          <NInput
            v-model:value="form.httpProxy"
            placeholder="请输入HTTP代理地址，如：http://127.0.0.1:7890"
          />
        </NFormItem>

        <NAlert type="info" class="info-alert">
          请根据所选服务器类型，在对应的图床网站上获取API Token。如果您还没有账号，可前往 <a href="javascript:;" @click="BrowserOpenURL('https://www.imgurl.org/user/register')" title="ImgURL图床">ImgURL</a> 注册。
        </NAlert>
      </NForm>

      <div class="action-buttons">
        <NButton type="primary" @click="handleSave" :loading="saving || validating">
          <template #icon>
            <NIcon>
              <SaveOutline />
            </NIcon>
          </template>
          保存设置
        </NButton>
        <!-- <NButton @click="handleReset" :loading="loading">
          <template #icon>
            <NIcon>
              <RefreshOutline />
            </NIcon>
          </template>
          重置
        </NButton> -->
      </div>
    </NCard>
  </div>
</template>

<style scoped>
.settings {
  padding: 18px;
  max-width: 100%;
  margin: 0 auto;
}

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

.setting-card {
  border-radius: 12px;
}

.info-alert {
  margin-top: 8px;
  margin-bottom: 24px;
}

.info-alert a {
  color: #4098fc;
  text-decoration: underline;
  font-weight: 600;
  cursor: pointer;
}

.info-alert a:hover {
  color: #1d77e5;
}

.action-buttons {
  display: flex;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}
</style>
