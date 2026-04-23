const api = require('../../services/api.js')
const app = getApp()

Page({
  data: {
    username: '',
    password: ''
  },

  onLoad() {},

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

  onWechatLogin() {
    wx.showToast({ title: '微信登录开发中', icon: 'none' })
  },

  onGoRegister() {
    wx.navigateTo({ url: '/pages/register/register' })
  }
})
