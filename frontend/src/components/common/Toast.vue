<script setup lang="ts">
import { useToast } from '@/composables/useToast'

const { toasts, remove } = useToast()

function iconForType(type: string) {
  switch (type) {
    case 'success': return '✓'
    case 'error': return '✕'
    case 'warning': return '!'
    case 'info': return 'i'
    default: return 'i'
  }
}
</script>

<template>
  <Teleport to="body">
    <div class="toast-container">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          :class="['toast-item', `toast-${toast.type}`]"
          @click="remove(toast.id)"
        >
          <span class="toast-icon">{{ iconForType(toast.type) }}</span>
          <span class="toast-message">{{ toast.message }}</span>
          <button class="toast-close" @click.stop="remove(toast.id)">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-container {
  position: fixed;
  top: calc(var(--header-height, 56px) + 16px);
  right: 20px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 10px;
  pointer-events: none;
}

.toast-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: var(--radius-lg);
  background: var(--color-bg-elevated);
  box-shadow: var(--shadow-lg);
  border: 1px solid var(--color-border);
  min-width: 280px;
  max-width: 420px;
  pointer-events: auto;
  cursor: pointer;
  transition: all var(--transition-base);
}
.toast-item:hover {
  box-shadow: var(--shadow-xl);
  transform: translateX(-2px);
}

.toast-icon {
  width: 24px;
  height: 24px;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
}

.toast-success .toast-icon {
  background: var(--color-success-light);
  color: var(--color-success);
}
.toast-error .toast-icon {
  background: var(--color-danger-light);
  color: var(--color-danger);
}
.toast-warning .toast-icon {
  background: var(--color-warning-light);
  color: var(--color-warning);
}
.toast-info .toast-icon {
  background: var(--color-info-light);
  color: var(--color-info);
}

.toast-message {
  flex: 1;
  font-size: var(--text-base);
  color: var(--color-text);
  line-height: var(--leading-tight);
  word-break: break-word;
}

.toast-close {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: var(--radius-full);
  color: var(--color-text-muted);
  transition: all var(--transition-fast);
}
.toast-close:hover {
  background: var(--color-bg-muted);
  color: var(--color-text);
}

/* 过渡动画 */
.toast-enter-active {
  transition: all 300ms ease-out;
}
.toast-leave-active {
  transition: all 200ms ease-in;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(30px);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(30px) scale(0.95);
}
.toast-move {
  transition: transform 200ms ease;
}
</style>
