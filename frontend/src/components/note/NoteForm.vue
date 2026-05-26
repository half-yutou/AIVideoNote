<script setup lang="ts">
import { ref } from 'vue'
import { useProviderStore } from '@/stores/provider'
import { ALLOWED_MEDIA_EXTENSIONS, MAX_UPLOAD_SIZE, validateFile } from '@/api/upload'

const props = defineProps<{
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
  'update:local-file': [f: File | null]
  'provider-change': [id: string]
  'model-select': [model: string]
  'test-connection': []
  'create-group': [name: string]
}>()

const providerStore = useProviderStore()

const platforms = [
  { value: 'bilibili', label: 'B站', icon: 'bilibili' },
  { value: 'local', label: '本地文件', icon: 'local' },
]

const styles = [
  { value: '', label: '默认' },
  { value: 'academic', label: '学术风' },
  { value: 'casual', label: '口语风' },
  { value: 'keypoint', label: '要点提取' },
]

const newGroupName = ref('')
const showNewGroup = ref(false)

// 本地文件上传相关
const localFile = ref<File | null>(null)
const dragOver = ref(false)
const fileError = ref('')
const fileInputRef = ref<HTMLInputElement | null>(null)

const acceptExts = ALLOWED_MEDIA_EXTENSIONS.join(',')

function handleFileDrop(e: DragEvent) {
  dragOver.value = false
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    selectFile(files[0])
  }
}

function handleFileSelect(e: Event) {
  const input = e.target as HTMLInputElement
  if (input.files && input.files.length > 0) {
    selectFile(input.files[0])
  }
}

function selectFile(file: File) {
  fileError.value = ''
  const err = validateFile(file)
  if (err) {
    fileError.value = err
    localFile.value = null
    emit('update:local-file', null)
    return
  }
  localFile.value = file
  emit('update:local-file', file)
}

function removeFile() {
  localFile.value = null
  fileError.value = ''
  emit('update:local-file', null)
  if (fileInputRef.value) {
    fileInputRef.value.value = ''
  }
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
}

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
    <div class="form-header">
      <h3>生成视频笔记</h3>
      <span class="form-badge">AI</span>
    </div>

    <!-- 平台选择（放在输入区域上方） -->
    <div class="form-group">
      <label class="form-label">平台</label>
      <div class="platform-tabs">
        <button
          v-for="p in platforms"
          :key="p.value"
          type="button"
          :class="['plat-btn', { active: form.platform === p.value }]"
          @click="emit('update:platform', p.value)"
        >
          <img v-if="p.icon === 'bilibili'" src="/bilibili-color.svg" class="plat-icon-svg" alt="bilibili" />
          <svg v-else-if="p.icon === 'local'" class="plat-icon-svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
          </svg>
          {{ p.label }}
        </button>
      </div>
    </div>

    <!-- 输入区域容器（固定高度防止跳变） -->
    <div class="input-area-container">
      <!-- B站链接输入 -->
      <div v-if="form.platform !== 'local'" class="form-group">
        <label class="form-label">视频链接</label>
        <div class="input-wrapper">
          <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/>
            <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>
          </svg>
          <input
            type="text"
            class="form-input has-icon"
            :value="form.video_url"
            placeholder="粘贴 B站视频链接"
            @input="emit('update:video-url', ($event.target as HTMLInputElement).value)"
          />
        </div>
      </div>

      <!-- 本地文件上传区域 -->
      <div v-else class="form-group">
        <label class="form-label">选择文件</label>
        <div
          v-if="!localFile"
          class="upload-zone"
          :class="{ 'drag-over': dragOver }"
          @dragover.prevent="dragOver = true"
          @dragleave.prevent="dragOver = false"
          @drop.prevent="handleFileDrop"
          @click="fileInputRef?.click()"
        >
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="17 8 12 3 7 8"/>
            <line x1="12" y1="3" x2="12" y2="15"/>
          </svg>
          <div class="upload-text">
            <p class="upload-hint">拖拽文件到此处，或 <span class="upload-link">点击选择</span></p>
            <p class="upload-formats">支持 MP4/MKV/AVI/MP3/WAV 等，最大 2GB</p>
          </div>
          <input
            ref="fileInputRef"
            type="file"
            :accept="acceptExts"
            class="upload-input-hidden"
            @change="handleFileSelect"
          />
        </div>
        <div v-else class="upload-file-info">
          <div class="file-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polygon points="23 7 16 12 23 17 23 7"/>
              <rect x="1" y="5" width="15" height="14" rx="2" ry="2"/>
            </svg>
          </div>
          <div class="file-details">
            <span class="file-name">{{ localFile.name }}</span>
            <span class="file-size">{{ formatFileSize(localFile.size) }}</span>
          </div>
          <button type="button" class="file-remove" @click="removeFile" title="移除文件">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <p v-if="fileError" class="form-error">{{ fileError }}</p>
      </div>
    </div>

    <div class="form-group">
      <label class="form-label">LLM 提供商</label>
      <select
        :value="form.provider_id"
        class="form-select"
        @change="emit('provider-change', ($event.target as HTMLSelectElement).value)"
      >
        <option value="">请选择提供商</option>
        <option
          v-for="p in providerStore.providers.filter((x) => x.enabled)"
          :key="p.id"
          :value="p.id"
        >
          {{ p.name }} ({{ p.type }})
        </option>
      </select>
      <div v-if="hasProvider" class="form-row-actions">
        <button
          type="button"
          class="test-btn"
          :disabled="testing"
          @click="emit('test-connection')"
        >
          <svg v-if="!testing" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
            <polyline points="22 4 12 14.01 9 11.01"/>
          </svg>
          <span v-if="testing" class="btn-spinner"></span>
          {{ testing ? '测试中...' : '测试连通性' }}
        </button>
      </div>
    </div>

    <div class="form-group">
      <label class="form-label">模型</label>
      <select
        :value="selectedModel"
        class="form-select"
        :disabled="models.length === 0"
        @change="emit('model-select', ($event.target as HTMLSelectElement).value)"
      >
        <option value="">{{ models.length === 0 ? '请先选择提供商' : '请选择模型' }}</option>
        <option v-for="m in models" :key="m" :value="m">{{ m }}</option>
      </select>
    </div>

    <div class="form-row">
      <div class="form-group flex-1">
        <label class="form-label">笔记风格</label>
        <select
          :value="form.style"
          class="form-select"
          @change="emit('update:style', ($event.target as HTMLSelectElement).value)"
        >
          <option v-for="s in styles" :key="s.value" :value="s.value">{{ s.label }}</option>
        </select>
      </div>

      <div class="form-group flex-1">
        <label class="form-label">分组</label>
        <select
          :value="form.group_id || 'default'"
          class="form-select"
          @change="emit('update:group-id', ($event.target as HTMLSelectElement).value)"
        >
          <option v-for="g in groups" :key="g.id" :value="g.id">{{ g.name }}</option>
        </select>
      </div>
    </div>

    <div class="form-row-actions">
      <button type="button" class="link-btn" @click="showNewGroup = true">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        新建分组
      </button>
    </div>

    <Transition name="slide">
      <div v-if="showNewGroup" class="new-group-row">
        <input
          type="text"
          v-model="newGroupName"
          placeholder="输入分组名称"
          class="form-input new-group-input"
          @keydown.enter="handleCreateGroup()"
        />
        <button type="button" class="confirm-btn" @click="handleCreateGroup()">确认</button>
        <button type="button" class="cancel-btn" @click="showNewGroup = false">取消</button>
      </div>
    </Transition>

    <p v-if="errors.error" class="form-error">{{ errors.error }}</p>

    <button type="submit" class="submit-btn" :disabled="loading">
      <svg v-if="!loading" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
      </svg>
      <span v-if="loading" class="btn-spinner"></span>
      {{ loading ? '处理中...' : '生成笔记' }}
    </button>
  </form>
</template>

<style scoped>
.note-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.form-header {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-bottom: var(--space-1);
}
.form-header h3 {
  font-size: var(--text-md);
  font-weight: 600;
  color: var(--color-text);
}
.form-badge {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 6px;
  background: linear-gradient(135deg, var(--color-primary), #8b5cf6);
  color: white;
  border-radius: var(--radius-sm);
  letter-spacing: 0.5px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.form-label {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-secondary);
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}
.input-icon {
  position: absolute;
  left: 12px;
  color: var(--color-text-muted);
  pointer-events: none;
}
.form-input {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--text-base);
  background: var(--color-bg-elevated);
  transition: all var(--transition-fast);
}
.form-input.has-icon {
  padding-left: 36px;
}
.form-input:focus {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-focus);
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
}
.form-select:disabled {
  background: var(--color-bg-muted);
  color: var(--color-text-muted);
  cursor: not-allowed;
}

/* 输入区域容器 - 固定高度防止切换平台时跳变 */
.input-area-container {
  min-height: 72px;
  display: flex;
  flex-direction: column;
}

.platform-tabs {
  display: flex;
  gap: var(--space-2);
}
.plat-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg-elevated);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
}
.plat-btn:hover {
  border-color: var(--color-primary-subtle);
  background: var(--color-primary-light);
}
.plat-btn.active {
  background: var(--color-primary);
  color: var(--color-text-inverse);
  border-color: var(--color-primary);
  box-shadow: var(--shadow-sm);
}
.plat-icon-svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.form-row {
  display: flex;
  gap: var(--space-3);
}
.flex-1 {
  flex: 1;
}

.form-row-actions {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.link-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  font-size: var(--text-sm);
  color: var(--color-primary);
  padding: 0;
  transition: opacity var(--transition-fast);
}
.link-btn:hover {
  opacity: 0.7;
}

.test-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  font-size: var(--text-xs);
  color: var(--color-success);
  border: 1px solid var(--color-success);
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  transition: all var(--transition-fast);
}
.test-btn:hover:not(:disabled) {
  background: var(--color-success);
  color: var(--color-text-inverse);
}
.test-btn:disabled {
  opacity: 0.5;
}

.form-error {
  color: var(--color-danger);
  font-size: var(--text-sm);
  padding: var(--space-2) var(--space-3);
  background: var(--color-danger-light);
  border-radius: var(--radius-sm);
}

.submit-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  margin-top: var(--space-2);
  padding: 12px;
  background: var(--color-primary);
  color: var(--color-text-inverse);
  border-radius: var(--radius-md);
  font-size: var(--text-base);
  font-weight: 600;
  transition: all var(--transition-fast);
}
.submit-btn:hover:not(:disabled) {
  background: var(--color-primary-hover);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}
.submit-btn:active:not(:disabled) {
  transform: translateY(0);
}
.submit-btn:disabled {
  opacity: 0.6;
}

.btn-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

.new-group-row {
  display: flex;
  gap: var(--space-2);
  align-items: center;
}
.new-group-input {
  flex: 1;
  padding: 8px 12px;
  font-size: var(--text-sm);
}
.confirm-btn {
  padding: 6px 14px;
  background: var(--color-primary);
  color: var(--color-text-inverse);
  border-radius: var(--radius-sm);
  font-size: var(--text-xs);
  font-weight: 500;
  transition: background var(--transition-fast);
}
.confirm-btn:hover {
  background: var(--color-primary-hover);
}
.cancel-btn {
  padding: 6px 14px;
  background: var(--color-bg-muted);
  color: var(--color-text-secondary);
  border-radius: var(--radius-sm);
  font-size: var(--text-xs);
  transition: background var(--transition-fast);
}
.cancel-btn:hover {
  background: var(--color-bg-subtle);
}

/* 过渡动画 */
.slide-enter-active {
  transition: all 200ms ease-out;
}
.slide-leave-active {
  transition: all 150ms ease-in;
}
.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

/* 文件上传区域 */
.upload-zone {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  border: 2px dashed var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg-elevated);
  cursor: pointer;
  transition: all var(--transition-fast);
}
.upload-zone svg {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
}
.upload-zone:hover {
  border-color: var(--color-primary-subtle);
  background: var(--color-primary-light);
}
.upload-zone.drag-over {
  border-color: var(--color-primary);
  background: var(--color-primary-light);
  box-shadow: var(--shadow-focus);
}
.upload-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.upload-hint {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  margin: 0;
}
.upload-link {
  color: var(--color-primary);
  font-weight: 500;
}
.upload-formats {
  font-size: var(--text-xs);
  color: var(--color-text-muted);
  margin: 0;
  line-height: 1.4;
}
.upload-input-hidden {
  display: none;
}

.upload-file-info {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg-elevated);
}
.file-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: var(--color-primary-light);
  border-radius: var(--radius-sm);
  color: var(--color-primary);
  flex-shrink: 0;
}
.file-details {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.file-name {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.file-size {
  font-size: var(--text-xs);
  color: var(--color-text-muted);
}
.file-remove {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: var(--radius-sm);
  color: var(--color-text-muted);
  transition: all var(--transition-fast);
  flex-shrink: 0;
}
.file-remove:hover {
  background: var(--color-danger-light);
  color: var(--color-danger);
}
</style>
