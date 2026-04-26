const api = require('../../services/api.js')
const image = require('../../utils/image')
const app = getApp()

Page({
  data: {
    productId: null, // 商品ID，如果有则显示该商品的评论
    reviews: [],
    page: 1,
    pageSize: 10,
    hasMore: true,
    loading: false
  },

  onLoad(options) {
    if (options.product_id) {
      this.setData({ productId: parseInt(options.product_id) })
      // 设置页面标题为商品评价
      wx.setNavigationBarTitle({ title: '商品评价' })
    } else {
      // 设置页面标题为我的评价
      wx.setNavigationBarTitle({ title: '我的评价' })
    }
    this.loadReviews()
  },

  onShow() {
    if (app.globalData.token || this.data.productId) {
      this.setData({ reviews: [], page: 1, hasMore: true })
      this.loadReviews()
    }
  },

  async loadReviews() {
    if (this.data.loading || !this.data.hasMore) return

    this.setData({ loading: true })

    try {
      let res
      
      // 如果有 productId，加载该商品的评论
      if (this.data.productId) {
        res = await api.getProductReviews(this.data.productId, {
          page: this.data.page,
          page_size: this.data.pageSize
        })
      } else {
        // 否则加载我的评价
        res = await api.getMyReviews({
          page: this.data.page,
          page_size: this.data.pageSize
        })
      }

      const list = (res.data.list || []).map(review => ({
        ...review,
        images: (review.images || []).map(img => image.formatImageUrl(img)),
        product: review.product ? {
          ...review.product,
          images: (review.product.images || []).map(img => image.formatImageUrl(img))
        } : null
      }))

      this.setData({
        reviews: this.data.page === 1 ? list : [...this.data.reviews, ...list],
        hasMore: list.length >= this.data.pageSize,
        loading: false
      })
    } catch (err) {
      console.error(err)
      this.setData({ loading: false })
    }
  },

  onReachBottom() {
    if (this.data.hasMore) {
      this.setData({ page: this.data.page + 1 })
      this.loadReviews()
    }
  },

  onProductTap(e) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({ url: `/pages/product/product?id=${id}` })
  },

  // 获取评分星星
  getStars(rating) {
    return Array.from({ length: 5 }, (_, i) => i < rating ? '★' : '☆').join('')
  }
})
