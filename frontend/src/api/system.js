import request from './request'

export const userApi = {
  list: (params) => request.get('/system/users', { params }),
  create: (data) => request.post('/system/users', data),
  update: (id, data) => request.put(`/system/users/${id}`, data),
  remove: (id) => request.delete(`/system/users/${id}`)
}

export const roleApi = {
  list: () => request.get('/system/roles'),
  create: (data) => request.post('/system/roles', data),
  update: (id, data) => request.put(`/system/roles/${id}`, data),
  remove: (id) => request.delete(`/system/roles/${id}`)
}

export const menuApi = {
  tree: () => request.get('/system/menus/tree'),
  create: (data) => request.post('/system/menus', data),
  update: (id, data) => request.put(`/system/menus/${id}`, data),
  remove: (id) => request.delete(`/system/menus/${id}`)
}

export const deptApi = {
  tree: () => request.get('/system/depts/tree'),
  create: (data) => request.post('/system/depts', data),
  update: (id, data) => request.put(`/system/depts/${id}`, data),
  remove: (id) => request.delete(`/system/depts/${id}`)
}
