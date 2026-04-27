const api = require('../../utils/request.js')
const image = require('../../utils/image')
const config = require('../../utils/config.js')
const app = getApp()

Page({
  data: {
    cartList: [],
    total: 0,
    allSelected: true,
    isEmpty: false
  },

  async onShow() {
    // 设置导航栏标题
    await config.setNavigationBarTitle('cart_page_title', '购物车')
    
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
      const list = (res.data.list || []).map(item => ({
        ...item,
        product: {
          ...item.product,
          images: image.formatImageUrls(item.product.images)
        }
      }))
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
    const { type } = e.currentTarget.dataset
    const currentItem = this.data.cartList.find(item => item.id === id)
    if (!currentItem) return

    let quantity = currentItem.quantity
    if (type === 'minus') {
      quantity = Math.max(1, quantity - 1)
    } else if (type === 'plus') {
      quantity = quantity + 1
    } else {
      quantity = Math.max(1, parseInt(e.detail.value, 10) || 1)
    }

    try {
      await api.updateCart(id, { quantity })
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
  },

  goShopping() {
    wx.switchTab({ url: '/pages/index/index' })
  }
})
