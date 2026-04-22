const api = require('../../utils/request.js')

Page({
  data: {
    orderNo: null,
    order: null,
    payTypes: [
      { id: 1, name: '微信支付', icon: '/static/pay/wechat.png' },
      { id: 2, name: '支付宝', icon: '/static/pay/alipay.png' }
    ],
    selectedPayType: 1
  },

  onLoad(options) {
    if (options.order_no) {
      this.setData({ orderNo: options.order_no })
      this.loadOrder()
    }
  },

  async loadOrder() {
    try {
      const res = await api.getOrder(this.data.orderNo)
      this.setData({ order: res.data })
    } catch (err) {
      console.error(err)
    }
  },

  onPayTypeChange(e) {
    const { id } = e.currentTarget.dataset
    this.setData({ selectedPayType: id })
  },

  async onPay() {
    try {
      const res = await api.getPayURL({
        order_id: this.data.order.id,
        pay_type: this.data.selectedPayType
      })

      if (res.data.pay_url) {
        wx.redirectTo({ url: `/pages/pay/pay?order_no=${this.data.orderNo}` })
      } else {
        wx.showToast({ title: '支付开发中', icon: 'none' })
      }
    } catch (err) {
      console.error(err)
    }
  }
})
