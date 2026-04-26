const api = require('../../utils/request.js')
const image = require('../../utils/image')
const app = getApp()

Page({
  data: {
    status: 0,
    type: '', // 'refund' 表示退款/售后
    orders: [],
    page: 1,
    pageSize: 10,
    hasMore: true,
    loading: false
  },

  onLoad(options) {
    if (options.status) {
      this.setData({ status: parseInt(options.status) })
    }
    if (options.type) {
      this.setData({ type: options.type })
    }
  },

  onShow() {
    if (app.globalData.token) {
      this.setData({ orders: [], page: 1, hasMore: true })
      this.loadOrders()
    }
  },

  async loadOrders() {
    if (this.data.loading || !this.data.hasMore) return

    this.setData({ loading: true })

    try {
      let res
      
      // 如果是退款/售后类型，需要特殊处理
      if (this.data.type === 'refund') {
        // 获取所有订单，然后在前端过滤
        res = await api.getOrders({
          page: this.data.page,
          page_size: 100 // 获取更多数据以便过滤
        })
        
        const allOrders = (res.data.list || []).map(order => ({
          ...order,
          items: (order.items || []).map(item => ({
            ...item,
            product_image: image.formatImageUrl(item.product_image)
          }))
        }))
        
        // 过滤出退款相关的订单（状态7或8）
        const refundOrders = allOrders.filter(order => 
          order.order_status === 7 || order.order_status === 8
        )
        
        this.setData({
          orders: this.data.page === 1 ? refundOrders : [...this.data.orders, ...refundOrders],
          hasMore: false, // 退款订单通常不多，一次加载完
          loading: false
        })
      } else {
        // 正常订单列表
        res = await api.getOrders({
          page: this.data.page,
          page_size: this.data.pageSize,
          status: this.data.status || undefined
        })

        const list = (res.data.list || []).map(order => ({
          ...order,
          items: (order.items || []).map(item => ({
            ...item,
            product_image: image.formatImageUrl(item.product_image)
          }))
        }))
        
        this.setData({
          orders: this.data.page === 1 ? list : [...this.data.orders, ...list],
          hasMore: list.length >= this.data.pageSize,
          loading: false
        })
      }
    } catch (err) {
      console.error(err)
      this.setData({ loading: false })
    }
  },

  onReachBottom() {
    if (this.data.hasMore) {
      this.setData({ page: this.data.page + 1 })
      this.loadOrders()
    }
  },

  onTabChange(e) {
    const dataset = e.currentTarget.dataset
    
    if (dataset.type === 'refund') {
      // 切换到退款/售后
      this.setData({
        type: 'refund',
        status: 0,
        orders: [],
        page: 1,
        hasMore: true
      })
    } else {
      // 切换到正常订单状态
      this.setData({
        type: '',
        status: parseInt(dataset.status),
        orders: [],
        page: 1,
        hasMore: true
      })
    }
    this.loadOrders()
  },

  onOrderTap(e) {
    const { orderid } = e.currentTarget.dataset
    wx.navigateTo({ url: `/pages/order-detail/order-detail?id=${orderid}` })
  },

  getStatusText(status) {
    const texts = { 
      1: '待付款', 
      2: '待发货', 
      3: '配送中', 
      4: '已收货', 
      5: '已完成', 
      6: '已取消',
      7: '退款中',
      8: '已退款'
    }
    return texts[status] || '未知'
  },

  getStatusClass(status) {
    const classes = { 
      1: 'status-pending', 
      2: 'status-paid', 
      3: 'status-shipped', 
      4: 'status-received',
      7: 'status-refunding',
      8: 'status-refunded'
    }
    return classes[status] || ''
  },

  // 去评价
  onGoToReview(e) {
    const { orderid, items } = e.currentTarget.dataset
    if (!items || items.length === 0) return
    
    // 跳转到第一个未评价的商品的评价页面
    const firstItem = items[0]
    wx.navigateTo({
      url: `/pages/write-review/write-review?order_id=${orderid}&product_id=${firstItem.product_id}&order_item_id=${firstItem.id}&product_name=${encodeURIComponent(firstItem.product_name)}&product_image=${encodeURIComponent(firstItem.product_image || '')}`
    })
  }
})
