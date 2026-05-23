<script setup lang="ts">
import { ref } from 'vue'
import { useProviderStore } from '@/stores/provider'

defineProps<{
  form: GenerateRequest
  models: string[]
  selectedModel: string
  groups: GroupData[]
  errors: { error: string }
  loading: boolean
  testing: boolean
  hasProvider: boolean
}>()

const emit = defineEmits<{
  submit: []
  'update:platform': [v: string]
  'update:video-url': [v: string]
  'update:style': [v: string]
  'update:group-id': [v: string]
  'provider-change': [id: string]
  'model-select': [model: string]
  'test-connection': []
  'create-group': [name: string]
}>()

const providerStore = useProviderStore()

const platforms = [
  { value: 'bilibili', label: 'B站' },
  { value: 'youtube', label: 'YouTube' },
  { value: 'douyin', label: '抖音' },
  { value: 'kuaishou', label: '快手' },
  { value: 'local', label: '本地文件' },
]

const styles = [
  { value: '', label: '默认' },
  { value: 'academic', label: '学术风' },
  { value: 'casual', label: '口语风' },
  { value: 'keypoint', label: '要点提取' },
]

const newGroupName = ref('')
const showNewGroup = ref(false)

function handleCreateGroup() {
  const name = newGroupName.value.trim()
  if (!name) return
  emit('create-group', name)
  newGroupName.value = ''
  showNewGroup.value = false
}
</script>

<template>
  <form class="note-form" @submit.prevent="emit('submit')">
    <h3>生成视频笔记</h3>

    <label>视频链接</label>
    <input
      type="text"
      :value="form.video_url"
      placeholder="粘贴 B站/YouTube/抖音/快手 链接"
      @input="emit('update:video-url', ($event.target as HTMLInputElement).value)"
    />

    <label>平台</label>
    <div class="platform-tabs">
      <button
        v-for="p in platforms"
        :key="p.value"
        type="button"
        :class="['plat-btn', { active: form.platform === p.value }]"
        @click="emit('update:platform', p.value)"
      >
        {{ p.label }}
      </button>
    </div>

    <label>LLM 提供商</label>
    <select
      :value="form.provider_id"
      class="select"
      @change="emit('provider-change', ($event.target as HTMLSelectElement).value)"
    >
      <option value="">请选择</option>
      <option
        v-for="p in providerStore.providers.filter((x) => x.enabled)"
        :key="p.id"
        :value="p.id"
      >
        {{ p.name }} ({{ p.type }})
      </option>
    </select>
    <div class="row-actions">
      <button
        v-if="hasProvider"
        type="button"
        class="test-btn"
        :disabled="testing"
        @click="emit('test-connection')"
      >
        {{ testing ? '测试中...' : '测试连通性' }}
      </button>
    </div>

    <label>模型</label>
    <select
      :value="selectedModel"
      class="select"
      :disabled="models.length === 0"
      @change="emit('model-select', ($event.target as HTMLSelectElement).value)"
    >
      <option value="">{{ models.length === 0 ? '请先选择提供商' : '请选择模型' }}</option>
      <option v-for="m in models" :key="m" :value="m">{{ m }}</option>
    </select>

    <label>笔记风格</label>
    <select
      :value="form.style"
      class="select"
      @change="emit('update:style', ($event.target as HTMLSelectElement).value)"
    >
      <option v-for="s in styles" :key="s.value" :value="s.value">{{ s.label }}</option>
    </select>

    <label>分组</label>
    <select
      :value="form.group_id || 'default'"
      class="select"
      @change="emit('update:group-id', ($event.target as HTMLSelectElement).value)"
    >
      <option v-for="g in groups" :key="g.id" :value="g.id">{{ g.name }}</option>
    </select>
    <div class="row-actions">
      <button type="button" class="link-btn" @click="showNewGroup = true">
        + 新建分组
      </button>
    </div>
    <div v-if="showNewGroup" class="new-group-row">
      <input
        type="text"
        v-model="newGroupName"
        placeholder="输入分组名称"
        class="new-group-input"
        @keydown.enter="handleCreateGroup()"
      />
      <button type="button" class="confirm-btn" @click="handleCreateGroup()">确认</button>
      <button type="button" class="cancel-btn" @click="showNewGroup = false">取消</button>
    </div>

    <p v-if="errors.error" class="error">{{ errors.error }}</p>

    <button type="submit" class="submit-btn" :disabled="loading">
      {{ loading ? '提交中...' : '生成笔记' }}
    </button>
  </form>
</template>

<style scoped>
.note-form {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
h3 {
  font-size: 16px;
  margin-bottom: 4px;
}
label {
  font-size: 13px;
  color: #666;
  margin-top: 4px;
}
input[type='text'] {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  outline: none;
}
input[type='text']:focus {
  border-color: #1a1a2e;
}
.select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  background: #fff;
  outline: none;
}
.platform-tabs {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}
.plat-btn {
  padding: 6px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: #fff;
  cursor: pointer;
  font-size: 13px;
  color: #666;
}
.plat-btn.active {
  background: #1a1a2e;
  color: #fff;
  border-color: #1a1a2e;
}
.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  cursor: pointer;
}
.row-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}
.link-btn {
  font-size: 13px;
  color: #1a1a2e;
  background: none;
  border: none;
  cursor: pointer;
  text-align: left;
  padding: 0;
}
.link-btn:hover {
  text-decoration: underline;
}
.test-btn {
  font-size: 12px;
  color: #4caf50;
  background: none;
  border: 1px solid #4caf50;
  padding: 3px 10px;
  border-radius: 4px;
  cursor: pointer;
}
.test-btn:hover:not(:disabled) {
  background: #4caf50;
  color: #fff;
}
.test-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.error {
  color: #e53935;
  font-size: 13px;
}
.submit-btn {
  margin-top: 8px;
  padding: 10px;
  background: #1a1a2e;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  font-weight: 600;
}
.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.submit-btn:hover:not(:disabled) {
  background: #2a2a4e;
}
.new-group-row {
  display: flex;
  gap: 6px;
  align-items: center;
}
.new-group-input {
  flex: 1;
  padding: 6px 10px;
  border: 1px solid #4a90d9;
  border-radius: 4px;
  font-size: 13px;
  outline: none;
}
.confirm-btn {
  padding: 4px 12px;
  background: #1a1a2e;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}
.cancel-btn {
  padding: 4px 12px;
  background: #eee;
  color: #666;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}
</style>
