App({
  globalData: {
    apiBase: 'http://localhost:8686/api',
    imageUrl: 'http://localhost:8686',  // 图片基础 URL
    userInfo: null,
    token: null,
    categoryType: '' // 分类页面类型：'recommend' 或 'new' 或空
  },

  onLaunch() {
    const token = wx.getStorageSync('token')
    const userInfo = wx.getStorageSync('userInfo')
    if (token) {
      this.globalData.token = token
      this.globalData.userInfo = userInfo
    }
  },

  setToken(token) {
    this.globalData.token = token
    wx.setStorageSync('token', token)
  },

  setUserInfo(userInfo) {
    this.globalData.userInfo = userInfo
    wx.setStorageSync('userInfo', userInfo)
  },

  clearSession() {
    this.globalData.token = null
    this.globalData.userInfo = null
    wx.removeStorageSync('token')
    wx.removeStorageSync('userInfo')
  },

  // 检查登录状态，未登录则跳转登录页
  checkLogin() {
    if (!this.globalData.token) {
      wx.navigateTo({ url: '/pages/login/login' })
      return false
    }
    return true
  },

  // 退出登录
  logout() {
    this.clearSession()
  }
})
