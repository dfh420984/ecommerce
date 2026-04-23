const formatNumber = n => {
  n = n.toString()
  return n[1] ? n : '0' + n
}

const formatTime = date => {
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hour = date.getHours()
  const minute = date.getMinutes()
  const second = date.getSeconds()
  return `${[year, month, day].map(formatNumber).join('-')} ${[hour, minute, second].map(formatNumber).join(':')}`
}

const formatPrice = price => {
  return '¥' + (Math.round(price * 100) / 100).toFixed(2)
}

const getStatusText = status => {
  const texts = {
    1: '待付款',
    2: '待发货',
    3: '配送中',
    4: '已收货',
    5: '已完成',
    6: '已取消',
    7: '退款中',
    8: '已退款'
  }
  return texts[status] || '未知'
}

const getStatusClass = status => {
  const classes = {
    1: 'status-pending',
    2: 'status-paid',
    3: 'status-shipped',
    4: 'status-received',
    5: 'status-completed',
    6: 'status-cancelled'
  }
  return classes[status] || ''
}

/**
 * 解析富文本内容，提取视频和HTML
 * @param {string} html - 富文本HTML内容
 * @returns {Array} - 包含视频和HTML片段的数组
 */
const parseRichContent = html => {
  if (!html) return []
  
  const segments = []
  // 匹配 WangEditor 生成的视频标签结构
  // 支持两种格式：
  // 1. <video src="...">...</video>
  // 2. <video><source src="..." /></video>
  const videoRegex = /<video[^>]*>([\s\S]*?)<\/video>/gi
  
  let lastIndex = 0
  let match
  
  while ((match = videoRegex.exec(html)) !== null) {
    // 添加视频前的HTML片段
    if (match.index > lastIndex) {
      const textSegment = html.substring(lastIndex, match.index)
      if (textSegment.trim()) {
        segments.push({
          type: 'html',
          content: textSegment
        })
      }
    }
    
    // 提取视频信息
    const videoTag = match[0]
    const videoContent = match[1]
    
    // 尝试从 <source> 标签获取 src
    let src = ''
    const sourceMatch = videoContent.match(/<source[^>]*src=["']([^"']+)["']/i)
    if (sourceMatch) {
      src = sourceMatch[1]
    } else {
      // 如果没有 <source>，尝试从 video 标签获取 src
      const videoSrcMatch = videoTag.match(/<video[^>]*src=["']([^"']+)["']/i)
      if (videoSrcMatch) {
        src = videoSrcMatch[1]
      }
    }
    
    // 提取 poster
    const posterMatch = videoTag.match(/poster=["']([^"']+)["']/i)
    const poster = posterMatch ? posterMatch[1] : ''
    
    // 只有当找到有效的 src 时才添加视频片段
    if (src) {
      segments.push({
        type: 'video',
        src: src,
        poster: poster
      })
    } else {
      // 如果没有找到视频地址，作为HTML处理
      segments.push({
        type: 'html',
        content: videoTag
      })
    }
    
    lastIndex = match.index + match[0].length
  }
  
  // 添加剩余的HTML片段
  if (lastIndex < html.length) {
    const textSegment = html.substring(lastIndex)
    if (textSegment.trim()) {
      segments.push({
        type: 'html',
        content: textSegment
      })
    }
  }
  
  // 如果没有视频，返回原始HTML
  if (segments.length === 0) {
    segments.push({
      type: 'html',
      content: html
    })
  }
  
  return segments
}

module.exports = {
  formatTime,
  formatPrice,
  getStatusText,
  getStatusClass,
  parseRichContent
}
