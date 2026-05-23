import request from './request'

export function getCookies() {
  return request.get<ApiResponse<CookieData[]>>('/cookies')
}

export function saveCookie(platform: string, content: string) {
  return request.post<ApiResponse<CookieData>>('/cookies', { platform, content })
}

export function deleteCookie(platform: string) {
  return request.delete<ApiResponse<null>>(`/cookies/${platform}`)
}
