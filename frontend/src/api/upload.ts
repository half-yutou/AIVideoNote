import axios from 'axios'

// 前后端统一的允许上传的媒体文件扩展名
export const ALLOWED_MEDIA_EXTENSIONS = [
  '.mp4', '.mkv', '.avi', '.mov', '.webm', '.flv',
  '.mp3', '.wav', '.flac', '.aac', '.ogg', '.m4a',
]

// 最大文件大小 2GB
export const MAX_UPLOAD_SIZE = 2 * 1024 * 1024 * 1024

export interface UploadResult {
  path: string
  filename: string
  size: number
}

/**
 * 校验文件格式和大小
 * @returns 错误信息，null 表示通过
 */
export function validateFile(file: File): string | null {
  const ext = '.' + file.name.split('.').pop()?.toLowerCase()
  if (!ALLOWED_MEDIA_EXTENSIONS.includes(ext)) {
    return `不支持的文件格式 ${ext}，允许: ${ALLOWED_MEDIA_EXTENSIONS.join(', ')}`
  }
  if (file.size > MAX_UPLOAD_SIZE) {
    return '文件过大，最大支持 2GB'
  }
  return null
}

/**
 * 上传文件到服务端
 * @param file 文件对象
 * @param onProgress 进度回调 (0-100)
 */
export function uploadFile(file: File, onProgress?: (percent: number) => void) {
  const form = new FormData()
  form.append('file', file)
  return axios.post<ApiResponse<UploadResult>>('/api/v1/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress(e) {
      if (onProgress && e.total) {
        onProgress(Math.round((e.loaded / e.total) * 100))
      }
    },
  })
}
