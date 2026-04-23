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
  },

  onShow() {
    if (app.globalData.token) {
      this.loadDefaultAddress()
    }
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
    const currentId = this.data.selectedAddress ? this.data.selectedAddress.id : ''
    wx.navigateTo({ url: `/pages/address/address?mode=select&id=${currentId}` })
  },

  onQuantityChange(e) {
    const { type } = e.currentTarget.dataset
    let quantity = this.data.quantity

    if (type === 'minus') {
      quantity = Math.max(1, quantity - 1)
    } else if (type === 'plus') {
      quantity = quantity + 1
    } else {
      quantity = parseInt(e.detail.value, 10) || 1
    }

    this.setData({ quantity })
  },

  async onAddCart() {
    if (!app.globalData.token) {
      wx.navigateTo({ url: '/pages/login/login' })
      return
    }

    try {
      await api.addCart({
        product_id: this.data.id,
        quantity: this.data.quantity
      })
      wx.showToast({ title: '已加入购物车', icon: 'success' })
    } catch (err) {
      console.error(err)
    }
  },

  onBuyNow() {
    if (!app.globalData.token) {
      wx.navigateTo({ url: '/pages/login/login' })
      return
    }
    wx.navigateTo({
      url: `/pages/order/order?product_id=${this.data.id}&quantity=${this.data.quantity}`
    })
  }
})
