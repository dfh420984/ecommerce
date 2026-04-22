const api = require('../../utils/request.js')
const app = getApp()

Page({
  data: {
    cartList: [],
    total: 0,
    allSelected: true,
    isEmpty: false
  },

  onShow() {
    if (app.globalData.token) {
      this.loadCart()
    } else {
      wx.showModal({
        title: '提示',
        content: '请先登录',
        success: (res) => {
          if (res.confirm) {
            wx.navigateTo({ url: '/pages/login/login' })
          }
        }
      })
    }
  },

  async loadCart() {
    try {
      const res = await api.getCart()
      const list = res.data.list || []
      const total = res.data.total || 0
      const allSelected = list.every(item => item.selected === 1)

      this.setData({
        cartList: list,
        total: total,
        allSelected: allSelected,
        isEmpty: list.length === 0
      })
    } catch (err) {
      console.error(err)
    }
  },

  async onSelectItem(e) {
    const { id } = e.currentTarget.dataset
    const { cartList } = this.data
    const item = cartList.find(c => c.id === id)
    if (item) {
      const selected = item.selected === 1 ? 0 : 1
      try {
        await api.selectCart(id, { selected })
        this.loadCart()
      } catch (err) {
        console.error(err)
      }
    }
  },

  async onSelectAll() {
    const selected = this.data.allSelected ? 0 : 1
    try {
      await api.selectAllCart({ selected })
      this.loadCart()
    } catch (err) {
      console.error(err)
    }
  },

  async onQuantityChange(e) {
    const { id } = e.currentTarget.dataset
    const { value } = e.detail
    try {
      await api.updateCart(id, { quantity: value })
      this.loadCart()
    } catch (err) {
      console.error(err)
    }
  },

  async onDeleteItem(e) {
    const { id } = e.currentTarget.dataset
    wx.showModal({
      title: '提示',
      content: '确定删除该商品?',
      success: async (res) => {
        if (res.confirm) {
          try {
            await api.deleteCart(id)
            this.loadCart()
          } catch (err) {
            console.error(err)
          }
        }
      }
    })
  },

  onCheckout() {
    const selectedItems = this.data.cartList.filter(item => item.selected === 1)
    if (selectedItems.length === 0) {
      wx.showToast({ title: '请选择商品', icon: 'none' })
      return
    }
    wx.navigateTo({ url: '/pages/order/order?cart_ids=' + selectedItems.map(i => i.id).join(',') })
  }
})
