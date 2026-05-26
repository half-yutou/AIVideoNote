<script setup lang="ts">
import { reactive } from 'vue'

const props = defineProps<{ provider: LLMProviderData | null }>()

const emit = defineEmits<{
  save: [data: ProviderRequest | ProviderUpdateRequest]
  close: []
}>()

const form = reactive<ProviderRequest>({
  name: props.provider?.name || '',
  api_key: '',
  base_url: props.provider?.base_url || '',
  type: props.provider?.type || '',
  logo: props.provider?.logo || '',
})

const enabled = props.provider ? props.provider.enabled : true

function handleSubmit() {
  if (!form.name || !form.base_url || !form.type) return
  if (!props.provider && !form.api_key) return

  const data: any = { ...form }
  if (props.provider) {
    if (!form.api_key) delete data.api_key
    data.enabled = enabled
    delete data.logo
    delete data.type
    Object.keys(data).forEach((k) => {
      if (data[k] === props.provider![k as keyof LLMProviderData]) delete data[k]
    })
  }
  emit('save', data)
}
</script>

<template>
  <Teleport to="body">
    <Transition name="overlay">
      <div class="dialog-overlay" @click.self="emit('close')">
        <Transition name="dialog" appear>
          <div class="dialog">
            <div class="dialog-header">
              <h3>{{ provider ? '编辑提供商' : '添加提供商' }}</h3>
              <button class="close-btn" @click="emit('close')">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
              </button>
            </div>

            <form @submit.prevent="handleSubmit">
              <div class="form-field">
                <label>名称 <span class="required">*</span></label>
                <input v-model="form.name" placeholder="如 OpenAI、DeepSeek" />
              </div>

              <div class="form-field">
                <label>API Key <span class="required">*</span></label>
                <input v-model="form.api_key" type="password" :placeholder="provider ? '留空则不修改' : 'sk-...'" />
              </div>

              <div class="form-field">
                <label>API 地址 <span class="required">*</span></label>
                <input v-model="form.base_url" placeholder="https://api.openai.com/v1" />
              </div>

              <div class="form-field">
                <label>类型 <span class="required">*</span></label>
                <select v-model="form.type">
                  <option value="">请选择类型</option>
                  <option value="openai">OpenAI</option>
                  <option value="deepseek">DeepSeek</option>
                  <option value="qwen">通义千问</option>
                  <option value="custom">自定义</option>
                </select>
              </div>

              <div class="form-field">
                <label>Logo URL</label>
                <input v-model="form.logo" placeholder="可选，提供商图标地址" />
              </div>

              <label v-if="provider" class="checkbox-label">
                <input type="checkbox" :checked="enabled" @change="enabled = ($event.target as HTMLInputElement).checked" />
                <span class="checkbox-text">启用该提供商</span>
              </label>

              <div class="dialog-actions">
                <button type="button" class="btn-cancel" @click="emit('close')">取消</button>
                <button type="submit" class="btn-save">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="20 6 9 17 4 12"/>
                  </svg>
                  保存
                </button>
              </div>
            </form>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.4);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
}

.dialog {
  background: var(--color-bg-elevated);
  border-radius: var(--radius-xl);
  padding: 0;
  width: 480px;
  max-width: 90vw;
  box-shadow: var(--shadow-xl);
  border: 1px solid var(--color-border);
  overflow: hidden;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-5) var(--space-6);
  border-bottom: 1px solid var(--color-border-light);
}
.dialog-header h3 {
  font-size: var(--text-lg);
  font-weight: 600;
  color: var(--color-text);
}
.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: var(--radius-md);
  color: var(--color-text-muted);
  transition: all var(--transition-fast);
}
.close-btn:hover {
  background: var(--color-bg-muted);
  color: var(--color-text);
}

form {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  padding: var(--space-5) var(--space-6);
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}
.form-field label {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-secondary);
}
.required {
  color: var(--color-danger);
}

.form-field input,
.form-field select {
  padding: 10px 14px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--text-base);
  background: var(--color-bg-elevated);
  transition: all var(--transition-fast);
  outline: none;
}
.form-field input:focus,
.form-field select:focus {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-focus);
}
.form-field select {
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%2364748b' d='M2.5 4.5L6 8l3.5-3.5'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  padding-right: 32px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  cursor: pointer;
  padding: var(--space-2) 0;
}
.checkbox-label input[type="checkbox"] {
  width: 16px;
  height: 16px;
  accent-color: var(--color-primary);
  cursor: pointer;
}
.checkbox-text {
  font-size: var(--text-base);
  color: var(--color-text);
}

.dialog-actions {
  display: flex;
  gap: var(--space-3);
  justify-content: flex-end;
  padding-top: var(--space-3);
  border-top: 1px solid var(--color-border-light);
}

.btn-cancel,
.btn-save {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-5);
  border-radius: var(--radius-md);
  font-size: var(--text-base);
  font-weight: 500;
  transition: all var(--transition-fast);
}
.btn-cancel {
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  color: var(--color-text-secondary);
}
.btn-cancel:hover {
  background: var(--color-bg-muted);
  border-color: var(--color-text-muted);
}
.btn-save {
  background: var(--color-primary);
  color: var(--color-text-inverse);
  border: none;
}
.btn-save:hover {
  background: var(--color-primary-hover);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

/* 过渡动画 */
.overlay-enter-active {
  transition: opacity 200ms ease-out;
}
.overlay-leave-active {
  transition: opacity 150ms ease-in;
}
.overlay-enter-from,
.overlay-leave-to {
  opacity: 0;
}

.dialog-enter-active {
  transition: all 300ms cubic-bezier(0.34, 1.56, 0.64, 1);
}
.dialog-leave-active {
  transition: all 150ms ease-in;
}
.dialog-enter-from {
  opacity: 0;
  transform: scale(0.95) translateY(10px);
}
.dialog-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
