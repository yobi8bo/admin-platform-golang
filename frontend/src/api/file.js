import request from './request'

export function fileList(params) {
  return request.get('/files', { params })
}

export function uploadFile(file) {
  const form = new FormData()
  form.append('file', file)
  return request.post('/files/upload', form)
}

export function uploadAvatar(file) {
  const form = new FormData()
  form.append('file', file)
  return request.post('/files/avatar', form)
}

export function fileUrl(id) {
  return request.get(`/files/${id}/url`)
}

export function fileDownloadUrl(id) {
  return request.get(`/files/${id}/download-url`)
}

export function deleteFile(id) {
  return request.delete(`/files/${id}`)
}

export function avatarUrl(id) {
  return request.get(`/files/avatar/${id}/url`)
}
