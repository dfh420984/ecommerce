// pages/logistics/logistics.js
const api = require('../../services/api')

Page({
  data: {
    orderId: null,
    orderNo: '',
    expressCompany: '',
    expressNo: '',
    tracks: [],
    loading: true
  },

  onLoad(options) {
    if (options.order_id) {
      this.setData({ orderId: options.order_id })
      this.loadLogistics(options.order_id)
    }
  },

  // 加载物流信息
  loadLogistics(orderId) {
    wx.showLoading({ title: '加载中' })
    
    api.getOrderLogistics(orderId).then(res => {
      const data = res.data
      
      this.setData({
        orderNo: data.order_no || '',
        expressCompany: data.express_company || '',
        expressNo: data.express_no || '',
        tracks: data.tracks || [],
        loading: false
      })
      
      wx.hideLoading()
    }).catch(err => {
      wx.hideLoading()
      this.setData({ loading: false })
      wx.showToast({ title: err.message || '加载失败', icon: 'none' })
    })
  },

  // 复制快递单号
  copyExpressNo() {
    wx.setClipboardData({
      data: this.data.expressNo,
      success: () => {
        wx.showToast({ title: '已复制', icon: 'success' })
      }
    })
  }
})
