// pages/refund-detail/refund-detail.js
const api = require('../../services/api')
const { getImageUrl } = require('../../utils/image')

Page({
  data: {
    refundId: null,
    refund: null,
    loading: true
  },

  onLoad(options) {
    if (options.id) {
      this.setData({ refundId: options.id })
      this.loadRefundDetail(options.id)
    }
  },

  // 加载退款详情
  loadRefundDetail(id) {
    wx.showLoading({ title: '加载中' })
    
    api.getRefundDetail(id).then(res => {
      const refund = res.data
      
      // 处理图片URL
      if (refund.images && refund.images.length > 0) {
        refund.imageUrls = refund.images.map(img => getImageUrl(img))
      }
      
      this.setData({
        refund,
        loading: false
      })
      
      wx.hideLoading()
    }).catch(err => {
      wx.hideLoading()
      this.setData({ loading: false })
      wx.showToast({ title: err.message || '加载失败', icon: 'none' })
    })
  },

  // 预览图片
  previewImage(e) {
    const urls = this.data.refund.imageUrls
    const current = e.currentTarget.dataset.url
    
    wx.previewImage({
      current,
      urls
    })
  },

  // 获取状态文本
  getStatusText(status) {
    const map = {
      1: '待审核',
      2: '已通过',
      3: '已拒绝',
      4: '退款中',
      5: '已退款'
    }
    return map[status] || '未知'
  },

  // 获取状态样式类
  getStatusClass(status) {
    const map = {
      1: 'status-pending',
      2: 'status-approved',
      3: 'status-rejected',
      4: 'status-refunding',
      5: 'status-refunded'
    }
    return map[status] || ''
  }
})
