import { ref, type Ref } from 'vue'

export interface ToastOptions {
  message: string
  type?: 'success' | 'error' | 'warning' | 'info'
  duration?: number
}

interface ToastItem {
  id: number
  message: string
  type: 'success' | 'error' | 'warning' | 'info'
}

const toasts: Ref<ToastItem[]> = ref([])
let nextId = 0

export function useToast() {
  function show(options: ToastOptions) {
    const id = nextId++
    const type = options.type || 'info'
    const duration = options.duration ?? 3500
    toasts.value.push({ id, message: options.message, type })
    if (duration > 0) {
      setTimeout(() => remove(id), duration)
    }
  }

  function success(message: string, duration?: number) {
    show({ message, type: 'success', duration })
  }

  function error(message: string, duration?: number) {
    show({ message, type: 'error', duration: duration ?? 5000 })
  }

  function warning(message: string, duration?: number) {
    show({ message, type: 'warning', duration })
  }

  function info(message: string, duration?: number) {
    show({ message, type: 'info', duration })
  }

  function remove(id: number) {
    const idx = toasts.value.findIndex((t) => t.id === id)
    if (idx !== -1) toasts.value.splice(idx, 1)
  }

  return { toasts, show, success, error, warning, info, remove }
}
