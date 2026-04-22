const api = require('../../utils/request.js')
const image = require('../../utils/image')
const app = getApp()

Page({
  data: {
    orderNo: null,
    order: null,
    addresses: [],
    selectedAddress: null,
    cart_ids: [],
    remark: ''
  },

  onLoad(options) {
    if (options.order_no) {
      this.setData({ orderNo: options.order_no })
      this.loadOrder()
    }
    if (options.cart_ids) {
      this.setData({ cart_ids: options.cart_ids.split(',').map(id => parseInt(id)) })
    }
    this.loadAddresses()
  },

  async loadOrder() {
    try {
      const res = await api.getOrder(this.data.orderNo)
      const order = res.data
      // 格式化订单商品图片
      if (order.items) {
        order.items = order.items.map(item => ({
          ...item,
          product_image: image.formatImageUrl(item.product_image)
        }))
      }
      this.setData({ order })
    } catch (err) {
      console.error(err)
    }
  },

  async loadAddresses() {
    try {
      const res = await api.getAddresses()
      const addresses = res.data || []
      const defaultAddr = addresses.find(a => a.is_default === 1) || addresses[0]
      this.setData({
        addresses: addresses,
        selectedAddress: defaultAddr || addresses[0]
      })
    } catch (err) {
      console.error(err)
    }
  },

  onAddressTap() {
    wx.navigateTo({ url: '/pages/address/address?mode=select' })
  },

  onRemarkInput(e) {
    this.setData({ remark: e.detail.value })
  },

  async onSubmit() {
    if (!this.data.selectedAddress) {
      wx.showToast({ title: '请选择收货地址', icon: 'none' })
      return
    }

    try {
      let res
      if (this.data.cart_ids.length > 0) {
        res = await api.createOrder({
          address_id: this.data.selectedAddress.id,
          cart_ids: this.data.cart_ids,
          remark: this.data.remark
        })
      } else {
        res = await api.createOrder({
          address_id: this.data.selectedAddress.id,
          remark: this.data.remark
        })
      }
      wx.navigateTo({ url: `/pages/pay/pay?order_no=${res.data.order_no}` })
    } catch (err) {
      console.error(err)
    }
  }
})
