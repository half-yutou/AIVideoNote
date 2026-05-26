<script setup lang="ts">
import { reactive, ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/task'
import { useProviderStore } from '@/stores/provider'
import { useGroupStore } from '@/stores/group'
import { useToast } from '@/composables/useToast'
import { listModels, testConnection } from '@/api/provider'
import { uploadFile } from '@/api/upload'
import NoteForm from '@/components/note/NoteForm.vue'
import MarkdownViewer from '@/components/note/MarkdownViewer.vue'
import StepBar from '@/components/note/StepBar.vue'
import TaskHistory from '@/components/task/TaskHistory.vue'

const STORAGE_KEY = 'videonote_form_config'

const router = useRouter()
const taskStore = useTaskStore()
const providerStore = useProviderStore()
const groupStore = useGroupStore()
const toast = useToast()

const models = ref<string[]>([])
const selectedModel = ref('')
const currentProviderId = ref('')
const testing = ref(false)
const localFile = ref<File | null>(null)
const uploadProgress = ref(0)
const uploading = ref(false)

function loadSavedConfig() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) return JSON.parse(raw)
  } catch { /* ignore */ }
  return {}
}

const saved = loadSavedConfig()

const form = reactive<GenerateRequest>({
  video_url: saved.video_url || '',
  platform: saved.platform || 'bilibili',
  quality: saved.quality || 'medium',
  model_name: saved.model_name || '',
  provider_id: saved.provider_id || '',
  style: saved.style || '',
  group_id: saved.group_id || 'default',
  link: saved.link || false,
})

watch(form, (val) => {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(val))
}, { deep: true })

async function handleProviderChange(providerId: string) {
  currentProviderId.value = providerId
  form.provider_id = providerId
  form.model_name = ''
  selectedModel.value = ''
  models.value = []
  if (!providerId) return
  try {
    const res = await listModels(providerId)
    models.value = res.data.data || []
  } catch {
    toast.error('获取模型列表失败')
  }
}

function handleModelSelect(model: string) {
  selectedModel.value = model
  form.model_name = model
}

async function handleSubmit() {
  if (form.platform === 'local') {
    // 本地文件模式：先上传文件获取服务端路径
    if (!localFile.value) {
      toast.warning('请选择要上传的视频/音频文件')
      return
    }
    if (!form.model_name || !form.provider_id) {
      toast.warning('请填写完整信息')
      return
    }
    try {
      uploading.value = true
      uploadProgress.value = 0
      toast.info('正在上传文件...')
      const res = await uploadFile(localFile.value, (p) => { uploadProgress.value = p })
      const serverPath = res.data.data?.path
      if (!serverPath) {
        toast.error('上传失败：未获取服务端路径')
        return
      }
      // 上传成功，用服务端路径提交任务
      await taskStore.submitTask({ ...form, video_url: serverPath })
      toast.success('任务已提交，正在处理中...')
    } catch (e: any) {
      toast.error(e?.response?.data?.message || e.message || '上传或提交失败')
    } finally {
      uploading.value = false
      uploadProgress.value = 0
    }
  } else {
    // B站链接模式
    if (!form.video_url || !form.platform || !form.model_name || !form.provider_id) {
      toast.warning('请填写完整信息')
      return
    }
    try {
      await taskStore.submitTask({ ...form })
      toast.success('任务已提交，正在处理中...')
    } catch (e: any) {
      toast.error(e.message || '提交失败')
    }
  }
}

function handleLocalFile(file: File | null) {
  localFile.value = file
}

async function handleCreateGroup(name: string) {
  try {
    const g = await groupStore.addGroup(name)
    form.group_id = g.id
    toast.success(`分组「${name}」创建成功`)
  } catch {
    toast.error('创建分组失败')
  }
}

function handleTaskSelect(taskId: string) {
  taskStore.loadTask(taskId)
}

function toSettings() {
  router.push('/settings')
}

function isFailedStatus(status: string) {
  return status === 'FAILED' || status.endsWith('_FAILED')
}

async function handleTestConnection() {
  if (!currentProviderId.value) return
  testing.value = true
  try {
    const res = await testConnection(currentProviderId.value)
    toast.success(res.data.message || '连接成功')
  } catch (e: any) {
    toast.error('连接失败: ' + e.message)
  } finally {
    testing.value = false
  }
}

onMounted(async () => {
  await providerStore.fetchProviders()
  groupStore.fetchGroups()
  taskStore.fetchTasks()
  const enabled = providerStore.providers.filter((p) => p.enabled)
  if (enabled.length === 0) return
  if (form.provider_id) {
    const found = enabled.find((p) => p.id === form.provider_id)
    if (found) {
      currentProviderId.value = found.id
      try {
        const res = await listModels(found.id)
        models.value = res.data.data || []
        if (form.model_name && models.value.includes(form.model_name)) {
          selectedModel.value = form.model_name
        }
      } catch { /* ignore */ }
      return
    }
  }
  const first = enabled[0]
  currentProviderId.value = first.id
  form.provider_id = first.id
  try {
    const res = await listModels(first.id)
    models.value = res.data.data || []
    if (form.model_name && models.value.includes(form.model_name)) {
      selectedModel.value = form.model_name
    }
  } catch { /* ignore */ }
})
</script>

<template>
  <div class="home">
    <div class="home-main">
      <aside class="sidebar">
      <div class="sidebar-top">
        <button class="settings-btn" @click="toSettings">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="3"/>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68 1.65 1.65 0 0 0 9 3V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
          </svg>
          <span>设置</span>
          <span class="settings-hint">提供商 · Cookie</span>
        </button>
      </div>

      <NoteForm
        :form="form"
        :models="models"
        :selected-model="selectedModel"
        :groups="groupStore.groups"
        :errors="{ error: taskStore.error }"
        :loading="taskStore.loading || uploading"
        :testing="testing"
        :has-provider="currentProviderId !== ''"
        @submit="handleSubmit"
        @update:platform="(v: string) => form.platform = v"
        @update:video-url="(v: string) => form.video_url = v"
        @update:style="(v: string) => form.style = v"
        @update:group-id="(v: string) => form.group_id = v"
        @update:local-file="handleLocalFile"
        @provider-change="handleProviderChange"
        @model-select="handleModelSelect"
        @test-connection="handleTestConnection"
        @create-group="handleCreateGroup"
      />

    </aside>

    <section class="content">
      <div v-if="!taskStore.currentMarkdown && !taskStore.loading" class="empty-state">
        <div class="empty-illustration">
          <svg width="80" height="80" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
            <polygon points="23 7 16 12 23 17 23 7"/>
            <rect x="1" y="5" width="15" height="14" rx="2" ry="2"/>
          </svg>
        </div>
        <h2 class="empty-title">AI 视频笔记</h2>
        <p class="empty-desc">粘贴视频链接，AI 自动生成结构化笔记</p>
        <div class="empty-steps">
          <div class="empty-step">
            <span class="step-num">1</span>
            <span>粘贴视频链接</span>
          </div>
          <div class="empty-step-arrow">→</div>
          <div class="empty-step">
            <span class="step-num">2</span>
            <span>AI 转录 & 分析</span>
          </div>
          <div class="empty-step-arrow">→</div>
          <div class="empty-step">
            <span class="step-num">3</span>
            <span>获取结构化笔记</span>
          </div>
        </div>
      </div>

      <div v-if="taskStore.loading && !taskStore.currentMarkdown" class="loading-state">
        <div class="loading-spinner"></div>
        <p>正在处理中，请稍候...</p>
      </div>

      <MarkdownViewer
        v-if="taskStore.currentMarkdown"
        :markdown="taskStore.currentMarkdown"
      />

      <div v-if="taskStore.currentTaskId && taskStore.currentMarkdown" class="actions">
        <a :href="`/api/v1/export/${taskStore.currentTaskId}?format=md`" class="btn-export" download>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/>
          </svg>
          导出 Markdown
        </a>
      </div>
    </section>

    <aside class="history-panel">
      <div class="history-header">
        <span class="history-title">历史记录</span>
        <span class="history-count">{{ taskStore.tasks.length }}</span>
      </div>
      <TaskHistory
        :tasks="taskStore.tasks"
        @select="handleTaskSelect"
        @delete="(id: string) => taskStore.removeTask(id)"
      />
    </aside>
    </div>

    <!-- 底部进度条 -->
    <Transition name="slide-up">
      <div v-if="taskStore.currentStatus" class="bottom-bar">
        <StepBar :status="taskStore.currentStatus" :error-message="taskStore.currentError" />
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  height: calc(100vh - var(--header-height));
  overflow: hidden;
}

.home-main {
  display: flex;
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

/* 底部进度条 */
.bottom-bar {
  width: 100%;
  background: var(--color-bg-elevated);
  border-top: 1px solid var(--color-border);
  flex-shrink: 0;
}

.slide-up-enter-active {
  transition: all 300ms ease-out;
}
.slide-up-leave-active {
  transition: all 200ms ease-in;
}
.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(8px);
}

.sidebar {
  width: var(--sidebar-width);
  min-width: var(--sidebar-width);
  background: var(--color-bg-elevated);
  border-right: 1px solid var(--color-border);
  overflow-y: auto;
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.sidebar-top {
  margin-bottom: var(--space-1);
}

.settings-btn {
  width: 100%;
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg-muted);
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
  transition: all var(--transition-fast);
}
.settings-btn:hover {
  border-color: var(--color-primary);
  background: var(--color-primary-light);
  color: var(--color-primary);
}
.settings-hint {
  margin-left: auto;
  font-size: var(--text-xs);
  color: var(--color-text-muted);
}

.content {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-6);
  min-width: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

/* 空状态 */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-4);
  animation: fadeInUp 600ms ease-out;
}
.empty-illustration {
  color: var(--color-primary-subtle);
  margin-bottom: var(--space-2);
}
.empty-title {
  font-size: var(--text-2xl);
  font-weight: 700;
  color: var(--color-text);
  letter-spacing: -0.02em;
}
.empty-desc {
  font-size: var(--text-md);
  color: var(--color-text-secondary);
}
.empty-steps {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-top: var(--space-6);
  padding: var(--space-5) var(--space-6);
  background: var(--color-bg-muted);
  border-radius: var(--radius-xl);
}
.empty-step {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}
.step-num {
  width: 24px;
  height: 24px;
  border-radius: var(--radius-full);
  background: var(--color-primary);
  color: var(--color-text-inverse);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-xs);
  font-weight: 600;
}
.empty-step-arrow {
  color: var(--color-text-muted);
  font-size: var(--text-md);
}

/* 加载状态 */
.loading-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-4);
  color: var(--color-text-secondary);
}
.loading-spinner {
  width: 36px;
  height: 36px;
  border: 3px solid var(--color-bg-subtle);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

/* 历史面板 */
.history-panel {
  width: var(--history-width);
  min-width: var(--history-width);
  background: var(--color-bg-elevated);
  border-left: 1px solid var(--color-border);
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}
.history-header {
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: sticky;
  top: 0;
  background: var(--color-bg-elevated);
  z-index: 1;
}
.history-title {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-secondary);
}
.history-count {
  font-size: var(--text-xs);
  color: var(--color-text-muted);
  background: var(--color-bg-muted);
  padding: 2px 8px;
  border-radius: var(--radius-full);
  font-weight: 500;
}

/* 操作按钮 */
.actions {
  display: flex;
  gap: var(--space-3);
  margin-top: var(--space-4);
  padding-top: var(--space-4);
  border-top: 1px solid var(--color-border-light);
}
.btn-export {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-5);
  background: var(--color-success);
  color: var(--color-text-inverse);
  border-radius: var(--radius-md);
  font-size: var(--text-base);
  font-weight: 500;
  text-decoration: none;
  transition: all var(--transition-fast);
}
.btn-export:hover {
  background: #059669;
  color: var(--color-text-inverse);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}


</style>
