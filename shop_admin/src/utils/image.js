// 图片URL处理工具

// 获取完整的图片URL
export function getImageUrl(url) {
  if (!url) return ''
  
  // 如果已经是完整URL，直接返回
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url
  }
  
  // 拼接服务器地址
  const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8686'
  return `${baseURL}${url}`
}

// 获取多个图片的完整URL
export function getImageUrls(urls) {
  if (!urls || !Array.isArray(urls)) return []
  return urls.map(url => getImageUrl(url))
}
