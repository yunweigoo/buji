const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export interface ApiEnvelope<T> {
  data: T
  requestId: string
}

type ApiRequestOptions = Omit<UniApp.RequestOptions, 'url'>

export async function apiRequest<T>(path: string, options: ApiRequestOptions = {}): Promise<T> {
  return new Promise((resolve, reject) => {
    uni.request({
      ...options,
      url: `${API_BASE_URL}${path}`,
      header: {
        'content-type': 'application/json',
        ...options.header,
      },
      success(response) {
        if (response.statusCode >= 200 && response.statusCode < 300) {
          resolve((response.data as ApiEnvelope<T>).data)
          return
        }
        reject(new Error(`请求失败：${response.statusCode}`))
      },
      fail(error) {
        reject(new Error(error.errMsg || '网络请求失败'))
      },
    })
  })
}

