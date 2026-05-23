import request from './request'

export function getGroupList() {
  return request.get<ApiResponse<GroupData[]>>('/group/list')
}

export function createGroup(name: string) {
  return request.post<ApiResponse<GroupData>>('/group/create', { name })
}

export function renameGroup(id: string, name: string) {
  return request.put<ApiResponse<null>>(`/group/${id}/name`, { name })
}

export function deleteGroup(id: string) {
  return request.delete<ApiResponse<null>>(`/group/${id}`)
}
