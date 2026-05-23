import { defineStore } from 'pinia'
import { ref } from 'vue'
import { generateNote, getTaskStatus, getTaskList, deleteTask, renameTask } from '@/api/task'

let pollToken: symbol | null = null

if (import.meta.hot) {
  import.meta.hot.dispose(() => {
    pollToken = null
  })
}

export const useTaskStore = defineStore('task', () => {
  const currentTaskId = ref<string | null>(null)
  const currentStatus = ref<string>('')
  const currentMarkdown = ref<string>('')
  const currentError = ref<string>('')
  const tasks = ref<TaskItemData[]>([])
  const loading = ref(false)
  const error = ref('')

  async function submitTask(req: GenerateRequest) {
    loading.value = true
    error.value = ''
    try {
      const res = await generateNote(req)
      currentTaskId.value = res.data.data!.task_id
      currentStatus.value = 'PENDING'
      currentMarkdown.value = ''
      startPolling()
      return res.data.data!.task_id
    } catch (e: any) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  function startPolling() {
    pollToken = Symbol()
    poll(pollToken)
  }

  async function poll(token: symbol) {
    if (pollToken !== token) return
    if (!currentTaskId.value) return
    const taskId = currentTaskId.value
    try {
      const res = await getTaskStatus(taskId)
      if (pollToken !== token) return
      if (currentTaskId.value !== taskId) return
      const data = res.data.data!
      currentStatus.value = data.status
      if (data.markdown) {
        currentMarkdown.value = data.markdown
      }
      if (data.error_message) {
        currentError.value = data.error_message
      }
      if (data.status === 'SUCCESS' || data.status?.endsWith('_FAILED') || data.status === 'FAILED') {
        pollToken = null
        fetchTasks()
        return
      }
    } catch {
      // 轮询失败不报错
    }
    if (pollToken === token) {
      setTimeout(() => poll(token), 3000)
    }
  }

  function stopPolling() {
    pollToken = null
  }

  async function fetchTasks() {
    try {
      const res = await getTaskList()
      tasks.value = res.data.data || []
    } catch {
      // ignore
    }
  }

  async function removeTask(id: string) {
    try {
      await deleteTask(id)
      tasks.value = tasks.value.filter((t) => t.id !== id)
    } catch {
      // ignore
    }
  }

  async function renameTaskAction(id: string, name: string) {
    await renameTask(id, name)
    const t = tasks.value.find((x) => x.id === id)
    if (t) t.name = name
  }

  async function loadTask(id: string) {
    stopPolling()
    currentTaskId.value = id
    currentError.value = ''
    try {
      const res = await getTaskStatus(id)
      const data = res.data.data!
      currentStatus.value = data.status
      currentMarkdown.value = data.markdown || ''
      currentError.value = data.error_message || ''
      if (data.status !== 'SUCCESS' && !data.status?.endsWith('_FAILED') && data.status !== 'FAILED') {
        startPolling()
      }
    } catch {
      // ignore
    }
  }

  return {
    currentTaskId,
    currentStatus,
    currentMarkdown,
    currentError,
    tasks,
    loading,
    error,
    submitTask,
    fetchTasks,
    removeTask,
    renameTask: renameTaskAction,
    loadTask,
    stopPolling,
  }
})
