import axios from 'axios'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 300000,
  headers: { 'Content-Type': 'application/json' },
})

request.interceptors.response.use(
  (res) => {
    const data = res.data as ApiResponse<unknown>
    if (data.code !== 0) {
      return Promise.reject(new Error(data.message || '请求失败'))
    }
    return res
  },
  (err) => {
    const msg = err.response?.data?.message || err.message || '网络错误'
    return Promise.reject(new Error(msg))
  }
)

export default request
