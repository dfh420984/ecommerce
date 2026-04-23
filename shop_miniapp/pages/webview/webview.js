Page({
  data: {
    url: ''
  },

  onLoad(options) {
    if (options.url) {
      this.setData({
        url: decodeURIComponent(options.url)
      })
    } else {
      wx.showToast({
        title: '链接地址无效',
        icon: 'none'
      })
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
    }
  }
})
