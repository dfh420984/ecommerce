const api = require('../../services/api.js')
const config = require('../../utils/config.js')

Page({
  data: {
    orderId: null,
    orderNo: null,
    order: null,
    payTypes: [
      { id: 1, name: '微信支付', icon: '/static/pay/wechat.png' },
      { id: 2, name: '支付宝', icon: '/static/pay/alipay.png' }
    ],
    selectedPayType: 1,
    paying: false,
    paid: false, // 是否已支付
    hasCoupon: false, // 是否使用了优惠券
    couponAmount: 0, // 优惠金额
    payDisabled: false // 支付按钮是否禁用
  },

  async onLoad(options) {
    // 设置导航栏标题
    await config.setNavigationBarTitle('pay_page_title', '支付')
    
    if (options.order_id) {
      this.setData({ orderId: parseInt(options.order_id) })
    }
    if (options.order_no) {
      this.setData({ orderNo: options.order_no })
    }
    if (this.data.orderId) {
      this.loadOrder()
    }
  },

  async loadOrder() {
    try {
      const res = await api.getOrder(this.data.orderId)
      const order = res.data
      this.setData({
        order: order,
        orderNo: order.order_no,
        // 如果订单已支付，更新状态
        paid: order.pay_status === 1,
        // 检查是否使用了优惠券
        hasCoupon: order.coupon_id && order.coupon_amount > 0,
        couponAmount: order.coupon_amount || 0
      })
    } catch (err) {
      console.error(err)
    }
  },

  onPayTypeChange(e) {
    const { id } = e.currentTarget.dataset
    console.log('切换支付方式:', id)
    this.setData({ 
      selectedPayType: parseInt(id)
    }, () => {
      console.log('当前选中支付方式:', this.data.selectedPayType)
    })
  },

  async onPay() {
    // 防止重复提交
    if (this.data.paying || this.data.paid || this.data.payDisabled || !this.data.order) {
      if (this.data.paid) {
        wx.showToast({ title: '订单已支付', icon: 'none' })
      }
      return
    }

    // 设置支付状态
    this.setData({ 
      paying: true,
      payDisabled: true 
    })

    try {
      const res = await api.getPayURL({
        order_id: this.data.order.id,
        pay_type: this.data.selectedPayType
      })

      // 检查是否是模拟支付
      if (res.data.mock_pay) {
        // 模拟支付成功
        wx.showToast({ 
          title: '支付成功', 
          icon: 'success',
          duration: 2000
        })
        
        // 更新页面状态
        this.setData({ 
          paid: true,
          paying: false
        })
        
        // 立即跳转，防止重复点击
        setTimeout(() => {
          wx.redirectTo({ 
            url: `/pages/order-detail/order-detail?id=${this.data.order.id}` 
          })
        }, 1500)
      } else if (res.data.pay_url) {
        // 真实支付，显示二维码或跳转
        wx.showModal({
          title: '提示',
          content: '请扫描支付二维码完成支付',
          showCancel: false,
          success: () => {
            this.setData({ 
              paying: false,
              payDisabled: false 
            })
          }
        })
      } else {
        wx.showToast({ title: '支付开发中', icon: 'none' })
        this.setData({ 
          paying: false,
          payDisabled: false 
        })
      }
    } catch (err) {
      console.error('支付失败:', err)
      const errorMsg = err.msg || err.message || '支付失败'
      
      // 如果是“订单已支付”的错误，标记为已支付
      if (errorMsg.includes('已支付')) {
        this.setData({ 
          paid: true,
          paying: false,
          payDisabled: true
        })
        wx.showToast({ 
          title: '订单已支付', 
          icon: 'success',
          duration: 2000
        })
        setTimeout(() => {
          wx.redirectTo({ 
            url: `/pages/order-detail/order-detail?id=${this.data.order.id}` 
          })
        }, 2000)
      } else {
        wx.showToast({ 
          title: errorMsg, 
          icon: 'none',
          duration: 2000
        })
        // 支付失败，恢复按钮状态
        this.setData({ 
          paying: false,
          payDisabled: false 
        })
      }
    }
  },

  goOrderList() {
    wx.redirectTo({ url: '/pages/order-list/order-list' })
  },

  goHome() {
    wx.switchTab({ url: '/pages/index/index' })
  },

  goOrderDetail() {
    wx.redirectTo({ 
      url: `/pages/order-detail/order-detail?id=${this.data.order.id}` 
    })
  },

  // 页面显示时检查订单状态
  onShow() {
    if (this.data.orderId) {
      this.loadOrder()
    }
  },

  // 页面卸载时重置状态
  onUnload() {
    this.setData({ 
      paying: false,
      payDisabled: false 
    })
  }
})
