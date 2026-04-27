const config = require('../../utils/config')
const configUtil = require('../../utils/config.js')
const request = require('../../utils/request')

Page({
  data: {
    searchKeyword: '',
    currentCategory: 0,
    expandedQuestionId: null,
    serviceTime: '09:00-21:00',
    customerServicePhone: '400-123-4567',
    categories: [
      { id: 0, name: '全部' }
    ],
    questions: [],
    filteredQuestions: []
  },

  async onLoad() {
    // 设置导航栏标题
    await configUtil.setNavigationBarTitle('help_page_title', '帮助中心')
    
    this.loadHelpConfig()
    this.loadCategories()
    this.loadQuestions()
  },

  // 加载帮助中心配置
  async loadHelpConfig() {
    try {
      const res = await request.get('/miniapp/help/config')
      if (res.code === 0 && res.data) {
        this.setData({
          serviceTime: res.data.service_time || '09:00-21:00',
          customerServicePhone: res.data.customer_service_phone || '400-123-4567'
        })
      }
    } catch (err) {
      console.error('加载配置失败', err)
    }
  },

  // 加载分类列表
  async loadCategories() {
    try {
      const res = await request.get('/miniapp/help/categories')
      if (res.code === 0 && res.data) {
        const categories = [{ id: 0, name: '全部' }, ...res.data]
        this.setData({ categories })
      }
    } catch (err) {
      console.error('加载分类失败', err)
    }
  },

  // 加载问题列表
  async loadQuestions(categoryId = 0) {
    try {
      let url = '/miniapp/help/questions/0'
      if (categoryId > 0) {
        url = `/miniapp/help/questions/${categoryId}`
      }
      
      const res = await request.get(url)
      if (res.code === 0 && res.data) {
        this.setData({ questions: res.data })
        this.filterQuestions()
      }
    } catch (err) {
      console.error('加载问题失败', err)
    }
  },

  onSearchInput(e) {
    const keyword = e.detail.value
    this.setData({
      searchKeyword: keyword
    })
    
    // 如果有搜索关键词，使用后端搜索 API
    if (keyword && keyword.trim()) {
      this.searchQuestions(keyword.trim())
    } else {
      // 清空搜索时，重新加载当前分类的问题
      const category = this.data.categories[this.data.currentCategory]
      this.loadQuestions(category.id)
    }
  },

  // 搜索问题
  async searchQuestions(keyword) {
    try {
      const category = this.data.categories[this.data.currentCategory]
      let url = `/miniapp/help/search?keyword=${encodeURIComponent(keyword)}`
      
      // 如果选择了分类，添加分类过滤
      if (category.id > 0) {
        url += `&category_id=${category.id}`
      }
      
      const res = await request.get(url)
      if (res.code === 0 && res.data) {
        this.setData({
          questions: res.data,
          filteredQuestions: res.data
        })
      }
    } catch (err) {
      console.error('搜索失败', err)
    }
  },

  onCategoryTap(e) {
    const index = e.currentTarget.dataset.index
    const category = this.data.categories[index]
    
    this.setData({
      currentCategory: index,
      expandedQuestionId: null
    })
    
    // 根据分类加载问题
    this.loadQuestions(category.id)
  },

  filterQuestions() {
    const { searchKeyword, questions } = this.data
    
    let filtered = questions
    
    // 按关键词搜索
    if (searchKeyword) {
      const keyword = searchKeyword.toLowerCase()
      filtered = filtered.filter(q => 
        q.title.toLowerCase().includes(keyword) || 
        q.answer.toLowerCase().includes(keyword)
      )
    }
    
    this.setData({
      filteredQuestions: filtered
    })
  },

  onQuestionTap(e) {
    const id = e.currentTarget.dataset.id
    const { expandedQuestionId } = this.data
    
    // 切换展开/收起
    this.setData({
      expandedQuestionId: expandedQuestionId === id ? null : id
    })
  },

  onOnlineService() {
    // TODO: 跳转到在线客服页面或打开客服会话
    wx.showToast({
      title: '在线客服功能开发中',
      icon: 'none'
    })
  },

  onPhoneService() {
    wx.makePhoneCall({
      phoneNumber: this.data.customerServicePhone
    })
  }
})
