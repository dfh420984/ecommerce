const api = require('../../services/api.js')
const image = require('../../utils/image')

Page({
  data: {
    title: '',
    type: '', // 'recommend' 或 'new'
    products: [],
    page: 1,
    pageSize: 10,
    hasMore: true,
    loading: false
  },

  onLoad(options) {
    const type = options.type || 'recommend'
    const title = type === 'recommend' ? '推荐商品' : '新品上市'
    
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
      if (this.data.type === 'recommend') {
        res = await api.getRecommendProducts()
      } else if (this.data.type === 'new') {
        res = await api.getNewProducts()
      }

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
    } catch (err) {
      console.error(err)
      this.setData({ loading: false })
    }
  },

  onProductTap(e) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({ url: `/pages/product/product?id=${id}` })
  },

  onReachBottom() {
    if (this.data.hasMore) {
      this.setData({ page: this.data.page + 1 })
      this.loadProducts()
    }
  }
})
