const config = require('../../utils/config')
const image = require('../../utils/image')

Page({
  data: {
    shopName: '电商小程序',
    logo: '/static/user/logo.png',
    version: '1.0.0',
    intro: '',
    techStack: '',
    contact: '',
    businessHours: '',
    returnPolicy: ''
  },

  async onLoad() {
    // 设置导航栏标题
    await config.setNavigationBarTitle('about_page_title', '关于我们')
    
    // 加载配置信息
    this.loadAboutInfo()
  },

  async loadAboutInfo() {
    try {
      // 批量获取配置
      const configs = await config.getConfigs([
        'shop_name',
        'shop_logo',
        'app_version',
        'about_intro',
        'about_tech_stack',
        'about_contact',
        'business_hours',
        'return_policy'
      ])

      this.setData({
        shopName: configs.shop_name || '电商小程序',
        logo: image.formatImageUrl(configs.shop_logo) || '/static/user/logo.png',
        version: configs.app_version || '1.0.0',
        intro: configs.about_intro || '',
        techStack: configs.about_tech_stack || '',
        contact: configs.about_contact || '',
        businessHours: configs.business_hours || '09:00-21:00',
        returnPolicy: configs.return_policy || '7天无理由退换货'
      })
    } catch (err) {
      console.error('加载关于信息失败', err)
    }
  }
})
