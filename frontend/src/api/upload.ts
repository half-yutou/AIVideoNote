import axios from 'axios'

export function uploadFile(file: File) {
  const form = new FormData()
  form.append('file', file)
  return axios.post<ApiResponse<{ url: string; filename: string }>>('/api/v1/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}
