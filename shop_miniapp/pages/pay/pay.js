const api = require('../../utils/request.js')

Page({
  data: {
    orderId: null,
    orderNo: null,
    order: null,
    payTypes: [
      { id: 1, name: '微信支付', icon: '/static/pay/wechat.png' },
      { id: 2, name: '支付宝', icon: '/static/pay/alipay.png' }
    ],
    selectedPayType: 1,
    paying: false,
    paid: false // 是否已支付
  },

  onLoad(options) {
    if (options.order_id) {
      this.setData({ orderId: parseInt(options.order_id) })
    }
    if (options.order_no) {
      this.setData({ orderNo: options.order_no })
    }
    if (this.data.orderId) {
      this.loadOrder()
    }
  },

  async loadOrder() {
    try {
      const res = await api.getOrder(this.data.orderId)
      this.setData({
        order: res.data,
        orderNo: res.data.order_no,
        // 如果订单已支付，更新状态
        paid: res.data.pay_status === 1
      })
    } catch (err) {
      console.error(err)
    }
  },

  onPayTypeChange(e) {
    const { id } = e.currentTarget.dataset
    console.log('切换支付方式:', id)
    this.setData({ 
      selectedPayType: parseInt(id)
    }, () => {
      console.log('当前选中支付方式:', this.data.selectedPayType)
    })
  },

  async onPay() {
    if (this.data.paying || this.data.paid || !this.data.order) return

    this.setData({ paying: true })
    try {
      const res = await api.getPayURL({
        order_id: this.data.order.id,
        pay_type: this.data.selectedPayType
      })

      // 检查是否是模拟支付
      if (res.data.mock_pay) {
        // 模拟支付成功
        wx.showToast({ 
          title: '支付成功', 
          icon: 'success',
          duration: 2000
        })
        
        // 更新页面状态
        this.setData({ 
          paid: true,
          paying: false 
        })
        
        // 延迟跳转至订单详情页
        setTimeout(() => {
          wx.redirectTo({ 
            url: `/pages/order-detail/order-detail?id=${this.data.order.id}` 
          })
        }, 2000)
      } else if (res.data.pay_url) {
        // 真实支付，显示二维码或跳转
        wx.showModal({
          title: '提示',
          content: '请扫描支付二维码完成支付',
          showCancel: false,
          success: () => {
            this.setData({ paying: false })
          }
        })
      } else {
        wx.showToast({ title: '支付开发中', icon: 'none' })
        this.setData({ paying: false })
      }
    } catch (err) {
      console.error('支付失败:', err)
      wx.showToast({ 
        title: err.msg || '支付失败', 
        icon: 'none' 
      })
      this.setData({ paying: false })
    }
  },

  goOrderList() {
    wx.redirectTo({ url: '/pages/order-list/order-list' })
  },

  goHome() {
    wx.switchTab({ url: '/pages/index/index' })
  },

  goOrderDetail() {
    wx.redirectTo({ 
      url: `/pages/order-detail/order-detail?id=${this.data.order.id}` 
    })
  }
})
