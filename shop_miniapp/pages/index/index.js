const api = require('../../utils/request')
const image = require('../../utils/image')
const app = getApp()

Page({
  data: {
    banners: [],
    categories: [],
    recommendProducts: [],
    newProducts: []
  },

  onLoad() {
    this.loadData()
  },

  onShow() {
    if (app.globalData.token) {
      this.updateCartBadge()
    }
  },

  async loadData() {
    try {
      const [banners, categories, recommend, newProducts] = await Promise.all([
        api.getBanners(),
        api.getCategories(),
        api.getRecommendProducts(),
        api.getNewProducts()
      ])

      // 格式化图片 URL
      const formattedBanners = (banners.data || []).map(item => ({
        ...item,
        image: image.formatImageUrl(item.image)
      }))

      const formattedCategories = (categories.data || []).map(item => ({
        ...item,
        icon: image.formatImageUrl(item.icon)
      }))

      const formattedRecommend = (recommend.data || []).map(item => ({
        ...item,
        images: image.formatImageUrls(item.images)
      }))

      const formattedNewProducts = (newProducts.data || []).map(item => ({
        ...item,
        images: image.formatImageUrls(item.images)
      }))

      this.setData({
        banners: formattedBanners,
        categories: formattedCategories,
        recommendProducts: formattedRecommend,
        newProducts: formattedNewProducts
      })
    } catch (err) {
      console.error(err)
    }
  },

  async updateCartBadge() {
    try {
      const res = await api.getCartCount()
      const count = res.data.count
      if (count > 0) {
        wx.setTabBarBadge({ index: 2, text: String(count) })
      } else {
        wx.removeTabBarBadge({ index: 2 })
      }
    } catch (err) {
      console.error(err)
    }
  },

  onBannerTap(e) {
    const { link, link_type: linkType, target_id: targetId } = e.currentTarget.dataset
    if (linkType === 1 && targetId) {
      wx.navigateTo({ url: `/pages/product/product?id=${targetId}` })
    } else if (link) {
      wx.navigateTo({ url: link })
    }
  },

  onCategoryTap(e) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({ url: `/pages/product/product?category_id=${id}` })
  },

  onProductTap(e) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({ url: `/pages/product/product?id=${id}` })
  },

  onMoreRecommendTap() {
    wx.navigateTo({ url: '/pages/product-list/product-list?type=recommend' })
  },

  onMoreNewProductsTap() {
    wx.navigateTo({ url: '/pages/product-list/product-list?type=new' })
  }
})
