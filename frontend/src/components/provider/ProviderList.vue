<script setup lang="ts">
defineProps<{ providers: LLMProviderData[]; loading: boolean }>()

defineEmits<{
  edit: [provider: LLMProviderData]
  delete: [id: string]
}>()
</script>

<template>
  <div class="provider-list">
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="providers.length === 0" class="empty">暂无提供商，请添加</div>
    <div
      v-for="p in providers"
      :key="p.id"
      class="provider-item"
    >
      <div class="provider-info">
        <span class="provider-name">{{ p.name }}</span>
        <span class="provider-type">{{ p.type }}</span>
        <span class="provider-url">{{ p.base_url }}</span>
        <span :class="['status-badge', p.enabled ? 'on' : 'off']">
          {{ p.enabled ? '启用' : '禁用' }}
        </span>
      </div>
      <div class="provider-actions">
        <button class="action-btn edit" @click="$emit('edit', p)">编辑</button>
        <button class="action-btn del" @click="$emit('delete', p.id)">删除</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.provider-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.loading,
.empty {
  text-align: center;
  color: #999;
  padding: 20px;
  font-size: 14px;
}
.provider-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
}
.provider-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}
.provider-name {
  font-weight: 600;
  font-size: 14px;
}
.provider-type {
  background: #e8e8e8;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  color: #666;
}
.provider-url {
  font-size: 12px;
  color: #999;
}
.status-badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
}
.status-badge.on { background: #e8f5e9; color: #2e7d32; }
.status-badge.off { background: #fbe9e7; color: #c62828; }
.provider-actions {
  display: flex;
  gap: 8px;
}
.action-btn {
  padding: 4px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
  cursor: pointer;
  font-size: 12px;
}
.action-btn.edit:hover { border-color: #1a1a2e; color: #1a1a2e; }
.action-btn.del:hover { border-color: #e53935; color: #e53935; }
</style>
