declare global {
  interface ApiResponse<T> {
    code: number
    message: string
    data?: T
  }

  interface GenerateRequest {
    video_url: string
    platform: string
    quality?: string
    model_name: string
    provider_id: string
    group_id?: string
    format?: string[]
    style?: string
    extras?: string
    link?: boolean
  }

  interface GroupData {
    id: string
    name: string
    created_at: string
  }

  interface TaskStatusData {
    task_id: string
    status: string
    video_url: string
    platform: string
    error_message?: string
    created_at: string
    markdown?: string
    note_id?: string
  }

  interface TaskItemData {
    id: string
    status: string
    video_url: string
    platform: string
    video_id: string
    name: string
    group_name: string
    error_message: string
    created_at: string
  }

  interface LLMProviderData {
    id: string
    name: string
    base_url: string
    type: string
    logo: string
    enabled: boolean
    created_at: string
    updated_at: string
  }

  interface ProviderRequest {
    name: string
    api_key: string
    base_url: string
    type: string
    logo?: string
  }

  interface ProviderUpdateRequest {
    name?: string
    api_key?: string
    base_url?: string
    type?: string
    logo?: string
    enabled?: boolean
  }

  interface CookieData {
    platform: string
    content: string
    created_at: string
    updated_at: string
  }
}

export {}