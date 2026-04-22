const api = require('../../utils/request.js')
const image = require('../../utils/image')
const app = getApp()

Page({
  data: {
    userInfo: null,
    menuItems: [
      [
        { icon: '/static/user/order.png', text: '全部订单', url: '/pages/order-list/order-list' }
      ],
      [
        { icon: '/static/user/address.png', text: '收货地址', url: '/pages/address/address' },
        { icon: '/static/user/coupon.png', text: '优惠券', url: '' }
      ],
      [
        { icon: '/static/user/help.png', text: '帮助中心', url: '' },
        { icon: '/static/user/about.png', text: '关于我们', url: '' }
      ]
    ],
    orderCounts: {
      pending: 0,
      paid: 0,
      shipped: 0,
      received: 0
    }
  },

  onLoad() {
    const userInfo = app.globalData.userInfo
    this.setData({ userInfo })
  },

  onShow() {
    if (app.globalData.token) {
      this.loadUserInfo()
      this.loadOrderCounts()
    }
  },

  async loadUserInfo() {
    try {
      const res = await api.getUserInfo()
      const userInfo = res.data
      // 格式化头像 URL
      if (userInfo.avatar) {
        userInfo.avatar = image.formatImageUrl(userInfo.avatar)
      }
      app.setUserInfo(userInfo)
      this.setData({ userInfo })
    } catch (err) {
      console.error(err)
    }
  },

  async loadOrderCounts() {
    try {
      const res = await api.getOrders({ page: 1, page_size: 100 })
      const orders = res.data.list || []
      const counts = {
        pending: orders.filter(o => o.order_status === 1).length,
        paid: orders.filter(o => o.order_status === 2).length,
        shipped: orders.filter(o => o.order_status === 3).length,
        received: orders.filter(o => o.order_status === 4).length
      }
      this.setData({ orderCounts: counts })
    } catch (err) {
      console.error(err)
    }
  },

  onMenuTap(e) {
    const { url } = e.currentTarget.dataset
    if (!url) {
      wx.showToast({ title: '功能开发中', icon: 'none' })
      return
    }
    if (!app.globalData.token) {
      wx.navigateTo({ url: '/pages/login/login' })
      return
    }
    wx.navigateTo({ url })
  },

  onOrderStatusTap(e) {
    const { status } = e.currentTarget.dataset
    if (!app.globalData.token) {
      wx.navigateTo({ url: '/pages/login/login' })
      return
    }
    wx.navigateTo({ url: `/pages/order-list/order-list?status=${status}` })
  },

  onLoginTap() {
    if (!app.globalData.token) {
      wx.navigateTo({ url: '/pages/login/login' })
    }
  }
})
