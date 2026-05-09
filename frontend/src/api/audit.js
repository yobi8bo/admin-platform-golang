import request from './request'

export function loginLogs(params) {
  return request.get('/audit/login-logs', { params })
}

export function operationLogs(params) {
  return request.get('/audit/operation-logs', { params })
}
