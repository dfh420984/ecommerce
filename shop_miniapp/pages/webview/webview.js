const config = require('../../utils/config.js')

Page({
  data: {
    url: ''
  },

  async onLoad(options) {
    // 设置导航栏标题
    await config.setNavigationBarTitle('webview_page_title', '网页')
    
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
