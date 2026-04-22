// 图片 URL 工具函数
const app = getApp()

/**
 * 格式化图片 URL
 * 如果已经是完整 URL 则直接返回，否则拼接 imageUrl 前缀
 */
function formatImageUrl(url) {
  if (!url) {
    return ''
  }
  
  // 如果已经是完整的 URL（http:// 或 https:// 或 data:image），直接返回
  if (url.startsWith('http://') || url.startsWith('https://') || url.startsWith('data:image')) {
    return url
  }
  
  // 拼接图片基础 URL
  return app.globalData.imageUrl + url
}

/**
 * 批量格式化图片数组
 */
function formatImageUrls(urls) {
  if (!urls || !Array.isArray(urls)) {
    return []
  }
  return urls.map(url => formatImageUrl(url))
}

module.exports = {
  formatImageUrl,
  formatImageUrls
}
