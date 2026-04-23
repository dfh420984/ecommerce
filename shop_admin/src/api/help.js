import request from '@/utils/request'

// ============ 分类管理 ============

// 获取分类列表
export function getHelpCategories() {
  return request({
    url: '/admin/help/categories',
    method: 'get'
  })
}

// 创建分类
export function createHelpCategory(data) {
  return request({
    url: '/admin/help/categories',
    method: 'post',
    data
  })
}

// 更新分类
export function updateHelpCategory(id, data) {
  return request({
    url: `/admin/help/categories/${id}`,
    method: 'put',
    data
  })
}

// 删除分类
export function deleteHelpCategory(id) {
  return request({
    url: `/admin/help/categories/${id}`,
    method: 'delete'
  })
}

// ============ 问题管理 ============

// 获取问题列表
export function getHelpQuestions(params) {
  return request({
    url: '/admin/help/questions',
    method: 'get',
    params
  })
}

// 获取问题详情
export function getHelpQuestion(id) {
  return request({
    url: `/admin/help/questions/${id}`,
    method: 'get'
  })
}

// 创建问题
export function createHelpQuestion(data) {
  return request({
    url: '/admin/help/questions',
    method: 'post',
    data
  })
}

// 更新问题
export function updateHelpQuestion(id, data) {
  return request({
    url: `/admin/help/questions/${id}`,
    method: 'put',
    data
  })
}

// 删除问题
export function deleteHelpQuestion(id) {
  return request({
    url: `/admin/help/questions/${id}`,
    method: 'delete'
  })
}

// 批量更新问题状态
export function batchUpdateQuestionsStatus(ids, status) {
  return request({
    url: '/admin/help/questions/batch-status',
    method: 'put',
    data: { ids, status }
  })
}

// 获取统计信息
export function getHelpStatistics() {
  return request({
    url: '/admin/help/statistics',
    method: 'get'
  })
}
