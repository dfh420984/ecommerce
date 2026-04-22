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

      this.setData({
        banners: banners.data || [],
        categories: categories.data || [],
        recommendProducts: recommend.data || [],
        newProducts: newProducts.data || []
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
  }
})
