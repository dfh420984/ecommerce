const api = require('../../services/api.js')
const config = require('../../utils/config.js')
const app = getApp()

Page({
  data: {
    username: '',
    password: ''
  },

  async onLoad() {
    // 设置导航栏标题
    await config.setNavigationBarTitle('login_page_title', '登录')
  },

  onUsernameInput(e) {
    this.setData({ username: e.detail.value })
  },

  onPasswordInput(e) {
    this.setData({ password: e.detail.value })
  },

  async onLogin() {
    const { username, password } = this.data

    if (!username) {
      wx.showToast({ title: '请输入用户名', icon: 'none' })
      return
    }
    if (!password) {
      wx.showToast({ title: '请输入密码', icon: 'none' })
      return
    }

    try {
      const res = await api.login({ username, password })
      app.setToken(res.data.token)
      app.setUserInfo(res.data.user)
      wx.showToast({ title: '登录成功', icon: 'success' })
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
    } catch (err) {
      console.error(err)
    }
  },

  async onWechatLogin() {
    try {
      // 显示加载提示
      wx.showLoading({ title: '登录中...', mask: true })
      
      // 获取微信登录code
      const loginRes = await new Promise((resolve, reject) => {
        wx.login({
          success: resolve,
          fail: reject
        })
      })
      
      if (!loginRes.code) {
        throw new Error('获取微信登录code失败')
      }
      
      // 调用后端微信登录接口
      const res = await api.wechatLogin({ code: loginRes.code })
      
      // 保存token和用户信息
      app.setToken(res.data.token)
      app.setUserInfo(res.data.user)
      
      wx.hideLoading()
      wx.showToast({ title: '登录成功', icon: 'success' })
      
      // 延迟跳转到上一页或首页
      setTimeout(() => {
        const pages = getCurrentPages()
        if (pages.length > 1) {
          wx.navigateBack()
        } else {
          wx.switchTab({ url: '/pages/index/index' })
        }
      }, 1500)
      
    } catch (err) {
      wx.hideLoading()
      console.error('微信登录失败:', err)
      wx.showToast({ title: '微信登录失败', icon: 'none' })
    }
  },

  onGoRegister() {
    wx.navigateTo({ url: '/pages/register/register' })
  }
})
