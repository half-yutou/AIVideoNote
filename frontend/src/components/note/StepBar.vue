<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ status: string; errorMessage?: string }>()

const sequentialSteps = [
  { key: 'PENDING', label: '排队中' },
  { key: 'PARSING', label: '解析链接' },
  { key: 'DOWNLOADING', label: '下载视频' },
  { key: 'TRANSCRIBING', label: '语音转录' },
  { key: 'GENERATING', label: 'AI 生成' },
  { key: 'POST_PROCESSING', label: '后处理' },
]

const stepFailedKeys = {
  PARSING_FAILED: 'PARSING',
  DOWNLOADING_FAILED: 'DOWNLOADING',
  TRANSCRIBING_FAILED: 'TRANSCRIBING',
  GENERATING_FAILED: 'GENERATING',
  POST_PROCESSING_FAILED: 'POST_PROCESSING',
}

const isSuccess = computed(() => props.status === 'SUCCESS')
const isFailed = computed(() => props.status === 'FAILED' || props.status?.endsWith('_FAILED'))
const isFinal = computed(() => isSuccess.value || isFailed.value)
const failedStepKey = computed(() => (stepFailedKeys as Record<string, string>)[props.status] || '')

const seqIdx = computed(() => {
  const idx = sequentialSteps.findIndex((s) => s.key === props.status)
  if (idx >= 0) return idx
  if (failedStepKey.value) return sequentialSteps.findIndex((s) => s.key === failedStepKey.value) + 1
  if (isSuccess.value) return sequentialSteps.length
  return -1
})
</script>

<template>
  <div v-if="status" class="step-bar">
    <div v-for="(step, i) in sequentialSteps" :key="step.key"
      :class="['step', {
        done: i < seqIdx,
        active: i === seqIdx && !isFinal,
        fail: i === seqIdx - 1 && failedStepKey && step.key === failedStepKey
      }]">
      <div class="dot">{{ i < seqIdx && step.key !== failedStepKey ? '✓' : step.key === failedStepKey ? '✗' : i + 1 }}</div>
      <span class="label">{{ step.label }}</span>
    </div>

    <div v-if="isFinal" class="branch">
      <div :class="['step', 'result', { done: isSuccess, fail: !isSuccess }]">
        <div class="dot">{{ isSuccess ? '✓' : '✗' }}</div>
        <span class="label">{{ isSuccess ? '成功' : '失败' }}</span>
      </div>
    </div>

    <div v-if="isFailed && errorMessage" class="error-msg">
      {{ errorMessage }}
    </div>
  </div>
</template>

<style scoped>
.step-bar {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 8px 0;
}
.step {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: #ccc;
}
.step.done { color: #4caf50; }
.step.active { color: #1a1a2e; font-weight: 600; }
.step.fail { color: #e53935; font-weight: 600; }
.step.result { color: #4caf50; font-weight: 600; }
.dot {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #eee;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  color: #999;
  flex-shrink: 0;
}
.step.done .dot { background: #4caf50; color: #fff; }
.step.active .dot { background: #1a1a2e; color: #fff; }
.step.fail .dot { background: #e53935; color: #fff; }
.step.result .dot { background: #4caf50; color: #fff; }
.step.result.fail .dot { background: #e53935; color: #fff; }
.label { white-space: nowrap; }
.branch { margin-top: 2px; }
.error-msg {
  margin-top: 6px;
  padding: 8px 10px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 6px;
  font-size: 12px;
  color: #dc2626;
  line-height: 1.5;
  word-break: break-all;
}
</style>
