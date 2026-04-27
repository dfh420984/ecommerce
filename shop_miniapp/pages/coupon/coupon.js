const api = require('../../services/api')
const config = require('../../utils/config.js')

Page({
  data: {
    activeTab: 0, // 0-可领取, 1-已领取
    subTab: 0, // 0-未使用, 1-已使用, 2-已过期
    availableCoupons: [], // 可领取的优惠券
    myCoupons: [], // 我的优惠券
    loading: false
  },

  async onLoad() {
    // 设置导航栏标题
    await config.setNavigationBarTitle('coupon_page_title', '优惠券')
    this.loadAvailableCoupons()
  },

  onShow() {
    // 每次显示时刷新数据
    if (this.data.activeTab === 0) {
      this.loadAvailableCoupons()
    } else {
      this.loadMyCoupons()
    }
  },

  onTabChange(e) {
    const index = parseInt(e.currentTarget.dataset.index)
    this.setData({ 
      activeTab: index,
      subTab: 0
    })
    if (index === 0) {
      this.loadAvailableCoupons()
    } else {
      this.loadMyCoupons()
    }
  },

  onSubTabChange(e) {
    const index = parseInt(e.currentTarget.dataset.index)
    this.setData({ subTab: index })
    this.loadMyCoupons()
  },

  // 加载可领取的优惠券
  async loadAvailableCoupons() {
    this.setData({ loading: true })
    try {
      const res = await api.getAvailableCoupons()
      this.setData({ 
        availableCoupons: res.data || []
      })
    } catch (error) {
      console.error('加载可领取优惠券失败:', error)
      wx.showToast({ title: '加载失败', icon: 'none' })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 加载我的优惠券
  async loadMyCoupons() {
    this.setData({ loading: true })
    try {
      // status: 1-未使用, 2-已使用, 3-已过期
      const statusMap = ['1', '2', '3']
      const status = statusMap[this.data.subTab]
      const res = await api.getMyCoupons(status)
      
      // 处理数据，添加格式化后的时间字段
      const processedCoupons = (res.data || []).map(coupon => {
        return {
          ...coupon,
          // 预先格式化时间，避免 WXML 中调用函数的问题
          expireTimeText: this.formatExpireTime(coupon.expire_time)
        }
      })
      
      console.log('处理后的优惠券数据:', processedCoupons[0])
      
      this.setData({ 
        myCoupons: processedCoupons
      })
    } catch (error) {
      console.error('加载我的优惠券失败:', error)
      wx.showToast({ title: '加载失败', icon: 'none' })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 领取优惠券
  async onReceiveCoupon(e) {
    const { id } = e.currentTarget.dataset
    
    wx.showLoading({ title: '领取中...' })
    try {
      await api.receiveCoupon(id)
      wx.hideLoading()
      wx.showToast({ title: '领取成功', icon: 'success' })
      // 刷新列表
      this.loadAvailableCoupons()
    } catch (error) {
      wx.hideLoading()
      console.error('领取优惠券失败:', error)
      wx.showToast({ title: error.message || '领取失败', icon: 'none' })
    }
  },

  // 去使用
  onUseCoupon(e) {
    const { id } = e.currentTarget.dataset
    wx.switchTab({
      url: '/pages/index/index'
    })
  },

  // 格式化过期时间（用于数据处理）
  formatExpireTime(timeStr) {
    if (!timeStr || typeof timeStr !== 'string') {
      return '未知'
    }
    // 提取日期部分 YYYY-MM-DD
    return timeStr.substring(0, 10)
  },

  // 格式化时间（保留用于其他地方）
  formatTime(timeStr) {
    if (!timeStr) return ''
    if (typeof timeStr === 'string') {
      return timeStr.substring(0, 10)
    }
    return ''
  }
})
