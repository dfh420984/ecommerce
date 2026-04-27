const api = require('../../services/api.js')
const image = require('../../utils/image.js')
const config = require('../../utils/config.js')
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
    images: [], // 存储已上传的图片URL
    uploadingImages: [], // 正在上传的图片
    isAnonymous: 0
  },

  async onLoad(options) {
    // 设置导航栏标题
    await config.setNavigationBarTitle('write_review_title', '写评价')
    
    this.setData({
      orderId: options.order_id,
      orderItemId: options.order_item_id,
      productId: options.product_id,
      productName: decodeURIComponent(options.product_name || ''),
      productImage: decodeURIComponent(options.product_image || '')
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
  async onUploadImage() {
    const maxImages = 9
    const currentCount = this.data.images.length + this.data.uploadingImages.length
    
    if (currentCount >= maxImages) {
      wx.showToast({ title: '最多上传9张图片', icon: 'none' })
      return
    }

    wx.chooseMedia({
      count: maxImages - currentCount,
      mediaType: ['image'],
      sizeType: ['compressed'],
      success: async (res) => {
        wx.showLoading({ title: '上传中...' })
        
        const tempFiles = res.tempFiles.map(f => f.tempFilePath)
        const uploadedUrls = []
        
        // 逐张上传图片
        for (const filePath of tempFiles) {
          try {
            const result = await api.uploadImage(filePath)
            // 格式化图片URL
            const imageUrl = image.formatImageUrl(result.data.url)
            uploadedUrls.push(imageUrl)
          } catch (err) {
            console.error('上传图片失败:', err)
            wx.hideLoading()
            wx.showToast({ title: '图片上传失败', icon: 'none' })
            return
          }
        }
        
        wx.hideLoading()
        const images = [...this.data.images, ...uploadedUrls]
        this.setData({ images })
        wx.showToast({ title: `成功上传${uploadedUrls.length}张图片`, icon: 'success' })
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
      await api.createReview({
        order_id: parseInt(this.data.orderId),
        product_id: parseInt(this.data.productId),
        order_item_id: parseInt(this.data.orderItemId),
        rating: this.data.rating,
        content: this.data.content,
        images: this.data.images, // 使用已上传的图片URL
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
