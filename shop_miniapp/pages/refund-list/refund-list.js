// pages/refund-list/refund-list.js
const api = require('../../services/api')

Page({
  data: {
    refundList: [],
    loading: false,
    page: 1,
    pageSize: 10,
    total: 0,
    hasMore: true,
    statusFilter: '' // 空表示全部
  },

  onLoad() {
    this.loadRefunds()
  },

  onPullDownRefresh() {
    this.setData({ page: 1, refundList: [], hasMore: true })
    this.loadRefunds(() => {
      wx.stopPullDownRefresh()
    })
  },

  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.loadMore()
    }
  },

  // 加载退款列表
  loadRefunds(callback) {
    this.setData({ loading: true })
    
    const params = {
      page: this.data.page,
      page_size: this.data.pageSize
    }
    
    if (this.data.statusFilter) {
      params.status = this.data.statusFilter
    }

    api.getMyRefunds(params).then(res => {
      const newList = res.data.list || []
      const total = res.data.total || 0
      
      this.setData({
        refundList: this.data.page === 1 ? newList : [...this.data.refundList, ...newList],
        total,
        hasMore: this.data.refundList.length < total,
        loading: false
      })
      
      if (callback) callback()
    }).catch(err => {
      this.setData({ loading: false })
      wx.showToast({ title: err.message || '加载失败', icon: 'none' })
      if (callback) callback()
    })
  },

  // 加载更多
  loadMore() {
    this.setData({ page: this.data.page + 1 })
    this.loadRefunds()
  },

  // 状态筛选
  onStatusChange(e) {
    const status = e.currentTarget.dataset.status
    this.setData({
      statusFilter: status,
      page: 1,
      refundList: [],
      hasMore: true
    })
    this.loadRefunds()
  },

  // 查看详情
  viewDetail(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/refund-detail/refund-detail?id=${id}`
    })
  },

  // 获取状态文本
  getStatusText(status) {
    const map = {
      1: '待审核',
      2: '已通过',
      3: '已拒绝',
      4: '退款中',
      5: '已退款'
    }
    return map[status] || '未知'
  },

  // 获取状态样式类
  getStatusClass(status) {
    const map = {
      1: 'status-pending',
      2: 'status-approved',
      3: 'status-rejected',
      4: 'status-refunding',
      5: 'status-refunded'
    }
    return map[status] || ''
  }
})
