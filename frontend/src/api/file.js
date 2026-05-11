import request from './request'

export const getFiles = (params) => request.get('/files', { params })
export const uploadFile = (data) => request.post('/files/upload', data)
export const uploadAvatar = (data) => request.post('/files/avatar', data)
export const getAvatarUrl = (id) => request.get(`/files/avatar/${id}/url`)
export const getFileUrl = (id) => request.get(`/files/${id}/url`)
export const getDownloadUrl = (id) => request.get(`/files/${id}/download-url`)
export const deleteFile = (id) => request.delete(`/files/${id}`)
