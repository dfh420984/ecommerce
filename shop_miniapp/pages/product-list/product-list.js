const api = require('../../services/api.js')
const image = require('../../utils/image')

Page({
  data: {
    title: '',
    type: '', // 'recommend' 或 'new' 或 'search'
    keyword: '', // 搜索关键词
    products: [],
    page: 1,
    pageSize: 10,
    hasMore: true,
    loading: false,
    searchTimer: null // 搜索防抖定时器
  },

  onLoad(options) {
    const type = options.type || 'recommend'
    let title = type === 'recommend' ? '推荐商品' : (type === 'new' ? '新品上市' : '搜索商品')
    
    this.setData({ 
      type,
      title
    })
    
    this.loadProducts()
  },

  async loadProducts() {
    if (this.data.loading || !this.data.hasMore) return

    this.setData({ loading: true })

    try {
      let res
      // 根据类型加载不同的商品
      if (this.data.type === 'search') {
        // 搜索商品
        const params = {
          page: this.data.page,
          page_size: this.data.pageSize,
          keyword: this.data.keyword
        }
        res = await api.getProducts(params)
      } else if (this.data.type === 'recommend') {
        res = await api.getRecommendProducts()
      } else if (this.data.type === 'new') {
        res = await api.getNewProducts()
      }

      // 处理搜索结果的分页
      if (this.data.type === 'search') {
        const list = (res.data.list || []).map(item => ({
          ...item,
          images: image.formatImageUrls(item.images)
        }))

        this.setData({
          products: this.data.page === 1 ? list : [...this.data.products, ...list],
          hasMore: list.length >= this.data.pageSize,
          loading: false
        })
      } else {
        // 模拟分页数据
        const allProducts = res.data || []
        const start = (this.data.page - 1) * this.data.pageSize
        const end = start + this.data.pageSize
        const list = allProducts.slice(start, end).map(item => ({
          ...item,
          images: image.formatImageUrls(item.images)
        }))

        this.setData({
          products: this.data.page === 1 ? list : [...this.data.products, ...list],
          hasMore: end < allProducts.length,
          loading: false
        })
      }
    } catch (err) {
      console.error(err)
      this.setData({ loading: false })
    }
  },

  onProductTap(e) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({ url: `/pages/product/product?id=${id}` })
  },

  // 搜索输入
  onSearchInput(e) {
    const keyword = e.detail.value
    this.setData({ keyword })
    
    // 清除之前的定时器
    if (this.data.searchTimer) {
      clearTimeout(this.data.searchTimer)
    }
    
    // 防抖处理：500ms后自动搜索
    const timer = setTimeout(() => {
      if (keyword.trim()) {
        this.setData({
          products: [],
          page: 1,
          hasMore: true
        })
        this.loadProducts()
      }
    }, 500)
    
    this.setData({ searchTimer: timer })
  },

  // 确认搜索（点击键盘搜索按钮）
  onSearchConfirm(e) {
    const keyword = e.detail.value.trim()
    if (!keyword) {
      wx.showToast({ title: '请输入搜索关键词', icon: 'none' })
      return
    }
    
    // 清除定时器
    if (this.data.searchTimer) {
      clearTimeout(this.data.searchTimer)
    }
    
    this.setData({
      keyword,
      products: [],
      page: 1,
      hasMore: true
    })
    this.loadProducts()
  },

  // 清除搜索
  onClearSearch() {
    this.setData({
      keyword: '',
      products: [],
      page: 1,
      hasMore: true
    })
  },

  onReachBottom() {
    if (this.data.hasMore) {
      this.setData({ page: this.data.page + 1 })
      this.loadProducts()
    }
  }
})
