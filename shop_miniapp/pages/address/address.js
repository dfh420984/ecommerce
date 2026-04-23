const api = require('../../utils/request.js')
const app = getApp()

Page({
  data: {
    mode: 'normal',
    addresses: [],
    selectedId: null
  },

  onLoad(options) {
    if (options.mode === 'select') {
      this.setData({ mode: 'select' })
    }
    if (options.id) {
      this.setData({ selectedId: parseInt(options.id) })
    }
  },

  onShow() {
    this.loadAddresses()
  },

  async loadAddresses() {
    try {
      const res = await api.getAddresses()
      this.setData({ addresses: res.data || [] })
    } catch (err) {
      console.error(err)
    }
  },

  onAddressTap(e) {
    const { id } = e.currentTarget.dataset
    const { mode } = this.data

    if (mode === 'select') {
      const selectedAddress = this.data.addresses.find(a => a.id === id)
      this.setData({ selectedId: id })
      const pages = getCurrentPages()
      const prevPage = pages[pages.length - 2]
      if (prevPage) {
        prevPage.setData({ selectedAddress })
      }
      wx.navigateBack()
    } else {
      wx.navigateTo({ url: `/pages/address-form/address-form?id=${id}` })
    }
  },

  onAddAddress() {
    wx.navigateTo({ url: '/pages/address-form/address-form' })
  },

  async onSetDefault(e) {
    const { id } = e.currentTarget.dataset
    try {
      await api.setDefaultAddress(id)
      this.loadAddresses()
    } catch (err) {
      console.error(err)
    }
  },

  async onDeleteAddress(e) {
    const { id } = e.currentTarget.dataset
    wx.showModal({
      title: '提示',
      content: '确定删除该地址?',
      success: async (res) => {
        if (res.confirm) {
          try {
            await api.deleteAddress(id)
            this.loadAddresses()
          } catch (err) {
            console.error(err)
          }
        }
      }
    })
  }
})
