const api = require('../../services/api.js')
const image = require('../../utils/image')
const app = getApp()

Page({
  data: {
    userInfo: null,
    isLoggedIn: false,
    menuItems: [
      [
        { icon: '/static/user/profile.png', text: '个人资料', url: '/pages/profile/profile' },
        { icon: '/static/user/order.png', text: '全部订单', url: '/pages/order-list/order-list' }
      ],
      [
        { icon: '/static/user/address.png', text: '收货地址', url: '/pages/address/address' },
        { icon: '/static/user/coupon.png', text: '优惠券', url: '/pages/coupon/coupon' }
      ],
      [
        { icon: '/static/user/help.png', text: '帮助中心', url: '/pages/help/help' },
        { icon: '/static/user/about.png', text: '关于我们', url: '/pages/about/about' }
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
    this.checkLoginStatus()
  },

  onShow() {
    this.checkLoginStatus()
    if (app.globalData.token) {
      this.loadUserInfo()
      this.loadOrderCounts()
    }
  },

  // 检查登录状态
  checkLoginStatus() {
    const token = app.globalData.token
    const userInfo = app.globalData.userInfo
    this.setData({
      isLoggedIn: !!token,
      userInfo: userInfo || null
    })
  },

  // 加载用户信息
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
      console.error('加载用户信息失败:', err)
    }
  },

  // 加载订单数量统计
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
      console.error('加载订单统计失败:', err)
    }
  },

  // 点击登录
  onLoginTap() {
    if (!app.globalData.token) {
      wx.navigateTo({ url: '/pages/login/login' })
    }
  },

  // 订单状态点击
  onOrderStatusTap(e) {
    if (!app.checkLogin()) return
    const { status } = e.currentTarget.dataset
    wx.navigateTo({ url: `/pages/order-list/order-list?status=${status}` })
  },

  // 菜单项点击
  onMenuTap(e) {
    const { url } = e.currentTarget.dataset
    if (!url) {
      wx.showToast({ title: '功能开发中', icon: 'none' })
      return
    }
    if (!app.checkLogin()) return
    wx.navigateTo({ url })
  },

  // 退出登录
  onLogout() {
    wx.showModal({
      title: '提示',
      content: '确定要退出登录吗？',
      success: (res) => {
        if (res.confirm) {
          app.logout()
          this.setData({
            userInfo: null,
            isLoggedIn: false,
            orderCounts: {
              pending: 0,
              paid: 0,
              shipped: 0,
              received: 0
            }
          })
          wx.showToast({ title: '已退出登录', icon: 'success' })
        }
      }
    })
  }
})
