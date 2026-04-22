const api = require('../../utils/request.js')
const image = require('../../utils/image')
const app = getApp()

Page({
  data: {
    status: 0,
    orders: [],
    page: 1,
    pageSize: 10,
    hasMore: true,
    loading: false
  },

  onLoad(options) {
    if (options.status) {
      this.setData({ status: parseInt(options.status) })
    }
  },

  onShow() {
    if (app.globalData.token) {
      this.setData({ orders: [], page: 1, hasMore: true })
      this.loadOrders()
    }
  },

  async loadOrders() {
    if (this.data.loading || !this.data.hasMore) return

    this.setData({ loading: true })

    try {
      const res = await api.getOrders({
        page: this.data.page,
        page_size: this.data.pageSize,
        status: this.data.status || undefined
      })

      const list = (res.data.list || []).map(order => ({
        ...order,
        items: (order.items || []).map(item => ({
          ...item,
          product_image: image.formatImageUrl(item.product_image)
        }))
      }))
      this.setData({
        orders: this.data.page === 1 ? list : [...this.data.orders, ...list],
        hasMore: list.length >= this.data.pageSize,
        loading: false
      })
    } catch (err) {
      console.error(err)
      this.setData({ loading: false })
    }
  },

  onReachBottom() {
    if (this.data.hasMore) {
      this.setData({ page: this.data.page + 1 })
      this.loadOrders()
    }
  },

  onOrderTap(e) {
    const { orderno } = e.currentTarget.dataset
    wx.navigateTo({ url: `/pages/order-list/order-list?order_no=${orderno}` })
  },

  getStatusText(status) {
    const texts = { 1: '待付款', 2: '待发货', 3: '配送中', 4: '已收货', 5: '已完成', 6: '已取消' }
    return texts[status] || '未知'
  },

  getStatusClass(status) {
    const classes = { 1: 'status-pending', 2: 'status-paid', 3: 'status-shipped', 4: 'status-received' }
    return classes[status] || ''
  }
})
