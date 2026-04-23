const api = require('../../utils/request.js')
const image = require('../../utils/image')
const app = getApp()

Page({
  data: {
    orderId: null,
    orderNo: null,
    order: null,
    orderItems: [],
    selectedAddress: null,
    cart_ids: [],
    productId: null,
    quantity: 1,
    remark: '',
    totalAmount: '0.00'
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
    if (options.cart_ids) {
      this.setData({ cart_ids: options.cart_ids.split(',').map(id => parseInt(id)) })
      this.loadCartPreview()
    }
    if (options.product_id) {
      this.setData({
        productId: parseInt(options.product_id),
        quantity: parseInt(options.quantity || '1')
      })
      this.loadProductPreview()
    }
    this.loadAddresses()
  },

  onShow() {
    if (app.globalData.token) {
      this.loadAddresses()
    }
  },

  async loadOrder() {
    try {
      const res = await api.getOrder(this.data.orderId)
      const order = res.data
      // 格式化订单商品图片
      if (order.items) {
        order.items = order.items.map(item => ({
          ...item,
          product_image: image.formatImageUrl(item.product_image)
        }))
      }
      this.setData({
        order,
        orderItems: order.items || [],
        totalAmount: Number(order.pay_amount || 0).toFixed(2)
      })
    } catch (err) {
      console.error(err)
    }
  },

  async loadCartPreview() {
    try {
      const res = await api.getCart()
      const cartIds = this.data.cart_ids
      const selectedList = (res.data.list || []).filter(item => cartIds.includes(item.id))
      const orderItems = selectedList.map(item => {
        const images = image.formatImageUrls(item.product.images || [])
        return {
          id: item.id,
          product_id: item.product_id,
          product_name: item.product.name,
          product_image: images[0] || '',
          price: item.product.price,
          quantity: item.quantity,
          subtotal: Number(item.product.price) * Number(item.quantity)
        }
      })
      const totalAmount = orderItems.reduce((sum, item) => sum + item.subtotal, 0)

      this.setData({
        orderItems,
        totalAmount: totalAmount.toFixed(2)
      })
    } catch (err) {
      console.error(err)
    }
  },

  async loadProductPreview() {
    try {
      const res = await api.getProduct(this.data.productId)
      const product = res.data
      const images = image.formatImageUrls(product.images || [])
      const quantity = this.data.quantity
      const subtotal = Number(product.price) * Number(quantity)

      this.setData({
        orderItems: [{
          id: product.id,
          product_id: product.id,
          product_name: product.name,
          product_image: images[0] || '',
          price: product.price,
          quantity,
          subtotal
        }],
        totalAmount: subtotal.toFixed(2)
      })
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
        selectedAddress: defaultAddr || addresses[0]
      })
    } catch (err) {
      console.error(err)
    }
  },

  onAddressTap() {
    const currentId = this.data.selectedAddress ? this.data.selectedAddress.id : ''
    wx.navigateTo({ url: `/pages/address/address?mode=select&id=${currentId}` })
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
      } else if (this.data.productId) {
        res = await api.createOrder({
          address_id: this.data.selectedAddress.id,
          product_id: this.data.productId,
          quantity: this.data.quantity,
          remark: this.data.remark
        })
      } else {
        wx.showToast({ title: '订单商品不能为空', icon: 'none' })
        return
      }
      wx.navigateTo({ url: `/pages/pay/pay?order_id=${res.data.order_id}&order_no=${res.data.order_no}` })
    } catch (err) {
      console.error(err)
    }
  }
})
