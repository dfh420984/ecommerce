const api = require('../../services/api.js')
const image = require('../../utils/image')
const config = require('../../utils/config.js')
const app = getApp()

Page({
  data: {
    favorites: [],
    page: 1,
    pageSize: 10,
    hasMore: true,
    loading: false
  },

  async onLoad() {
    // 设置导航栏标题
    await config.setNavigationBarTitle('favorite_page_title', '我的收藏')
    this.loadFavorites()
  },

  onShow() {
    if (app.globalData.token) {
      this.setData({ favorites: [], page: 1, hasMore: true })
      this.loadFavorites()
    }
  },

  async loadFavorites() {
    if (this.data.loading || !this.data.hasMore) return

    this.setData({ loading: true })

    try {
      const res = await api.getMyFavorites({
        page: this.data.page,
        page_size: this.data.pageSize
      })

      const list = (res.data.list || []).map(item => ({
        ...item,
        product: item.product ? {
          ...item.product,
          images: (item.product.images || []).map(img => image.formatImageUrl(img))
        } : null
      }))

      this.setData({
        favorites: this.data.page === 1 ? list : [...this.data.favorites, ...list],
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
      this.loadFavorites()
    }
  },

  onProductTap(e) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({ url: `/pages/product/product?id=${id}` })
  },

  // 取消收藏
  async onRemoveFavorite(e) {
    const { id } = e.currentTarget.dataset
    
    wx.showModal({
      title: '提示',
      content: '确定要取消收藏吗？',
      success: async (res) => {
        if (res.confirm) {
          try {
            await api.unfavoriteProduct(id)
            wx.showToast({ title: '已取消收藏', icon: 'success' })
            
            // 从列表中移除
            const favorites = this.data.favorites.filter(item => item.product_id !== parseInt(id))
            this.setData({ favorites })
          } catch (err) {
            wx.showToast({ title: '操作失败', icon: 'none' })
          }
        }
      }
    })
  }
})
