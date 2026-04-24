const api = require('../../services/api.js')
const image = require('../../utils/image')
const config = require('../../utils/config.js')
const app = getApp()

Page({
  data: {
    userInfo: null,
    formData: {
      nickname: '',
      avatar: '',
      phone: '',
      email: ''
    },
    saving: false
  },

  onLoad() {
    this.loadUserInfo()
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
      
      this.setData({
        userInfo: userInfo,
        formData: {
          nickname: userInfo.nickname || '',
          avatar: userInfo.avatar || '',
          phone: userInfo.phone || '',
          email: userInfo.email || ''
        }
      })
    } catch (err) {
      console.error('加载用户信息失败:', err)
      wx.showToast({ title: '加载失败', icon: 'none' })
    }
  },

  // 选择头像
  onChooseAvatar() {
    wx.chooseMedia({
      count: 1,
      mediaType: ['image'],
      sizeType: ['compressed'],
      sourceType: ['album', 'camera'],
      success: (res) => {
        const tempFilePath = res.tempFiles[0].tempFilePath
        this.uploadAvatar(tempFilePath)
      },
      fail: (err) => {
        console.error('选择图片失败:', err)
      }
    })
  },

  // 上传头像
  async uploadAvatar(filePath) {
    wx.showLoading({ title: '上传中...' })
    
    try {
      const token = app.globalData.token
      const baseUrl = config.BASE_URL
      
      const uploadRes = await new Promise((resolve, reject) => {
        wx.uploadFile({
          url: `${baseUrl}/upload/image`,
          filePath: filePath,
          name: 'file',
          header: {
            'Authorization': `Bearer ${token}`
          },
          success: (res) => {
            const data = JSON.parse(res.data)
            if (data.code === 200) {
              resolve(data)
            } else {
              reject(new Error(data.message || '上传失败'))
            }
          },
          fail: (err) => {
            reject(err)
          }
        })
      })
      
      const avatarUrl = uploadRes.data.url
      const formattedUrl = image.formatImageUrl(avatarUrl)
      
      this.setData({
        'formData.avatar': formattedUrl
      })
      
      wx.hideLoading()
      wx.showToast({ title: '上传成功', icon: 'success' })
    } catch (err) {
      wx.hideLoading()
      console.error('上传头像失败:', err)
      wx.showToast({ title: '上传失败', icon: 'none' })
    }
  },

  // 输入框变化
  onInputChange(e) {
    const field = e.currentTarget.dataset.field
    const value = e.detail.value
    this.setData({
      [`formData.${field}`]: value
    })
  },

  // 保存
  async onSave() {
    const { formData } = this.data
    
    // 验证
    if (!formData.nickname || !formData.nickname.trim()) {
      wx.showToast({ title: '请输入昵称', icon: 'none' })
      return
    }

    if (formData.phone && !/^1[3-9]\d{9}$/.test(formData.phone)) {
      wx.showToast({ title: '手机号格式不正确', icon: 'none' })
      return
    }

    if (formData.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      wx.showToast({ title: '邮箱格式不正确', icon: 'none' })
      return
    }

    this.setData({ saving: true })

    try {
      // 获取原始头像路径（去除域名部分）
      let avatarPath = formData.avatar
      if (avatarPath) {
        const baseUrl = config.BASE_URL
        avatarPath = avatarPath.replace(baseUrl, '')
      }

      await api.updateUserInfo({
        nickname: formData.nickname.trim(),
        avatar: avatarPath,
        phone: formData.phone,
        email: formData.email
      })
      
      wx.showToast({ title: '保存成功', icon: 'success' })
      
      // 延迟返回上一页
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
    } catch (err) {
      console.error('保存失败:', err)
      wx.showToast({ title: err.message || '保存失败', icon: 'none' })
    } finally {
      this.setData({ saving: false })
    }
  }
})
