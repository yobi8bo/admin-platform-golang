import request from './request'

export const getLoginLogs = (params) => request.get('/audit/login-logs', { params })
export const getOperationLogs = (params) => request.get('/audit/operation-logs', { params })
