const api = require('../../utils/request.js')
const app = getApp()

Page({
  data: {
    categories: [],
    selectedCategoryId: null,
    products: [],
    page: 1,
    pageSize: 10,
    hasMore: true,
    loading: false
  },

  onLoad() {
    this.loadCategories()
    this.loadProducts()
  },

  async loadCategories() {
    try {
      const res = await api.getCategories()
      const categories = res.data || []
      this.setData({
        categories: categories,
        selectedCategoryId: categories.length > 0 ? categories[0].id : null
      })
    } catch (err) {
      console.error(err)
    }
  },

  async loadProducts() {
    if (this.data.loading || !this.data.hasMore) return

    this.setData({ loading: true })

    try {
      const res = await api.getProducts({
        page: this.data.page,
        page_size: this.data.pageSize,
        category_id: this.data.selectedCategoryId || undefined
      })

      const list = res.data.list || []
      this.setData({
        products: this.data.page === 1 ? list : [...this.data.products, ...list],
        hasMore: list.length >= this.data.pageSize,
        loading: false
      })
    } catch (err) {
      console.error(err)
      this.setData({ loading: false })
    }
  },

  onCategoryChange(e) {
    const { id } = e.currentTarget.dataset
    this.setData({
      selectedCategoryId: id,
      products: [],
      page: 1,
      hasMore: true
    })
    this.loadProducts()
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
