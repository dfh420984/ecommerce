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

module.exports = {
  formatTime,
  formatPrice,
  getStatusText,
  getStatusClass
}
