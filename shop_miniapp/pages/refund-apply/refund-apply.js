// pages/refund-apply/refund-apply.js
const api = require('../../services/api')

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
      this.setData({ orderId: options.order_id })
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
        that.setData({
          images: [...that.data.images, ...tempFiles]
        })
      }
    })
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

    // TODO: 上传图片到服务器，获取URL
    // 这里简化处理，直接使用本地路径（实际应该先上传）
    const imageUrls = images // 实际应该调用上传接口

    api.applyRefundNew({
      order_id: orderId,
      refund_type: refundType,
      reason: reason,
      refund_amount: refundAmount,
      images: imageUrls
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
