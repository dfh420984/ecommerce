import request from '@/utils/request'

// 获取配置列表
export const getConfigs = () => {
  return request.get('/admin/configs')
}

// 获取单个配置
export const getConfig = (id) => {
  return request.get(`/admin/configs/${id}`)
}

// 创建配置
export const createConfig = (data) => {
  return request.post('/admin/configs', data)
}

// 更新配置
export const updateConfig = (id, data) => {
  return request.put(`/admin/configs/${id}`, data)
}

// 删除配置
export const deleteConfig = (id) => {
  return request.delete(`/admin/configs/${id}`)
}
