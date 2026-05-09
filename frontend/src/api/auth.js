import request from './request'

export function login(data) {
  return request.post('/auth/login', data)
}

export function logout() {
  return request.post('/auth/logout')
}

export function profile() {
  return request.get('/auth/profile')
}

export function updateProfile(data) {
  return request.put('/auth/profile', data)
}

export function updatePassword(data) {
  return request.put('/auth/password', data)
}

export function permissions() {
  return request.get('/auth/permissions')
}

export function myMenus() {
  return request.get('/system/menus/my-tree')
}
