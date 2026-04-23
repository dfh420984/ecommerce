Page({
  data: {
    activeTab: 0,
    couponList: []
  },

  onLoad() {
    this.loadCoupons()
  },

  onTabChange(e) {
    const { index } = e.currentTarget.dataset
    this.setData({ activeTab: index })
    this.loadCoupons()
  },

  async loadCoupons() {
    // TODO: 调用后端 API 获取优惠券列表
    // 暂时使用模拟数据
    const mockCoupons = [
      {
        id: 1,
        name: '新人专享优惠券',
        amount: 20,
        minAmount: 100,
        startTime: '2024-01-01',
        endTime: '2024-12-31',
        status: 'unused',
        statusText: '未使用'
      },
      {
        id: 2,
        name: '满200减50',
        amount: 50,
        minAmount: 200,
        startTime: '2024-01-01',
        endTime: '2024-12-31',
        status: 'used',
        statusText: '已使用'
      }
    ]

    this.setData({ couponList: mockCoupons })
  },

  onUseCoupon(e) {
    const { id } = e.currentTarget.dataset
    wx.showToast({ title: '功能开发中', icon: 'none' })
  }
})
