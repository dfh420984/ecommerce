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

      console.log('Banner原始数据:', banners.data)

      // 格式化图片 URL
      const formattedBanners = (banners.data || []).map(item => ({
        ...item,
        image: image.formatImageUrl(item.image)
      }))

      console.log('Banner格式化后:', formattedBanners)

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
    const item = e.currentTarget.dataset.item
    const { link, link_type, target_id } = item
    
    console.log('Banner点击:', { link, link_type, target_id })
    
    // 如果没有配置任何跳转信息，不处理
    if (!link_type && !link && !target_id) {
      return
    }
    
    // link_type === 1: 跳转到商品详情
    if (link_type === 1 && target_id) {
      wx.navigateTo({ 
        url: `/pages/product/product?id=${target_id}` 
      })
      return
    }
    
    // link_type === 2: 跳转到分类页面
    if (link_type === 2 && target_id) {
      app.globalData.selectedCategoryId = target_id
      wx.switchTab({ 
        url: '/pages/category/category' 
      })
      return
    }
    
    // link_type === 3: 跳转到指定页面（内部页面路径）
    if (link_type === 3 && link) {
      // 判断是否是 tabBar 页面
      const tabBarPages = [
        '/pages/index/index',
        '/pages/category/category',
        '/pages/cart/cart',
        '/pages/user/user'
      ]
      
      if (tabBarPages.includes(link)) {
        wx.switchTab({ url: link })
      } else {
        wx.navigateTo({ url: link })
      }
      return
    }
    
    // 有 link 但没有 link_type：直接跳转（支持内部页面和外部链接）
    if (link) {
      // 判断是否是外部链接（http 或 https 开头）
      if (link.startsWith('http://') || link.startsWith('https://')) {
        // 外部链接，使用 web-view 打开
        wx.navigateTo({
          url: `/pages/webview/webview?url=${encodeURIComponent(link)}`
        })
      } else {
        // 内部页面路径
        const tabBarPages = [
          '/pages/index/index',
          '/pages/category/category',
          '/pages/cart/cart',
          '/pages/user/user'
        ]
        
        if (tabBarPages.includes(link)) {
          wx.switchTab({ url: link })
        } else {
          wx.navigateTo({ url: link })
        }
      }
    }
  },

  onCategoryTap(e) {
    const { id } = e.currentTarget.dataset
    // 将选中的分类ID存储到全局变量
    app.globalData.selectedCategoryId = id
    // 使用 switchTab 跳转到分类页面
    wx.switchTab({ url: '/pages/category/category' })
  },

  onSearchTap() {
    // 跳转到商品列表页面，显示搜索框
    wx.navigateTo({ url: '/pages/product-list/product-list?type=search' })
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
