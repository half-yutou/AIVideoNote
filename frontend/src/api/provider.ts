import request from './request'

export function createProvider(data: ProviderRequest) {
  return request.post<ApiResponse<LLMProviderData>>('/provider', data)
}

export function getProviders() {
  return request.get<ApiResponse<LLMProviderData[]>>('/provider')
}

export function getProviderById(id: string) {
  return request.get<ApiResponse<LLMProviderData>>(`/provider/${id}`)
}

export function updateProvider(id: string, data: ProviderUpdateRequest) {
  return request.put<ApiResponse<LLMProviderData>>(`/provider/${id}`, data)
}

export function deleteProvider(id: string) {
  return request.delete<ApiResponse<null>>(`/provider/${id}`)
}

export function listModels(providerId: string) {
  return request.get<ApiResponse<string[]>>('/model/list', {
    params: { provider_id: providerId },
  })
}

export function testConnection(providerId: string) {
  return request.get<ApiResponse<null>>(`/provider/${providerId}/test`)
}
