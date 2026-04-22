const api = require('../../services/api.js')
const image = require('../../utils/image')
const app = getApp()

Page({
  data: {
    categories: [], // 一级分类列表
    subCategories: [], // 子分类列表
    selectedCategoryId: null, // 当前选中的一级分类ID
    selectedSubCategoryId: null, // 当前选中的子分类ID（如果有）
    products: [],
    page: 1,
    pageSize: 10,
    hasMore: true,
    loading: false,
    hasSubCategories: false // 是否有子分类
  },

  onLoad() {
    this.loadCategories()
    this.loadProducts()
  },

  onShow() {
    // 分类页面只负责分类浏览
  },

  async loadCategories() {
    try {
      const res = await api.getCategories()
      const categories = (res.data || []).map(item => ({
        ...item,
        icon: image.formatImageUrl(item.icon)
      }))
      this.setData({
        categories: categories,
        selectedCategoryId: null, // 默认不选中
        subCategories: [],
        hasSubCategories: false
      })
      
      // 默认选中第一个分类并加载商品
      if (categories.length > 0) {
        const firstCategory = categories[0]
        this.setData({
          selectedCategoryId: firstCategory.id
        })
        await this.loadSubCategoriesAndProducts(firstCategory.id)
      }
    } catch (err) {
      console.error(err)
    }
  },

  async loadSubCategoriesAndProducts(categoryId) {
    try {
      const res = await api.getSubCategories(categoryId)
      const subCategories = (res.data || []).map(item => ({
        ...item,
        icon: image.formatImageUrl(item.icon)
      }))
      
      if (subCategories.length > 0) {
        // 有子分类，展示子分类
        this.setData({
          subCategories: subCategories,
          hasSubCategories: true
        })
      } else {
        // 没有子分类
        this.setData({
          hasSubCategories: false
        })
      }
      
      // 加载父级分类及其所有子分类的商品
      this.loadProducts()
    } catch (err) {
      console.error(err)
      // 如果获取子分类失败，直接查询商品
      this.loadProducts()
    }
  },

  async loadProducts() {
    if (this.data.loading || !this.data.hasMore) return

    this.setData({ loading: true })

    try {
      const params = {
        page: this.data.page,
        page_size: this.data.pageSize
      }

      // 按分类加载商品
      if (this.data.selectedSubCategoryId) {
        // 如果选中了子分类，查询子分类的商品
        params.category_id = this.data.selectedSubCategoryId
      } else if (this.data.selectedCategoryId) {
        // 否则查询父级分类及其子分类的商品
        params.category_id = this.data.selectedCategoryId
      } else {
        // 如果没有选中分类，不加载商品
        this.setData({ 
          products: [],
          loading: false 
        })
        return
      }

      const res = await api.getProducts(params)
      const list = (res.data.list || []).map(item => ({
        ...item,
        images: image.formatImageUrls(item.images)
      }))

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

  async onCategoryChange(e) {
    const { id } = e.currentTarget.dataset
    
    // 如果点击的是已选中的分类，不做处理
    if (this.data.selectedCategoryId === id) {
      return
    }
    
    this.setData({
      selectedCategoryId: id,
      selectedSubCategoryId: null,
      subCategories: [],
      hasSubCategories: false,
      products: [],
      page: 1,
      hasMore: true
    })
    
    await this.loadSubCategoriesAndProducts(id)
  },

  onSubCategoryChange(e) {
    const { id } = e.currentTarget.dataset
    
    // 如果点击的是已选中的子分类，不做处理
    if (this.data.selectedSubCategoryId === id) {
      return
    }
    
    this.setData({
      selectedSubCategoryId: id,
      products: [],
      page: 1,
      hasMore: true
    })
    
    // 加载子分类下的商品
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
