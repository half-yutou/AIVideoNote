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
  <div class="dialog-overlay" @click.self="emit('close')">
    <div class="dialog">
      <h3>{{ provider ? '编辑提供商' : '添加提供商' }}</h3>
      <form @submit.prevent="handleSubmit">
        <label>名称 *</label>
        <input v-model="form.name" placeholder="如 OpenAI" />

        <label>API Key *</label>
        <input v-model="form.api_key" type="password" :placeholder="provider ? '留空则不修改' : 'sk-...'" />

        <label>API 地址 *</label>
        <input v-model="form.base_url" placeholder="https://api.openai.com/v1" />

        <label>类型 *</label>
        <select v-model="form.type">
          <option value="">请选择</option>
          <option value="openai">OpenAI</option>
          <option value="deepseek">DeepSeek</option>
          <option value="qwen">通义千问</option>
          <option value="custom">自定义</option>
        </select>

        <label>Logo URL</label>
        <input v-model="form.logo" placeholder="可选" />

        <label v-if="provider" class="checkbox-label">
          <input type="checkbox" :checked="enabled" @change="enabled = ($event.target as HTMLInputElement).checked" />
          启用该提供商
        </label>

        <div class="dialog-actions">
          <button type="button" class="btn-cancel" @click="emit('close')">取消</button>
          <button type="submit" class="btn-save">保存</button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}
.dialog {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  width: 460px;
  max-width: 90vw;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
}
h3 {
  font-size: 18px;
  margin-bottom: 16px;
}
form {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
label {
  font-size: 13px;
  color: #666;
}
input,
select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  outline: none;
}
input:focus,
select:focus {
  border-color: #1a1a2e;
}
.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  cursor: pointer;
}
.dialog-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  margin-top: 12px;
}
.btn-cancel,
.btn-save {
  padding: 8px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
}
.btn-cancel {
  background: #fff;
  border: 1px solid #ddd;
}
.btn-save {
  background: #1a1a2e;
  color: #fff;
  border: none;
}
.btn-save:hover {
  background: #2a2a4e;
}
</style>
