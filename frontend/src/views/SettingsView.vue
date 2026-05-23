<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useProviderStore } from '@/stores/provider'
import { getCookies, saveCookie, deleteCookie } from '@/api/cookie'
import ProviderList from '@/components/provider/ProviderList.vue'
import ProviderDialog from '@/components/provider/ProviderDialog.vue'

const PLATFORM_LABELS: Record<string, string> = {
  bilibili: 'B站',
  youtube: 'YouTube',
  douyin: '抖音',
  kuaishou: '快手',
}

const router = useRouter()
const providerStore = useProviderStore()
const showDialog = ref(false)
const editingProvider = ref<LLMProviderData | null>(null)

const cookies = ref<CookieData[]>([])
const cookieLoading = ref(false)
const editingCookie = ref<{ platform: string; content: string } | null>(null)
const cookieContent = ref('')
const cookieFormat = ref<'header' | 'netscape'>('header')

onMounted(() => {
  providerStore.fetchProviders()
  fetchCookies()
})

async function fetchCookies() {
  cookieLoading.value = true
  try {
    const res = await getCookies()
    cookies.value = res.data.data || []
  } catch { /* ignore */ }
  finally { cookieLoading.value = false }
}

function handleAddCookie() {
  editingCookie.value = { platform: '', content: '' }
  cookieContent.value = ''
}

function handleEditCookie(c: CookieData) {
  editingCookie.value = { platform: c.platform, content: c.content }
  cookieContent.value = c.content
}

async function handleSaveCookie() {
  if (!editingCookie.value) return
  const platform = editingCookie.value.platform
  const content = cookieContent.value.trim()
  if (!platform || !content) {
    alert('请填写平台和 Cookie 内容')
    return
  }
  try {
    await saveCookie(platform, content)
    editingCookie.value = null
    await fetchCookies()
  } catch (e: any) {
    alert('保存失败: ' + e.message)
  }
}

async function handleDeleteCookie(platform: string) {
  if (!confirm(`确定删除 ${PLATFORM_LABELS[platform] || platform} 的 Cookie？`)) return
  try {
    await deleteCookie(platform)
    await fetchCookies()
  } catch (e: any) {
    alert('删除失败: ' + e.message)
  }
}

function handleAdd() {
  editingProvider.value = null
  showDialog.value = true
}

function handleEdit(provider: LLMProviderData) {
  editingProvider.value = provider
  showDialog.value = true
}

async function handleSave(data: ProviderRequest | ProviderUpdateRequest) {
  try {
    if (editingProvider.value) {
      await providerStore.saveProvider(editingProvider.value.id, data)
    } else {
      await providerStore.addProvider(data as ProviderRequest)
    }
    showDialog.value = false
  } catch (e: any) {
    alert(e.message)
  }
}

async function handleDelete(id: string) {
  if (!confirm('确定删除此提供商？')) return
  await providerStore.removeProvider(id)
}

function goBack() {
  router.push('/')
}
</script>

<template>
  <div class="settings">
    <div class="settings-header">
      <button class="back-btn" @click="goBack">← 返回</button>
      <h2>设置</h2>
    </div>

    <div class="section">
      <div class="section-header">
        <h3>平台 Cookie</h3>
        <span class="section-hint">为各平台提供 Cookie 可避免 403 反爬拦截</span>
      </div>

      <div v-if="cookies.length > 0" class="cookie-list">
        <div v-for="c in cookies" :key="c.platform" class="cookie-item">
          <div class="cookie-info">
            <span class="cookie-platform">{{ PLATFORM_LABELS[c.platform] || c.platform }}</span>
            <span class="cookie-preview">{{ c.content.substring(0, 60) }}{{ c.content.length > 60 ? '...' : '' }}</span>
          </div>
          <div class="cookie-actions">
            <button class="action-btn" @click="handleEditCookie(c)">编辑</button>
            <button class="action-btn danger" @click="handleDeleteCookie(c.platform)">删除</button>
          </div>
        </div>
      </div>
      <div v-else class="cookie-empty">
        暂无 Cookie 配置，添加后可避免部分平台的反爬拦截
      </div>

      <button class="add-btn" style="margin-top: 12px" @click="handleAddCookie">+ 添加 Cookie</button>

      <div v-if="editingCookie" class="cookie-editor">
        <label>平台</label>
        <select v-model="editingCookie.platform" class="select">
          <option value="">选择平台</option>
          <option value="bilibili">B站</option>
          <option value="youtube">YouTube</option>
          <option value="douyin">抖音</option>
          <option value="kuaishou">快手</option>
        </select>

        <label>Cookie 格式</label>
        <div class="format-tabs">
          <button
            type="button"
            :class="['format-tab', { active: cookieFormat === 'header' }]"
            @click="cookieFormat = 'header'"
          >
            请求头格式（推荐）
          </button>
          <button
            type="button"
            :class="['format-tab', { active: cookieFormat === 'netscape' }]"
            @click="cookieFormat = 'netscape'"
          >
            Netscape 格式
          </button>
        </div>

        <label v-if="cookieFormat === 'header'">粘贴请求头中的 Cookie 值</label>
        <label v-else>粘贴 Netscape 格式 Cookie 文本</label>
        <textarea
          v-model="cookieContent"
          class="cookie-textarea"
          rows="6"
          :placeholder="cookieFormat === 'header'
            ? '打开任意 bilibili.com 页面 → F12 → Network → 点一个请求 → 找到 Request Headers → 复制 Cookie: 后面全部内容'
            : 'Chrome 插件 Get cookies.txt LOCALLY 导出后粘贴内容'" />
        <div class="cookie-actions">
          <button class="action-btn primary" @click="handleSaveCookie">保存</button>
          <button class="action-btn" @click="editingCookie = null">取消</button>
        </div>
        <div v-if="cookieFormat === 'header'" class="cookie-help">
          如何获取：
          <ol>
            <li>浏览器打开 bilibili.com 并登录</li>
            <li>按 <code>F12</code> → <code>Network</code>（网络）标签</li>
            <li>刷新页面，点击任意请求</li>
            <li>在 <code>Request Headers</code> 中找到 <code>Cookie:</code> 一行</li>
            <li>复制冒号后面的全部内容粘贴到这里</li>
          </ol>
        </div>
        <div v-else class="cookie-help">
          获取方式：
          <ul>
            <li>Chrome 插件: <code>Get cookies.txt LOCALLY</code> 导出</li>
          </ul>
        </div>
      </div>
    </div>

    <div class="section" style="margin-top: 24px">
      <div class="section-header">
        <h3>LLM 模型提供商</h3>
        <button class="add-btn" @click="handleAdd">+ 添加提供商</button>
      </div>
      <ProviderList
        :providers="providerStore.providers"
        :loading="providerStore.loading"
        @edit="handleEdit"
        @delete="handleDelete"
      />
    </div>

    <ProviderDialog
      v-if="showDialog"
      :provider="editingProvider"
      @save="handleSave"
      @close="showDialog = false"
    />
  </div>
</template>

<style scoped>
.settings {
  max-width: 800px;
  margin: 0 auto;
  padding: 24px;
}
.settings-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}
.back-btn {
  padding: 6px 16px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: #fff;
  cursor: pointer;
  font-size: 14px;
}
.back-btn:hover {
  background: #f5f5f5;
}
h2 {
  font-size: 22px;
  font-weight: 600;
}
.section {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.section-header h3 {
  font-size: 16px;
  font-weight: 600;
}
.section-hint {
  font-size: 12px;
  color: #999;
}
.add-btn {
  padding: 6px 16px;
  border: 1px solid #4A90D9;
  border-radius: 6px;
  background: #4A90D9;
  color: #fff;
  cursor: pointer;
  font-size: 14px;
}
.add-btn:hover {
  background: #357ABD;
}
.select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  margin: 8px 0 12px;
}

.cookie-list {
  margin-bottom: 8px;
}
.cookie-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
}
.cookie-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.cookie-platform {
  font-weight: 600;
  font-size: 14px;
}
.cookie-preview {
  font-size: 12px;
  color: #999;
  font-family: monospace;
}
.cookie-empty {
  color: #999;
  font-size: 13px;
  margin-bottom: 8px;
}
.cookie-editor {
  margin-top: 16px;
  padding: 16px;
  background: #f9f9f9;
  border-radius: 8px;
}
.cookie-editor label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 4px;
}
.format-tabs {
  display: flex;
  gap: 0;
  margin-bottom: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  overflow: hidden;
}
.format-tab {
  flex: 1;
  padding: 6px 0;
  border: none;
  background: #f5f5f5;
  cursor: pointer;
  font-size: 13px;
  color: #666;
  transition: background 0.2s;
}
.format-tab:first-child {
  border-right: 1px solid #ddd;
}
.format-tab.active {
  background: #4A90D9;
  color: #fff;
}
.cookie-textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-family: monospace;
  font-size: 12px;
  resize: vertical;
  box-sizing: border-box;
}
.cookie-help {
  margin-top: 12px;
  font-size: 12px;
  color: #666;
}
.cookie-help ul {
  margin: 4px 0 0 16px;
  padding: 0;
}
.cookie-help code {
  background: #eee;
  padding: 1px 4px;
  border-radius: 3px;
  font-size: 11px;
}

.action-btn {
  padding: 4px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
  cursor: pointer;
  font-size: 12px;
  margin-left: 8px;
}
.action-btn:hover {
  background: #f0f0f0;
}
.action-btn.primary {
  background: #4A90D9;
  color: #fff;
  border-color: #4A90D9;
}
.action-btn.primary:hover {
  background: #357ABD;
}
.action-btn.danger {
  color: #e74c3c;
  border-color: #e74c3c;
}
.action-btn.danger:hover {
  background: #fdf0ef;
}
.cookie-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 8px;
}
</style>
