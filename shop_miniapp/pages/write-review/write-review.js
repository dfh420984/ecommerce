const api = require('../../services/api.js')
const app = getApp()

Page({
  data: {
    orderId: '',
    orderItemId: '',
    productId: '',
    productName: '',
    productImage: '',
    rating: 5,
    content: '',
    images: [],
    isAnonymous: 0
  },

  onLoad(options) {
    this.setData({
      orderId: options.order_id,
      orderItemId: options.order_item_id,
      productId: options.product_id,
      productName: options.product_name,
      productImage: options.product_image
    })
  },

  // 选择评分
  onRatingChange(e) {
    const rating = parseInt(e.currentTarget.dataset.rating)
    this.setData({ rating })
  },

  // 内容输入
  onContentInput(e) {
    this.setData({ content: e.detail.value })
  },

  // 上传图片
  onUploadImage() {
    if (this.data.images.length >= 9) {
      wx.showToast({ title: '最多上传9张图片', icon: 'none' })
      return
    }

    wx.chooseMedia({
      count: 9 - this.data.images.length,
      mediaType: ['image'],
      sizeType: ['compressed'],
      success: (res) => {
        const images = [...this.data.images, ...res.tempFiles.map(f => f.tempFilePath)]
        this.setData({ images })
      }
    })
  },

  // 删除图片
  onDeleteImage(e) {
    const index = e.currentTarget.dataset.index
    const images = this.data.images.filter((_, i) => i !== index)
    this.setData({ images })
  },

  // 匿名评价
  onAnonymousChange(e) {
    this.setData({ isAnonymous: e.detail.value ? 1 : 0 })
  },

  // 提交评价
  async onSubmit() {
    if (!this.data.content.trim()) {
      wx.showToast({ title: '请输入评价内容', icon: 'none' })
      return
    }

    wx.showLoading({ title: '提交中...' })

    try {
      // TODO: 上传图片到服务器获取URL
      // 这里简化处理，直接使用本地路径
      await api.createReview({
        order_id: parseInt(this.data.orderId),
        product_id: parseInt(this.data.productId),
        order_item_id: parseInt(this.data.orderItemId),
        rating: this.data.rating,
        content: this.data.content,
        images: this.data.images,
        is_anonymous: this.data.isAnonymous
      })

      wx.hideLoading()
      wx.showToast({ title: '评价成功', icon: 'success' })

      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
    } catch (err) {
      wx.hideLoading()
      wx.showToast({ title: err.msg || '提交失败', icon: 'none' })
    }
  }
})
