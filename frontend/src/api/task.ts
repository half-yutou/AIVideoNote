import request from './request'

export function generateNote(data: GenerateRequest) {
  return request.post<ApiResponse<{ task_id: string }>>('/note/generate', data)
}

export function getTaskStatus(id: string) {
  return request.get<ApiResponse<TaskStatusData>>(`/task/${id}/status`)
}

export function getTaskList() {
  return request.get<ApiResponse<TaskItemData[]>>('/task/list')
}

export function deleteTask(id: string) {
  return request.delete<ApiResponse<null>>(`/task/${id}`)
}

export function renameTask(id: string, name: string) {
  return request.put<ApiResponse<null>>(`/task/${id}/name`, { name })
}

export function getHealth() {
  return request.get('/health')
}
