import axios from 'axios'

export interface ApiResponse<T> {
  success: boolean
  message: string
  data: T
}

export const http = axios.create({
  baseURL: '/api',
  timeout: 120000,
})

http.interceptors.response.use((resp) => resp.data)

export function analyzeRelease(payload: FormData) {
  return http.post<ApiResponse<any>>('/releases/analyze', payload, {
    timeout: 300000,
  })
}

export function getReleaseReport(requestNo: string) {
  return http.get<ApiResponse<any>>(`/releases/${requestNo}`)
}

export function approvalAction(requestNo: string, data: { approver: string; action: string; comment?: string }) {
  return http.post<ApiResponse<any>>(`/releases/${requestNo}/approval`, data)
}

export function getBaseline(env: string, appName: string) {
  return http.get<ApiResponse<Record<string, string>>>(`/baselines/${env}/${appName}`)
}

export function updateBaseline(env: string, appName: string, values: Record<string, string>) {
  return http.post<ApiResponse<string>>(`/baselines/${env}/${appName}`, { values })
}
