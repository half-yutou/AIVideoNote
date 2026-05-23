<script setup lang="ts">
import { reactive, ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/task'
import { useProviderStore } from '@/stores/provider'
import { useGroupStore } from '@/stores/group'
import { listModels, testConnection } from '@/api/provider'
import NoteForm from '@/components/note/NoteForm.vue'
import MarkdownViewer from '@/components/note/MarkdownViewer.vue'
import StepBar from '@/components/note/StepBar.vue'
import TaskHistory from '@/components/task/TaskHistory.vue'

const STORAGE_KEY = 'videonote_form_config'

const router = useRouter()
const taskStore = useTaskStore()
const providerStore = useProviderStore()
const groupStore = useGroupStore()

const models = ref<string[]>([])
const selectedModel = ref('')
const currentProviderId = ref('')
const testing = ref(false)

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
  screenshot: saved.screenshot || false,
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
    // ignore
  }
}

function handleModelSelect(model: string) {
  selectedModel.value = model
  form.model_name = model
}

async function handleSubmit() {
  if (!form.video_url || !form.platform || !form.model_name || !form.provider_id) {
    alert('请填写完整信息')
    return
  }
  try {
    await taskStore.submitTask({ ...form })
  } catch (e: any) {
    alert(e.message)
  }
}

async function handleCreateGroup(name: string) {
  try {
    const g = await groupStore.addGroup(name)
    form.group_id = g.id
  } catch {
    // ignore
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
    alert(res.data.message || '连接成功')
  } catch (e: any) {
    alert('连接失败: ' + e.message)
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
    <aside class="sidebar">
      <button class="sidebar-link-btn top" @click="toSettings">
        ⚙ 设置 <span class="link-hint">提供商 · 平台 Cookie</span>
      </button>

      <NoteForm
        :form="form"
        :models="models"
        :selected-model="selectedModel"
        :groups="groupStore.groups"
        :errors="{ error: taskStore.error }"
        :loading="taskStore.loading"
        :testing="testing"
        :has-provider="currentProviderId !== ''"
        @submit="handleSubmit"
        @update:platform="(v: string) => form.platform = v"
        @update:video-url="(v: string) => form.video_url = v"
        @update:style="(v: string) => form.style = v"
        @update:group-id="(v: string) => form.group_id = v"
        @provider-change="handleProviderChange"
        @model-select="handleModelSelect"
        @test-connection="handleTestConnection"
        @create-group="handleCreateGroup"
      />

      <StepBar :status="taskStore.currentStatus" :error-message="taskStore.currentError" />

      <div v-if="taskStore.currentError && isFailedStatus(taskStore.currentStatus)" class="error-banner">
        <span class="error-icon">⚠</span>
        <div>
          <strong>任务执行失败</strong>
          <p>{{ taskStore.currentError }}</p>
        </div>
      </div>
    </aside>

    <section class="content">
      <div v-if="!taskStore.currentMarkdown && !taskStore.loading" class="placeholder">
        <p>输入视频链接，AI 帮你生成结构化笔记</p>
      </div>

      <MarkdownViewer
        v-if="taskStore.currentMarkdown"
        :markdown="taskStore.currentMarkdown"
      />

      <div v-if="taskStore.currentTaskId && taskStore.currentMarkdown" class="actions">
        <a :href="`/api/v1/export/${taskStore.currentTaskId}?format=md`" class="btn-export" download>
          导出笔记
        </a>
      </div>
    </section>

    <aside class="history-panel">
      <div class="history-header">
        <span class="history-title">任务历史</span>
        <span class="history-count">{{ taskStore.tasks.length }}</span>
      </div>
      <TaskHistory
        :tasks="taskStore.tasks"
        @select="handleTaskSelect"
        @delete="(id: string) => taskStore.removeTask(id)"
      />
    </aside>
  </div>
</template>

<style scoped>
.home {
  display: flex;
  height: calc(100vh - 52px);
}
.sidebar {
  width: 320px;
  min-width: 320px;
  background: #fff;
  border-right: 1px solid #e0e0e0;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.sidebar-link-btn {
  width: 100%;
  padding: 8px 12px;
  border: 1px dashed #ccc;
  border-radius: 6px;
  background: #fafafa;
  cursor: pointer;
  font-size: 13px;
  color: #555;
  text-align: left;
  transition: all 0.2s;
}
.sidebar-link-btn:hover {
  border-color: #4A90D9;
  background: #f0f7ff;
  color: #4A90D9;
}
.sidebar-link-btn.top {
  margin-bottom: 4px;
}
.link-hint {
  font-size: 11px;
  color: #aaa;
  font-weight: normal;
}
.content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  min-width: 0;
}
.history-panel {
  width: 260px;
  min-width: 260px;
  background: #fafafa;
  border-left: 1px solid #e0e0e0;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}
.history-header {
  padding: 12px 14px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: sticky;
  top: 0;
  background: #fafafa;
  z-index: 1;
}
.history-title {
  font-size: 13px;
  font-weight: 600;
  color: #666;
}
.history-count {
  font-size: 12px;
  color: #999;
  background: #e0e0e0;
  padding: 1px 8px;
  border-radius: 10px;
}
.placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #999;
  font-size: 16px;
}
.actions {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}
.btn-export {
  padding: 8px 20px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  text-decoration: none;
  display: inline-block;
}
.btn-export {
  background: #4caf50;
  color: #fff;
}
.error-banner {
  display: flex;
  gap: 10px;
  padding: 12px;
  background: #fff3f3;
  border: 1px solid #f5c6cb;
  border-radius: 6px;
  font-size: 13px;
}
.error-banner strong {
  color: #c0392b;
}
.error-banner p {
  margin: 4px 0 0;
  color: #666;
  word-break: break-all;
}
.error-icon {
  font-size: 18px;
  flex-shrink: 0;
  margin-top: 2px;
}
</style>
