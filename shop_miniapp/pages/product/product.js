const api = require('../../utils/request.js')
const image = require('../../utils/image')
const app = getApp()

Page({
  data: {
    id: null,
    product: null,
    quantity: 1,
    selectedAddress: null
  },

  onLoad(options) {
    if (options.id) {
      this.setData({ id: parseInt(options.id) })
      this.loadProduct()
    }
    if (options.cart_ids) {
      this.setData({ cart_ids: options.cart_ids.split(',').map(id => parseInt(id)) })
    }
    this.loadDefaultAddress()
  },

  async loadProduct() {
    try {
      const res = await api.getProduct(this.data.id)
      const product = res.data
      // 格式化图片 URL
      if (product.images) {
        product.images = image.formatImageUrls(product.images)
      }
      this.setData({ product })
    } catch (err) {
      console.error(err)
    }
  },

  async loadDefaultAddress() {
    try {
      const res = await api.getAddresses()
      const addresses = res.data || []
      const defaultAddr = addresses.find(a => a.is_default === 1) || addresses[0]
      this.setData({ selectedAddress: defaultAddr })
    } catch (err) {
      console.error(err)
    }
  },

  onAddressTap() {
    wx.navigateTo({ url: '/pages/address/address?mode=select' })
  },

  onQuantityChange(e) {
    const { value } = e.detail
    this.setData({ quantity: parseInt(value) || 1 })
  },

  onAddCart() {
    if (!app.globalData.token) {
      wx.navigateTo({ url: '/pages/login/login' })
      return
    }
    this.createOrder(true)
  },

  onBuyNow() {
    if (!app.globalData.token) {
      wx.navigateTo({ url: '/pages/login/login' })
      return
    }
    if (!this.data.selectedAddress) {
      wx.showToast({ title: '请选择收货地址', icon: 'none' })
      return
    }
    this.createOrder(false)
  },

  async createOrder(isCart) {
    try {
      let res
      if (isCart) {
        res = await api.addCart({
          product_id: this.data.id,
          quantity: this.data.quantity
        })
        wx.showToast({ title: '已加入购物车', icon: 'success' })
      } else {
        res = await api.createOrder({
          address_id: this.data.selectedAddress.id,
          product_id: this.data.id,
          quantity: this.data.quantity
        })
        wx.navigateTo({ url: `/pages/pay/pay?order_no=${res.data.order_no}` })
      }
    } catch (err) {
      console.error(err)
    }
  }
})
