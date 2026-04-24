// pages/refund-apply/refund-apply.js
const api = require('../../services/api')
const app = getApp()
const image = require('../../utils/image')

Page({
  data: {
    orderId: null,
    orderNo: '',
    payAmount: 0,
    refundType: 'refund_only',
    refundAmount: 0,
    reason: '',
    images: [],
    submitting: false
  },

  onLoad(options) {
    if (options.order_id) {
      // 将字符串转换为数字
      this.setData({ orderId: parseInt(options.order_id, 10) })
      this.loadOrderDetail(options.order_id)
    }
  },

  // 加载订单详情
  loadOrderDetail(orderId) {
    wx.showLoading({ title: '加载中' })
    api.getOrder(orderId).then(res => {
      const order = res.data
      this.setData({
        orderNo: order.order_no,
        payAmount: order.pay_amount,
        refundAmount: order.pay_amount
      })
      wx.hideLoading()
    }).catch(err => {
      wx.hideLoading()
      wx.showToast({ title: err.message || '加载失败', icon: 'none' })
    })
  },

  // 退款类型选择
  onRefundTypeChange(e) {
    this.setData({ refundType: e.detail.value })
  },

  // 退款金额输入
  onRefundAmountInput(e) {
    const value = parseFloat(e.detail.value) || 0
    this.setData({ refundAmount: value })
  },

  // 退款原因输入
  onReasonInput(e) {
    this.setData({ reason: e.detail.value })
  },

  // 选择图片
  chooseImage() {
    const that = this
    wx.chooseMedia({
      count: 3 - this.data.images.length,
      mediaType: ['image'],
      sizeType: ['compressed'],
      sourceType: ['album', 'camera'],
      success(res) {
        const tempFiles = res.tempFiles.map(file => file.tempFilePath)
        // 立即上传图片
        that.uploadImages(tempFiles)
      }
    })
  },

  // 上传图片到服务器
  async uploadImages(tempFilePaths) {
    wx.showLoading({ title: '上传中...' })
    
    try {
      const token = app.globalData.token
      const apiBase = app.globalData.apiBase
      const uploadedUrls = []
      
      // 逐个上传图片
      for (const filePath of tempFilePaths) {
        const uploadRes = await new Promise((resolve, reject) => {
          wx.uploadFile({
            url: `${apiBase}/upload`,
            filePath: filePath,
            name: 'file',
            header: {
              'Authorization': `Bearer ${token}`
            },
            success: (res) => {
              const data = JSON.parse(res.data)
              if (data.code === 0) {
                resolve(data)
              } else {
                reject(new Error(data.msg || '上传失败'))
              }
            },
            fail: (err) => {
              reject(err)
            }
          })
        })
        
        uploadedUrls.push(uploadRes.data.url)
      }
      
      wx.hideLoading()
      
      // 将新上传的图片添加到列表
      this.setData({
        images: [...this.data.images, ...uploadedUrls]
      })
      
      wx.showToast({ title: '上传成功', icon: 'success' })
    } catch (err) {
      wx.hideLoading()
      console.error('上传图片失败:', err)
      wx.showToast({ title: '上传失败', icon: 'none' })
    }
  },

  // 删除图片
  deleteImage(e) {
    const index = e.currentTarget.dataset.index
    const images = this.data.images.filter((_, i) => i !== index)
    this.setData({ images })
  },

  // 提交退款申请
  submitRefund() {
    const { orderId, refundType, reason, refundAmount, images, payAmount, submitting } = this.data

    if (!orderId) {
      wx.showToast({ title: '订单信息错误', icon: 'none' })
      return
    }

    if (!reason || reason.trim() === '') {
      wx.showToast({ title: '请填写退款原因', icon: 'none' })
      return
    }

    if (refundAmount <= 0 || refundAmount > payAmount) {
      wx.showToast({ title: '退款金额不正确', icon: 'none' })
      return
    }

    if (submitting) return

    this.setData({ submitting: true })
    wx.showLoading({ title: '提交中' })

    // 确保 order_id 是数字类型
    api.applyRefundNew({
      order_id: parseInt(orderId, 10),
      refund_type: refundType,
      reason: reason,
      refund_amount: parseFloat(refundAmount),
      images: images || []
    }).then(() => {
      wx.hideLoading()
      wx.showToast({ title: '申请成功', icon: 'success' })
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
    }).catch(err => {
      wx.hideLoading()
      this.setData({ submitting: false })
      wx.showToast({ title: err.message || '申请失败', icon: 'none' })
    })
  }
})
