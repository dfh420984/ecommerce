const api = require('../../services/api.js')
const image = require('../../utils/image')
const util = require('../../utils/util.js')
const config = require('../../utils/config.js')
const app = getApp()

Page({
  data: {
    id: null,
    product: null,
    quantity: 1,
    selectedAddress: null,
    contentSegments: [], // 富文本内容片段（包含视频和HTML）
    isLiked: false, // 是否已点赞
    isFavorited: false, // 是否已收藏
    reviewStats: null // 评论统计
  },

  async onLoad(options) {
    // 设置导航栏标题
    await config.setNavigationBarTitle('product_page_title', '商品详情')
    
    if (options.id) {
      this.setData({ id: parseInt(options.id) })
      this.loadProduct()
    }
  },

  onShow() {
    if (app.globalData.token) {
      this.loadDefaultAddress()
    }
  },

  async loadProduct() {
    try {
      const res = await api.getProduct(this.data.id)
      const product = res.data
      // 格式化图片 URL
      if (product.images) {
        product.images = image.formatImageUrls(product.images)
      }
      
      // 解析富文本内容，提取视频
      let contentSegments = []
      if (product.content) {
        console.log('原始富文本内容:', product.content)
        contentSegments = util.parseRichContent(product.content)
        console.log('解析后的内容片段:', contentSegments)
        
        // 格式化视频URL
        contentSegments = contentSegments.map(segment => {
          if (segment.type === 'video') {
            const formattedSegment = {
              ...segment,
              src: image.formatImageUrl(segment.src),
              poster: segment.poster ? image.formatImageUrl(segment.poster) : ''
            }
            console.log('格式化后的视频片段:', formattedSegment)
            return formattedSegment
          }
          return segment
        })
      }
      
      this.setData({ 
        product,
        contentSegments
      })
      
      // 加载点赞和收藏状态
      if (app.globalData.token) {
        this.loadLikeStatus()
        this.loadFavoriteStatus()
      }
      
      // 加载评论统计
      this.loadReviewStats()
    } catch (err) {
      console.error('加载商品详情失败:', err)
    }
  },

  async loadDefaultAddress() {
    try {
      const res = await api.getAddresses()
      const addresses = res.data || []
      const defaultAddr = addresses.find(a => a.is_default === 1) || addresses[0]
      this.setData({ selectedAddress: defaultAddr })
    } catch (err) {
      console.error(err)
    }
  },

  onAddressTap() {
    const currentId = this.data.selectedAddress ? this.data.selectedAddress.id : ''
    wx.navigateTo({ url: `/pages/address/address?mode=select&id=${currentId}` })
  },

  onQuantityChange(e) {
    const { type } = e.currentTarget.dataset
    let quantity = this.data.quantity

    if (type === 'minus') {
      quantity = Math.max(1, quantity - 1)
    } else if (type === 'plus') {
      quantity = quantity + 1
    } else {
      quantity = parseInt(e.detail.value, 10) || 1
    }

    this.setData({ quantity })
  },

  async onAddCart() {
    if (!app.globalData.token) {
      wx.navigateTo({ url: '/pages/login/login' })
      return
    }

    try {
      await api.addCart({
        product_id: this.data.id,
        quantity: this.data.quantity
      })
      wx.showToast({ title: '已加入购物车', icon: 'success' })
    } catch (err) {
      console.error(err)
    }
  },

  onBuyNow() {
    if (!app.globalData.token) {
      wx.navigateTo({ url: '/pages/login/login' })
      return
    }
    wx.navigateTo({
      url: `/pages/order/order?product_id=${this.data.id}&quantity=${this.data.quantity}`
    })
  },

  // 加载点赞状态
  async loadLikeStatus() {
    try {
      const res = await api.checkLikeStatus(this.data.id)
      this.setData({ isLiked: res.data.is_liked })
    } catch (err) {
      console.error('加载点赞状态失败:', err)
    }
  },

  // 加载收藏状态
  async loadFavoriteStatus() {
    try {
      const res = await api.checkFavoriteStatus(this.data.id)
      this.setData({ isFavorited: res.data.is_favorited })
    } catch (err) {
      console.error('加载收藏状态失败:', err)
    }
  },

  // 点赞/取消点赞
  async onLikeTap() {
    if (!app.globalData.token) {
      wx.navigateTo({ url: '/pages/login/login' })
      return
    }

    try {
      if (this.data.isLiked) {
        // 取消点赞
        await api.unlikeProduct(this.data.id)
        this.setData({ 
          isLiked: false,
          'product.like_count': Math.max(0, (this.data.product.like_count || 0) - 1)
        })
        wx.showToast({ title: '已取消点赞', icon: 'none' })
      } else {
        // 点赞
        await api.likeProduct(this.data.id)
        this.setData({ 
          isLiked: true,
          'product.like_count': (this.data.product.like_count || 0) + 1
        })
        wx.showToast({ title: '点赞成功', icon: 'success' })
      }
    } catch (err) {
      console.error('点赞操作失败:', err)
      wx.showToast({ title: err.msg || '操作失败', icon: 'none' })
    }
  },

  // 收藏/取消收藏
  async onFavoriteTap() {
    if (!app.globalData.token) {
      wx.navigateTo({ url: '/pages/login/login' })
      return
    }

    try {
      if (this.data.isFavorited) {
        // 取消收藏
        await api.unfavoriteProduct(this.data.id)
        this.setData({ isFavorited: false })
        wx.showToast({ title: '已取消收藏', icon: 'none' })
      } else {
        // 收藏
        await api.favoriteProduct(this.data.id)
        this.setData({ isFavorited: true })
        wx.showToast({ title: '收藏成功', icon: 'success' })
      }
    } catch (err) {
      console.error('收藏操作失败:', err)
      wx.showToast({ title: err.msg || '操作失败', icon: 'none' })
    }
  },

  // 加载评论统计
  async loadReviewStats() {
    try {
      console.log('开始加载评论统计，商品ID:', this.data.id)
      const res = await api.getReviewStats(this.data.id)
      console.log('评论统计结果:', res.data)
      this.setData({ reviewStats: res.data })
    } catch (err) {
      console.error('加载评论统计失败:', err)
    }
  },

  // 查看评论
  onViewReviews() {
    console.log('点击查看评论，商品ID:', this.data.id)
    wx.navigateTo({
      url: `/pages/reviews/reviews?product_id=${this.data.id}`
    })
  }
})
