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
    paying: false
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
        orderNo: res.data.order_no
      })
    } catch (err) {
      console.error(err)
    }
  },

  onPayTypeChange(e) {
    const { id } = e.currentTarget.dataset
    this.setData({ selectedPayType: id })
  },

  async onPay() {
    if (this.data.paying || !this.data.order) return

    this.setData({ paying: true })
    try {
      const res = await api.getPayURL({
        order_id: this.data.order.id,
        pay_type: this.data.selectedPayType
      })

      if (res.data.pay_url) {
        wx.showModal({
          title: '开发环境支付',
          content: '已生成支付链接。是否直接模拟支付成功，并将订单状态切换为待发货？',
          confirmText: '模拟成功',
          cancelText: '稍后支付',
          success: async (modalRes) => {
            if (modalRes.confirm) {
              try {
                await api.mockPaySuccess(this.data.order.id)
                wx.showToast({ title: '支付成功', icon: 'success' })
                setTimeout(() => {
                  wx.redirectTo({ url: `/pages/order-detail/order-detail?id=${this.data.order.id}` })
                }, 800)
              } catch (mockErr) {
                console.error(mockErr)
              }
            }
          },
          complete: () => {
            this.setData({ paying: false })
          }
        })
      } else {
        wx.showToast({ title: '支付开发中', icon: 'none' })
        this.setData({ paying: false })
      }
    } catch (err) {
      console.error(err)
      this.setData({ paying: false })
    }
  },

  goOrderList() {
    wx.redirectTo({ url: '/pages/order-list/order-list' })
  },

  goHome() {
    wx.switchTab({ url: '/pages/index/index' })
  }
})
