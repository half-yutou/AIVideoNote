<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useProviderStore } from '@/stores/provider'
import { useToast } from '@/composables/useToast'
import { getCookies, saveCookie, deleteCookie } from '@/api/cookie'
import ProviderList from '@/components/provider/ProviderList.vue'
import ProviderDialog from '@/components/provider/ProviderDialog.vue'

const PLATFORM_LABELS: Record<string, string> = {
  bilibili: 'B站',
}

const router = useRouter()
const providerStore = useProviderStore()
const toast = useToast()
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
    toast.warning('请填写平台和 Cookie 内容')
    return
  }
  try {
    await saveCookie(platform, content)
    editingCookie.value = null
    await fetchCookies()
    toast.success('Cookie 保存成功')
  } catch (e: any) {
    toast.error('保存失败: ' + e.message)
  }
}

async function handleDeleteCookie(platform: string) {
  if (!confirm(`确定删除 ${PLATFORM_LABELS[platform] || platform} 的 Cookie？`)) return
  try {
    await deleteCookie(platform)
    await fetchCookies()
    toast.success('Cookie 已删除')
  } catch (e: any) {
    toast.error('删除失败: ' + e.message)
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
    toast.success(editingProvider.value ? '提供商已更新' : '提供商已添加')
  } catch (e: any) {
    toast.error(e.message)
  }
}

async function handleDelete(id: string) {
  if (!confirm('确定删除此提供商？')) return
  await providerStore.removeProvider(id)
  toast.success('提供商已删除')
}

function goBack() {
  router.push('/')
}
</script>

<template>
  <div class="settings">
    <div class="settings-header">
      <button class="back-btn" @click="goBack">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="19" y1="12" x2="5" y2="12"/><polyline points="12 19 5 12 12 5"/>
        </svg>
        返回
      </button>
      <div class="settings-title-group">
        <h2>设置</h2>
        <span class="settings-subtitle">管理提供商和平台配置</span>
      </div>
    </div>

    <!-- Cookie 配置 -->
    <div class="section animate-fade-in-up">
      <div class="section-header">
        <div class="section-title-group">
          <h3>
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 2a10 10 0 1 0 10 10 4 4 0 0 1-5-5 4 4 0 0 1-5-5"/>
              <path d="M8.5 8.5v.01"/><path d="M16 15.5v.01"/><path d="M12 12v.01"/><path d="M11 17v.01"/><path d="M7 14v.01"/>
            </svg>
            平台 Cookie
          </h3>
          <span class="section-hint">为各平台提供 Cookie 可避免 403 反爬拦截</span>
        </div>
        <button class="add-btn" @click="handleAddCookie">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          添加
        </button>
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
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
        <span>暂无 Cookie 配置，添加后可避免部分平台的反爬拦截</span>
      </div>

      <Transition name="slide">
        <div v-if="editingCookie" class="cookie-editor">
          <div class="editor-field">
            <label>平台</label>
            <select v-model="editingCookie.platform" class="form-select">
              <option value="">选择平台</option>
              <option value="bilibili">B站</option>
            </select>
          </div>

          <div class="editor-field">
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
          </div>

          <div class="editor-field">
            <label v-if="cookieFormat === 'header'">粘贴请求头中的 Cookie 值</label>
            <label v-else>粘贴 Netscape 格式 Cookie 文本</label>
            <textarea
              v-model="cookieContent"
              class="cookie-textarea"
              rows="5"
              :placeholder="cookieFormat === 'header'
                ? '打开 bilibili.com → F12 → Network → 复制 Cookie 值'
                : 'Chrome 插件 Get cookies.txt LOCALLY 导出后粘贴'" />
          </div>

          <div class="editor-actions">
            <button class="action-btn primary" @click="handleSaveCookie">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="20 6 9 17 4 12"/>
              </svg>
              保存
            </button>
            <button class="action-btn" @click="editingCookie = null">取消</button>
          </div>

          <div class="cookie-help">
            <template v-if="cookieFormat === 'header'">
              <strong>如何获取：</strong>
              <ol>
                <li>浏览器打开 bilibili.com 并登录</li>
                <li>按 <code>F12</code> → <code>Network</code> 标签</li>
                <li>刷新页面，点击任意请求</li>
                <li>在 Request Headers 中找到 <code>Cookie:</code> 一行</li>
                <li>复制冒号后面的全部内容粘贴到这里</li>
              </ol>
            </template>
            <template v-else>
              <strong>获取方式：</strong>
              <p>Chrome 插件 <code>Get cookies.txt LOCALLY</code> 导出后粘贴内容</p>
            </template>
          </div>
        </div>
      </Transition>
    </div>

    <!-- LLM 提供商 -->
    <div class="section animate-fade-in-up" style="animation-delay: 100ms">
      <div class="section-header">
        <div class="section-title-group">
          <h3>
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
              <polyline points="3.27 6.96 12 12.01 20.73 6.96"/>
              <line x1="12" y1="22.08" x2="12" y2="12"/>
            </svg>
            LLM 模型提供商
          </h3>
        </div>
        <button class="add-btn" @click="handleAdd">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          添加提供商
        </button>
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
  padding: var(--space-6);
}

.settings-header {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  margin-bottom: var(--space-6);
}

.back-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg-elevated);
  color: var(--color-text-secondary);
  font-size: var(--text-base);
  transition: all var(--transition-fast);
}
.back-btn:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
  background: var(--color-primary-light);
}

.settings-title-group {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.settings-title-group h2 {
  font-size: var(--text-xl);
  font-weight: 700;
  color: var(--color-text);
}
.settings-subtitle {
  font-size: var(--text-sm);
  color: var(--color-text-muted);
}

.section {
  background: var(--color-bg-elevated);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border);
  margin-bottom: var(--space-5);
}

.section-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: var(--space-4);
}

.section-title-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}
.section-title-group h3 {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-md);
  font-weight: 600;
  color: var(--color-text);
}
.section-hint {
  font-size: var(--text-xs);
  color: var(--color-text-muted);
}

.add-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-3);
  border: none;
  border-radius: var(--radius-md);
  background: var(--color-primary);
  color: var(--color-text-inverse);
  font-size: var(--text-sm);
  font-weight: 500;
  transition: all var(--transition-fast);
}
.add-btn:hover {
  background: var(--color-primary-hover);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

.form-select {
  width: 100%;
  padding: 10px 14px;
  padding-right: 32px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--text-base);
  background: var(--color-bg-elevated);
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%2364748b' d='M2.5 4.5L6 8l3.5-3.5'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  transition: all var(--transition-fast);
}
.form-select:focus {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-focus);
  outline: none;
}

/* Cookie 列表 */
.cookie-list {
  margin-bottom: var(--space-2);
}
.cookie-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-3) 0;
  border-bottom: 1px solid var(--color-border-light);
}
.cookie-item:last-child {
  border-bottom: none;
}
.cookie-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}
.cookie-platform {
  font-weight: 600;
  font-size: var(--text-base);
  color: var(--color-text);
}
.cookie-preview {
  font-size: var(--text-xs);
  color: var(--color-text-muted);
  font-family: var(--font-mono);
}
.cookie-empty {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  color: var(--color-text-muted);
  font-size: var(--text-sm);
  padding: var(--space-3) 0;
}

/* Cookie 编辑器 */
.cookie-editor {
  margin-top: var(--space-4);
  padding: var(--space-4);
  background: var(--color-bg-muted);
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}
.editor-field {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}
.editor-field label {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-secondary);
}

.format-tabs {
  display: flex;
  gap: 0;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
}
.format-tab {
  flex: 1;
  padding: var(--space-2) 0;
  background: var(--color-bg-elevated);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
  border-right: 1px solid var(--color-border);
}
.format-tab:last-child {
  border-right: none;
}
.format-tab.active {
  background: var(--color-primary);
  color: var(--color-text-inverse);
}

.cookie-textarea {
  width: 100%;
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  resize: vertical;
  box-sizing: border-box;
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
  line-height: 1.5;
}
.cookie-textarea:focus {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-focus);
  outline: none;
}

.editor-actions {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.cookie-help {
  font-size: var(--text-xs);
  color: var(--color-text-muted);
  line-height: 1.6;
  padding: var(--space-3);
  background: var(--color-bg-elevated);
  border-radius: var(--radius-md);
}
.cookie-help strong {
  color: var(--color-text-secondary);
}
.cookie-help ol, .cookie-help ul {
  margin: var(--space-1) 0 0 var(--space-4);
  padding: 0;
}
.cookie-help code {
  background: var(--color-bg-muted);
  padding: 1px 5px;
  border-radius: 3px;
  font-size: 11px;
  font-family: var(--font-mono);
}

/* 操作按钮 */
.action-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-elevated);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
}
.action-btn:hover {
  background: var(--color-bg-muted);
  border-color: var(--color-text-muted);
}
.action-btn.primary {
  background: var(--color-primary);
  color: var(--color-text-inverse);
  border-color: var(--color-primary);
}
.action-btn.primary:hover {
  background: var(--color-primary-hover);
}
.action-btn.danger {
  color: var(--color-danger);
  border-color: var(--color-danger);
}
.action-btn.danger:hover {
  background: var(--color-danger-light);
}
.cookie-actions {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

/* 过渡动画 */
.slide-enter-active {
  transition: all 300ms ease-out;
}
.slide-leave-active {
  transition: all 200ms ease-in;
}
.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
