const api = require('../../utils/request.js')
const image = require('../../utils/image')

Page({
  data: {
    id: null,
    order: null
  },

  onLoad(options) {
    if (options.id) {
      this.setData({ id: parseInt(options.id, 10) })
      this.loadOrder()
    }
  },

  onShow() {
    if (this.data.id) {
      this.loadOrder()
    }
  },

  async loadOrder() {
    try {
      const res = await api.getOrder(this.data.id)
      const order = {
        ...res.data,
        items: (res.data.items || []).map(item => ({
          ...item,
          product_image: image.formatImageUrl(item.product_image)
        }))
      }
      this.setData({ order })
    } catch (err) {
      console.error(err)
    }
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
    return texts[status] || '未知状态'
  },

  async onPay() {
    if (!this.data.order) return
    wx.navigateTo({
      url: `/pages/pay/pay?order_id=${this.data.order.id}&order_no=${this.data.order.order_no}`
    })
  },

  async onCancelOrder() {
    wx.showModal({
      title: '提示',
      content: '确定取消该订单吗？',
      success: async (res) => {
        if (!res.confirm) return
        try {
          await api.cancelOrder(this.data.order.id)
          wx.showToast({ title: '已取消', icon: 'success' })
          this.loadOrder()
        } catch (err) {
          console.error(err)
        }
      }
    })
  },

  async onConfirmReceive() {
    wx.showModal({
      title: '提示',
      content: '确认已收到商品？',
      success: async (res) => {
        if (!res.confirm) return
        try {
          await api.confirmReceive(this.data.order.id)
          wx.showToast({ title: '确认成功', icon: 'success' })
          this.loadOrder()
        } catch (err) {
          console.error(err)
        }
      }
    })
  },

  async onDeleteOrder() {
    wx.showModal({
      title: '提示',
      content: '确定删除该订单吗？',
      success: async (res) => {
        if (!res.confirm) return
        try {
          await api.deleteOrder(this.data.order.id)
          wx.showToast({ title: '删除成功', icon: 'success' })
          setTimeout(() => {
            wx.navigateBack()
          }, 1200)
        } catch (err) {
          console.error(err)
        }
      }
    })
  },

  onProductTap(e) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({ url: `/pages/product/product?id=${id}` })
  },

  // 查看物流
  onViewLogistics() {
    if (!this.data.order || !this.data.order.id) return
    wx.navigateTo({
      url: `/pages/logistics/logistics?order_id=${this.data.order.id}`
    })
  },

  // 申请退款
  onApplyRefund() {
    if (!this.data.order || !this.data.order.id) return
    wx.navigateTo({
      url: `/pages/refund-apply/refund-apply?order_id=${this.data.order.id}`
    })
  }
})
