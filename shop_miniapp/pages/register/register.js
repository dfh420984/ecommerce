const api = require('../../services/api.js')
const config = require('../../utils/config.js')
const app = getApp()

Page({
  data: {
    username: '',
    password: '',
    confirmPassword: '',
    phone: ''
  },

  async onLoad() {
    // 设置导航栏标题
    await config.setNavigationBarTitle('register_page_title', '注册')
  },

  onUsernameInput(e) {
    this.setData({ username: e.detail.value })
  },

  onPasswordInput(e) {
    this.setData({ password: e.detail.value })
  },

  onConfirmPasswordInput(e) {
    this.setData({ confirmPassword: e.detail.value })
  },

  onPhoneInput(e) {
    this.setData({ phone: e.detail.value })
  },

  async onRegister() {
    const { username, password, confirmPassword, phone } = this.data

    // 表单验证
    if (!username) {
      wx.showToast({ title: '请输入用户名', icon: 'none' })
      return
    }

    if (username.length < 3 || username.length > 50) {
      wx.showToast({ title: '用户名长度为3-50个字符', icon: 'none' })
      return
    }

    if (!password) {
      wx.showToast({ title: '请输入密码', icon: 'none' })
      return
    }

    if (password.length < 6 || password.length > 50) {
      wx.showToast({ title: '密码长度为6-50个字符', icon: 'none' })
      return
    }

    if (password !== confirmPassword) {
      wx.showToast({ title: '两次输入的密码不一致', icon: 'none' })
      return
    }

    if (phone && !/^1[3-9]\d{9}$/.test(phone)) {
      wx.showToast({ title: '请输入正确的手机号', icon: 'none' })
      return
    }

    try {
      wx.showLoading({ title: '注册中...' })
      
      const registerData = { username, password }
      if (phone) {
        registerData.phone = phone
      }
      
      await api.register(registerData)
      
      wx.hideLoading()
      wx.showToast({ title: '注册成功', icon: 'success' })
      
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
    } catch (err) {
      wx.hideLoading()
      console.error('注册失败:', err)
    }
  },

  onGoLogin() {
    wx.navigateBack()
  }
})
