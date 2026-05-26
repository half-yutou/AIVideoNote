<script setup lang="ts">
import { computed, ref } from 'vue'
import { useTaskStore } from '@/stores/task'
import { useGroupStore } from '@/stores/group'

const props = defineProps<{ tasks: TaskItemData[] }>()
const emit = defineEmits<{
  select: [id: string]
  delete: [id: string]
}>()

const taskStore = useTaskStore()
const groupStore = useGroupStore()

const collapsedGroups = ref<Record<string, boolean>>({})
const editingTaskId = ref<string | null>(null)
const editingGroupName = ref<string | null>(null)
const editValue = ref('')

const statusLabel: Record<string, string> = {
  PENDING: '排队',
  PARSING: '解析',
  DOWNLOADING: '下载',
  TRANSCRIBING: '转录',
  GENERATING: '生成',
  POST_PROCESSING: '处理',
  PARSING_FAILED: '解析失败',
  DOWNLOADING_FAILED: '下载失败',
  TRANSCRIBING_FAILED: '转录失败',
  GENERATING_FAILED: '生成失败',
  POST_PROCESSING_FAILED: '处理失败',
  SUCCESS: '完成',
  FAILED: '失败',
}

function statusClass(s: string) {
  if (s === 'SUCCESS' || s === 'POST_PROCESSING') return 's-success'
  if (s === 'FAILED' || s?.endsWith('_FAILED')) return 's-failed'
  return 's-running'
}

function taskLabel(task: TaskItemData) {
  return task.name || task.video_url.slice(0, 40)
}

interface TaskGroup {
  name: string
  tasks: TaskItemData[]
}

const groups = computed<TaskGroup[]>(() => {
  const map = new Map<string, TaskItemData[]>()
  for (const t of props.tasks) {
    const key = t.group_name || '默认'
    if (!map.has(key)) map.set(key, [])
    map.get(key)!.push(t)
  }
  const result: TaskGroup[] = []
  for (const [name, tasks] of map) {
    result.push({ name, tasks })
  }
  result.sort((a, b) => {
    if (a.name === '默认') return -1
    if (b.name === '默认') return 1
    return a.name.localeCompare(b.name)
  })
  return result
})

function toggleGroup(name: string) {
  collapsedGroups.value[name] = !collapsedGroups.value[name]
}

function startEditTask(task: TaskItemData) {
  editingTaskId.value = task.id
  editingGroupName.value = null
  editValue.value = task.name || ''
}

function startEditGroup(groupName: string) {
  editingGroupName.value = groupName
  editingTaskId.value = null
  editValue.value = groupName === '默认' ? '' : groupName
}

async function confirmEditTask(task: TaskItemData) {
  const name = editValue.value.trim()
  if (name && name !== task.name) {
    await taskStore.renameTask(task.id, name)
  }
  editingTaskId.value = null
}

async function confirmEditGroup(oldName: string) {
  const newName = editValue.value.trim()
  if (!newName || newName === oldName) {
    editingGroupName.value = null
    return
  }
  const g = groupStore.groups.find((x) => x.name === oldName)
  if (g) {
    await groupStore.updateGroupName(g.id, newName)
  }
  editingGroupName.value = null
}

function cancelEdit() {
  editingTaskId.value = null
  editingGroupName.value = null
}

function handleEditKeydown(e: KeyboardEvent, confirmFn: () => void) {
  if (e.key === 'Enter') confirmFn()
  if (e.key === 'Escape') cancelEdit()
}
</script>

<template>
  <div class="task-history">
    <div v-if="tasks.length === 0" class="empty">
      <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
        <line x1="3" y1="9" x2="21" y2="9"/>
        <line x1="9" y1="21" x2="9" y2="9"/>
      </svg>
      <span>暂无任务记录</span>
    </div>

    <div v-for="group in groups" :key="group.name" class="group">
      <div class="group-header" @click="toggleGroup(group.name)">
        <svg
          :class="['group-arrow', { collapsed: collapsedGroups[group.name] }]"
          width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
        >
          <polyline points="6 9 12 15 18 9"/>
        </svg>
        <template v-if="editingGroupName === group.name">
          <input
            class="edit-input"
            v-model="editValue"
            @click.stop
            @keydown="handleEditKeydown($event, () => confirmEditGroup(group.name))"
            @blur="confirmEditGroup(group.name)"
            autofocus
          />
        </template>
        <span v-else class="group-name">{{ group.name }}</span>
        <span class="group-count">{{ group.tasks.length }}</span>
        <button
          class="icon-btn"
          @click.stop="startEditGroup(group.name)"
          title="重命名分组"
        >
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
          </svg>
        </button>
      </div>

      <TransitionGroup name="list" tag="div" v-show="!collapsedGroups[group.name]" class="group-tasks">
        <div
          v-for="task in group.tasks"
          :key="task.id"
          class="task-item"
          @click="$emit('select', task.id)"
        >
          <span :class="['status-dot', statusClass(task.status)]" />
          <div class="task-info">
            <template v-if="editingTaskId === task.id">
              <input
                class="edit-input"
                v-model="editValue"
                @click.stop
                @keydown="handleEditKeydown($event, () => confirmEditTask(task))"
                @blur="confirmEditTask(task)"
                autofocus
              />
            </template>
            <template v-else>
              <span class="task-name">{{ taskLabel(task) }}</span>
              <span class="task-meta">
                <span class="platform-badge">
                  <img v-if="task.platform === 'bilibili'" src="/bilibili-color.svg" class="platform-logo" alt="B站" />
                  <template v-else>{{ task.platform }}</template>
                </span>
                {{ statusLabel[task.status] || task.status }}
              </span>
            </template>
          </div>
          <div class="task-actions">
            <button
              class="icon-btn"
              @click.stop="startEditTask(task)"
              title="重命名"
            >
              <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
              </svg>
            </button>
            <button class="icon-btn danger" @click.stop="$emit('delete', task.id)" title="删除">
              <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
              </svg>
            </button>
          </div>
        </div>
      </TransitionGroup>
    </div>
  </div>
</template>

<style scoped>
.task-history {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-1) 0;
}

.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-8) var(--space-4);
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}

.group {
  border-bottom: 1px solid var(--color-border-light);
}

.group-header {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  cursor: pointer;
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-secondary);
  user-select: none;
  transition: background var(--transition-fast);
}
.group-header:hover {
  background: var(--color-bg-muted);
}

.group-arrow {
  flex-shrink: 0;
  transition: transform var(--transition-fast);
  color: var(--color-text-muted);
}
.group-arrow.collapsed {
  transform: rotate(-90deg);
}

.group-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.group-count {
  font-size: var(--text-xs);
  color: var(--color-text-muted);
  background: var(--color-bg-muted);
  padding: 1px 7px;
  border-radius: var(--radius-full);
  flex-shrink: 0;
}

.group-tasks {
  /* tasks container */
}

.task-item {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3) var(--space-2) var(--space-6);
  cursor: pointer;
  transition: background var(--transition-fast);
  border-bottom: 1px solid var(--color-border-light);
}
.task-item:hover {
  background: var(--color-bg-muted);
}
.task-item:hover .task-actions {
  opacity: 1;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: var(--radius-full);
  flex-shrink: 0;
}
.s-success { background: var(--color-success); }
.s-failed { background: var(--color-danger); }
.s-running {
  background: var(--color-warning);
  animation: pulse 2s ease-in-out infinite;
}

.task-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.task-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: var(--text-sm);
  color: var(--color-text);
  font-weight: 500;
}

.task-meta {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  font-size: var(--text-xs);
  color: var(--color-text-muted);
}

.platform-badge {
  display: inline-flex;
  align-items: center;
  padding: 0px 5px;
  background: var(--color-bg-subtle);
  border-radius: 3px;
  font-size: 10px;
  color: var(--color-text-muted);
}
.platform-logo {
  width: 14px;
  height: 14px;
}

.task-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity var(--transition-fast);
  flex-shrink: 0;
}

.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: var(--radius-sm);
  color: var(--color-text-muted);
  transition: all var(--transition-fast);
}
.icon-btn:hover {
  background: var(--color-bg-subtle);
  color: var(--color-text-secondary);
}
.icon-btn.danger:hover {
  background: var(--color-danger-light);
  color: var(--color-danger);
}

.edit-input {
  flex: 1;
  padding: 3px 8px;
  font-size: var(--text-sm);
  border: 1px solid var(--color-primary);
  border-radius: var(--radius-sm);
  outline: none;
  background: var(--color-bg-elevated);
  box-shadow: var(--shadow-focus);
}

/* 列表动画 */
.list-enter-active {
  transition: all 200ms ease-out;
}
.list-leave-active {
  transition: all 150ms ease-in;
}
.list-enter-from {
  opacity: 0;
  transform: translateX(-8px);
}
.list-leave-to {
  opacity: 0;
  transform: translateX(8px);
}
</style>
