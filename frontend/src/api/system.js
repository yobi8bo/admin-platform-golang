import request from './request'

export const getUsers = (params) => request.get('/system/users', { params })
export const createUser = (data) => request.post('/system/users', data)
export const updateUser = (id, data) => request.put(`/system/users/${id}`, data)
export const deleteUser = (id) => request.delete(`/system/users/${id}`)

export const getRoles = () => request.get('/system/roles')
export const createRole = (data) => request.post('/system/roles', data)
export const updateRole = (id, data) => request.put(`/system/roles/${id}`, data)
export const deleteRole = (id) => request.delete(`/system/roles/${id}`)

export const getMenuTree = () => request.get('/system/menus/tree')
export const getMyMenuTree = () => request.get('/system/menus/my-tree')
export const createMenu = (data) => request.post('/system/menus', data)
export const updateMenu = (id, data) => request.put(`/system/menus/${id}`, data)
export const deleteMenu = (id) => request.delete(`/system/menus/${id}`)

export const getDeptTree = () => request.get('/system/depts/tree')
export const createDept = (data) => request.post('/system/depts', data)
export const updateDept = (id, data) => request.put(`/system/depts/${id}`, data)
export const deleteDept = (id) => request.delete(`/system/depts/${id}`)
