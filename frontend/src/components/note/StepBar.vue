<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ status: string; errorMessage?: string }>()

const sequentialSteps = [
  { key: 'PENDING', label: '排队中', icon: '⏳' },
  { key: 'PARSING', label: '解析链接', icon: '🔗' },
  { key: 'DOWNLOADING', label: '下载视频', icon: '⬇' },
  { key: 'TRANSCRIBING', label: '语音转录', icon: '🎙' },
  { key: 'GENERATING', label: 'AI 生成', icon: '✨' },
  { key: 'POST_PROCESSING', label: '后处理', icon: '📝' },
]

const stepFailedKeys: Record<string, string> = {
  PARSING_FAILED: 'PARSING',
  DOWNLOADING_FAILED: 'DOWNLOADING',
  TRANSCRIBING_FAILED: 'TRANSCRIBING',
  GENERATING_FAILED: 'GENERATING',
  POST_PROCESSING_FAILED: 'POST_PROCESSING',
}

const isSuccess = computed(() => props.status === 'SUCCESS')
const isFailed = computed(() => props.status === 'FAILED' || props.status?.endsWith('_FAILED'))
const isFinal = computed(() => isSuccess.value || isFailed.value)
const failedStepKey = computed(() => stepFailedKeys[props.status] || '')

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
    <div class="step-track">
      <div
        v-for="(step, i) in sequentialSteps"
        :key="step.key"
        :class="['step', {
          done: i < seqIdx,
          active: i === seqIdx && !isFinal,
          fail: failedStepKey && step.key === failedStepKey
        }]"
      >
        <div class="step-indicator">
          <div class="dot">
            <template v-if="i < seqIdx && step.key !== failedStepKey">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                <polyline points="20 6 9 17 4 12"/>
              </svg>
            </template>
            <template v-else-if="step.key === failedStepKey">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </template>
            <template v-else>
              <span class="dot-num">{{ i + 1 }}</span>
            </template>
          </div>
          <div v-if="i < sequentialSteps.length - 1" :class="['connector', { filled: i < seqIdx - 1 || (i < seqIdx && !failedStepKey) }]"></div>
        </div>
        <span class="step-label">{{ step.label }}</span>
        <div v-if="i === seqIdx && !isFinal" class="active-pulse"></div>
      </div>
    </div>

    <div v-if="isFinal" class="result-badge" :class="{ success: isSuccess, fail: !isSuccess }">
      <svg v-if="isSuccess" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
        <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
        <polyline points="22 4 12 14.01 9 11.01"/>
      </svg>
      <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
        <circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/>
      </svg>
      <span>{{ isSuccess ? '生成完成' : '任务失败' }}</span>
    </div>

    <div v-if="isFailed && errorMessage" class="error-msg">
      {{ errorMessage }}
    </div>
  </div>
</template>

<style scoped>
.step-bar {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-4);
  padding: var(--space-3) var(--space-4);
}

.step-track {
  display: flex;
  align-items: flex-start;
  gap: 0;
  max-width: 600px;
  flex: 1;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  position: relative;
  flex: 1;
}

.step:last-child {
  flex: 0;
}

.step-indicator {
  display: flex;
  align-items: center;
  width: 100%;
}

.dot {
  width: 22px;
  height: 22px;
  border-radius: var(--radius-full);
  background: var(--color-bg-muted);
  border: 2px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all var(--transition-base);
  position: relative;
  z-index: 1;
}
.dot-num {
  font-size: 9px;
  font-weight: 600;
  color: var(--color-text-muted);
}

.connector {
  flex: 1;
  height: 2px;
  background: var(--color-border);
  transition: background var(--transition-base);
}
.connector.filled {
  background: var(--color-success);
}

.step-label {
  font-size: var(--text-xs);
  color: var(--color-text-muted);
  white-space: nowrap;
  transition: color var(--transition-fast);
  margin-top: 4px;
  transform: translateX(calc(-50% + 11px));
}

/* 已完成 */
.step.done .dot {
  background: var(--color-success);
  border-color: var(--color-success);
  color: white;
}
.step.done .step-label {
  color: var(--color-success);
  font-weight: 500;
}

/* 进行中 */
.step.active .dot {
  background: var(--color-primary);
  border-color: var(--color-primary);
  color: white;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}
.step.active .dot-num {
  color: white;
}
.step.active .step-label {
  color: var(--color-primary);
  font-weight: 600;
}

.active-pulse {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 22px;
  height: 22px;
  border-radius: var(--radius-full);
  background: var(--color-primary);
  opacity: 0.2;
  animation: pulse-ring 1.5s ease-out infinite;
}

@keyframes pulse-ring {
  0% { transform: translateX(-50%) scale(1); opacity: 0.3; }
  100% { transform: translateX(-50%) scale(1.8); opacity: 0; }
}

/* 失败 */
.step.fail .dot {
  background: var(--color-danger);
  border-color: var(--color-danger);
  color: white;
}
.step.fail .step-label {
  color: var(--color-danger);
  font-weight: 600;
}

/* 结果标记 */
.result-badge {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-1) var(--space-3);
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  font-weight: 600;
  white-space: nowrap;
  animation: fadeIn 300ms ease-out;
}
.result-badge.success {
  background: var(--color-success-light);
  color: var(--color-success);
}
.result-badge.fail {
  background: var(--color-danger-light);
  color: var(--color-danger);
}

.error-msg {
  padding: var(--space-1) var(--space-3);
  background: var(--color-danger-light);
  border: 1px solid #fecaca;
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  color: var(--color-danger);
  line-height: var(--leading-relaxed);
  word-break: break-all;
  max-width: 300px;
}
</style>
