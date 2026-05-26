<script setup lang="ts">
defineProps<{ providers: LLMProviderData[]; loading: boolean }>()

defineEmits<{
  edit: [provider: LLMProviderData]
  delete: [id: string]
}>()
</script>

<template>
  <div class="provider-list">
    <div v-if="loading" class="loading">
      <div class="loading-spinner"></div>
      <span>加载中...</span>
    </div>
    <div v-else-if="providers.length === 0" class="empty">
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
        <line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/>
      </svg>
      <span>暂无提供商，点击上方按钮添加</span>
    </div>
    <div
      v-for="p in providers"
      :key="p.id"
      class="provider-item"
    >
      <div class="provider-main">
        <div class="provider-header">
          <span class="provider-name">{{ p.name }}</span>
          <span class="provider-type">{{ p.type }}</span>
          <span :class="['status-badge', p.enabled ? 'on' : 'off']">
            {{ p.enabled ? '启用' : '禁用' }}
          </span>
        </div>
        <span class="provider-url">{{ p.base_url }}</span>
      </div>
      <div class="provider-actions">
        <button class="action-btn" @click="$emit('edit', p)">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
          </svg>
          编辑
        </button>
        <button class="action-btn danger" @click="$emit('delete', p.id)">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3 6 5 6 21 6"/>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
          </svg>
          删除
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.provider-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-6);
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}
.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid var(--color-bg-subtle);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-6);
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}

.provider-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  transition: all var(--transition-fast);
}
.provider-item:hover {
  border-color: var(--color-primary-subtle);
  box-shadow: var(--shadow-xs);
}

.provider-main {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  flex: 1;
  min-width: 0;
}

.provider-header {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.provider-name {
  font-weight: 600;
  font-size: var(--text-base);
  color: var(--color-text);
}

.provider-type {
  background: var(--color-bg-muted);
  padding: 2px 8px;
  border-radius: var(--radius-sm);
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
  font-weight: 500;
}

.provider-url {
  font-size: var(--text-xs);
  color: var(--color-text-muted);
  font-family: var(--font-mono);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-badge {
  padding: 2px 8px;
  border-radius: var(--radius-full);
  font-size: 11px;
  font-weight: 500;
}
.status-badge.on {
  background: var(--color-success-light);
  color: var(--color-success);
}
.status-badge.off {
  background: var(--color-danger-light);
  color: var(--color-danger);
}

.provider-actions {
  display: flex;
  gap: var(--space-2);
  flex-shrink: 0;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-1) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-elevated);
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
}
.action-btn:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
  background: var(--color-primary-light);
}
.action-btn.danger:hover {
  border-color: var(--color-danger);
  color: var(--color-danger);
  background: var(--color-danger-light);
}
</style>
