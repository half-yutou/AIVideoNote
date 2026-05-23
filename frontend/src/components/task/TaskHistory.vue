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
    <div v-if="tasks.length === 0" class="empty">暂无任务记录</div>

    <div v-for="group in groups" :key="group.name" class="group">
      <div class="group-header" @click="toggleGroup(group.name)">
        <span class="group-arrow">{{ collapsedGroups[group.name] ? '▸' : '▾' }}</span>
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
          class="edit-btn"
          @click.stop="startEditGroup(group.name)"
          title="重命名分组"
        >✎</button>
      </div>

      <div v-show="!collapsedGroups[group.name]" class="group-tasks">
        <div
          v-for="task in group.tasks"
          :key="task.id"
          class="task-item"
          @click="$emit('select', task.id)"
        >
          <span :class="['status-dot', statusClass(task.status)]" />
          <span class="platform-badge">{{ task.platform }}</span>
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
                {{ statusLabel[task.status] || task.status }}
                <span v-if="task.video_id" class="vid"> · {{ task.video_id.slice(0, 8) }}</span>
              </span>
            </template>
          </div>
          <span class="task-time">{{ task.created_at.slice(-8) }}</span>
          <button
            class="edit-btn"
            @click.stop="startEditTask(task)"
            title="重命名任务"
          >✎</button>
          <button class="del-btn" @click.stop="$emit('delete', task.id)">×</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.task-history {
  max-height: calc(100vh - 60px);
  overflow-y: auto;
}
.empty {
  color: #999;
  font-size: 13px;
  text-align: center;
  padding: 12px;
}
.group {
  border-bottom: 1px solid #eee;
}
.group-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 10px;
  cursor: pointer;
  background: #f8f9fa;
  font-size: 13px;
  font-weight: 500;
  color: #444;
  user-select: none;
}
.group-header:hover {
  background: #f0f1f2;
}
.group-arrow {
  font-size: 10px;
  width: 12px;
  color: #999;
  flex-shrink: 0;
}
.group-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.group-count {
  color: #bbb;
  font-size: 11px;
  flex-shrink: 0;
}
.group-tasks {
  /* tasks indented under group */
}
.task-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 10px 7px 28px;
  cursor: pointer;
  font-size: 12px;
  border-bottom: 1px solid #f5f5f5;
}
.task-item:hover {
  background: #f5f7fa;
}
.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}
.s-success { background: #4caf50; }
.s-failed { background: #e53935; }
.s-running { background: #ff9800; }
.platform-badge {
  padding: 1px 5px;
  background: #eee;
  border-radius: 3px;
  font-size: 10px;
  color: #888;
  flex-shrink: 0;
}
.task-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 1px;
}
.task-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #333;
  font-weight: 500;
}
.task-meta {
  font-size: 10px;
  color: #999;
}
.vid {
  color: #ccc;
}
.task-time {
  color: #ccc;
  font-size: 10px;
  flex-shrink: 0;
}
.edit-btn {
  border: none;
  background: none;
  color: #ccc;
  cursor: pointer;
  font-size: 12px;
  padding: 0 2px;
  flex-shrink: 0;
}
.edit-btn:hover {
  color: #666;
}
.del-btn {
  border: none;
  background: none;
  color: #ccc;
  cursor: pointer;
  font-size: 14px;
  padding: 0 2px;
  flex-shrink: 0;
}
.del-btn:hover {
  color: #e53935;
}
.edit-input {
  flex: 1;
  padding: 3px 6px;
  font-size: 12px;
  border: 1px solid #4a90d9;
  border-radius: 3px;
  outline: none;
  background: #fff;
}
</style>
