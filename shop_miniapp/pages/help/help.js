const config = require('../../utils/config')

Page({
  data: {
    searchKeyword: '',
    currentCategory: 0,
    expandedQuestionId: null,
    serviceTime: '',
    categories: [
      { id: 0, name: '全部' },
      { id: 1, name: '购物指南' },
      { id: 2, name: '配送说明' },
      { id: 3, name: '售后服务' },
      { id: 4, name: '账户安全' }
    ],
    questions: [
      // 购物指南
      {
        id: 1,
        categoryId: 1,
        title: '如何下单购买？',
        answer: '1. 浏览商品，点击商品进入详情页\n2. 选择规格和数量，点击"加入购物车"\n3. 进入购物车，勾选商品\n4. 点击"结算"，选择收货地址\n5. 提交订单并完成支付'
      },
      {
        id: 2,
        categoryId: 1,
        title: '如何修改收货地址？',
        answer: '在订单提交前，可以在确认订单页面点击"收货地址"进行修改。\n\n如果订单已提交但未发货，请联系客服修改。\n\n如果订单已发货，则无法修改地址。'
      },
      {
        id: 3,
        categoryId: 1,
        title: '如何使用优惠券？',
        answer: '1. 在"我的-优惠券"中查看可用优惠券\n2. 下单时，系统会自动匹配可用优惠券\n3. 也可以在结算页面手动选择优惠券\n4. 注意查看优惠券的使用条件和有效期'
      },
      {
        id: 4,
        categoryId: 1,
        title: '支持哪些支付方式？',
        answer: '目前支持以下支付方式：\n• 微信支付\n• 支付宝支付\n• 银行卡支付\n\n部分商品可能不支持某些支付方式，请以实际为准。'
      },
      
      // 配送说明
      {
        id: 5,
        categoryId: 2,
        title: '配送范围有哪些？',
        answer: '我们目前支持全国大部分地区配送。\n\n偏远地区（如西藏、新疆等）可能需要额外配送时间，具体以结算页面显示为准。'
      },
      {
        id: 6,
        categoryId: 2,
        title: '多久能收到货？',
        answer: '• 同城配送：1-2天\n• 省内配送：2-3天\n• 跨省配送：3-5天\n• 偏远地区：5-7天\n\n具体时间以物流公司实际配送为准。'
      },
      {
        id: 7,
        categoryId: 2,
        title: '运费是多少？',
        answer: '• 订单满{{freeShippingAmount}}元包邮\n• 不满包邮金额，收取固定运费10元\n• 特殊商品（如生鲜、大件）可能有额外运费\n\n具体运费以结算页面显示为准。'
      },
      {
        id: 8,
        categoryId: 2,
        title: '如何查询物流信息？',
        answer: '1. 进入"我的-我的订单"\n2. 找到对应订单，点击"查看物流"\n3. 即可看到详细的物流跟踪信息\n\n也可以通过物流公司官网查询。'
      },
      
      // 售后服务
      {
        id: 9,
        categoryId: 3,
        title: '退换货政策是什么？',
        answer: '我们提供{{returnPolicy}}服务。\n\n退换货条件：\n• 商品未使用、未拆封\n• 包装完整、配件齐全\n• 在退换货有效期内\n\n以下情况不支持退换货：\n• 定制类商品\n• 鲜活易腐类商品\n• 数字化商品'
      },
      {
        id: 10,
        categoryId: 3,
        title: '如何申请退款？',
        answer: '1. 进入"我的-我的订单"\n2. 找到需要退款的订单\n3. 点击"申请退款"\n4. 选择退款原因，提交申请\n5. 等待商家审核\n6. 审核通过后，退款将原路返回\n\n退款到账时间：1-7个工作日'
      },
      {
        id: 11,
        categoryId: 3,
        title: '退款多久到账？',
        answer: '退款审核通过后：\n• 微信支付：1-3个工作日\n• 支付宝：1-3个工作日\n• 银行卡：3-7个工作日\n\n具体到账时间以银行处理为准。'
      },
      {
        id: 12,
        categoryId: 3,
        title: '商品有质量问题怎么办？',
        answer: '如果收到商品有质量问题：\n1. 请在签收后24小时内联系客服\n2. 提供商品照片和视频\n3. 客服会为您安排退换货\n4. 来回运费由我们承担\n\n请保留好商品包装和配件。'
      },
      
      // 账户安全
      {
        id: 13,
        categoryId: 4,
        title: '如何修改密码？',
        answer: '1. 进入"我的-设置"\n2. 点击"修改密码"\n3. 输入原密码和新密码\n4. 确认修改\n\n如果忘记密码，可以点击"忘记密码"通过手机号重置。'
      },
      {
        id: 14,
        categoryId: 4,
        title: '如何绑定手机号？',
        answer: '1. 进入"我的-设置"\n2. 点击"绑定手机"\n3. 输入手机号和验证码\n4. 完成绑定\n\n一个手机号只能绑定一个账号。'
      },
      {
        id: 15,
        categoryId: 4,
        title: '如何注销账号？',
        answer: '如需注销账号：\n1. 进入"我的-设置"\n2. 点击"注销账号"\n3. 阅读注销须知\n4. 确认注销\n\n注意：注销后数据无法恢复，请谨慎操作。\n注销前请确保：\n• 无进行中的订单\n• 无未处理的售后\n• 账户余额已清零'
      }
    ],
    filteredQuestions: []
  },

  onLoad() {
    this.loadConfig()
    this.filterQuestions()
  },

  async loadConfig() {
    try {
      const configs = await config.getConfigs([
        'free_shipping_amount',
        'return_policy',
        'business_hours'
      ])
      
      this.setData({
        freeShippingAmount: configs.free_shipping_amount || '99',
        returnPolicy: configs.return_policy || '7天无理由退换货',
        serviceTime: configs.business_hours || '09:00-21:00'
      })
      
      // 更新问题中的配置变量
      this.updateQuestionsWithConfig()
    } catch (err) {
      console.error('加载配置失败', err)
    }
  },

  updateQuestionsWithConfig() {
    const { freeShippingAmount, returnPolicy } = this.data
    
    // 更新问题中的占位符
    const updatedQuestions = this.data.questions.map(q => {
      let answer = q.answer
      answer = answer.replace(/\{\{freeShippingAmount\}\}/g, freeShippingAmount)
      answer = answer.replace(/\{\{returnPolicy\}\}/g, returnPolicy)
      return { ...q, answer }
    })
    
    this.setData({
      questions: updatedQuestions
    })
    
    // 重新过滤
    this.filterQuestions()
  },

  onSearchInput(e) {
    this.setData({
      searchKeyword: e.detail.value
    })
    this.filterQuestions()
  },

  onCategoryTap(e) {
    const index = e.currentTarget.dataset.index
    this.setData({
      currentCategory: index,
      expandedQuestionId: null
    })
    this.filterQuestions()
  },

  filterQuestions() {
    const { searchKeyword, currentCategory, questions } = this.data
    
    let filtered = questions
    
    // 按分类过滤
    if (currentCategory > 0) {
      const categoryId = this.data.categories[currentCategory].id
      filtered = filtered.filter(q => q.categoryId === categoryId)
    }
    
    // 按关键词搜索
    if (searchKeyword) {
      const keyword = searchKeyword.toLowerCase()
      filtered = filtered.filter(q => 
        q.title.toLowerCase().includes(keyword) || 
        q.answer.toLowerCase().includes(keyword)
      )
    }
    
    this.setData({
      filteredQuestions: filtered
    })
  },

  onQuestionTap(e) {
    const id = e.currentTarget.dataset.id
    const { expandedQuestionId } = this.data
    
    // 切换展开/收起
    this.setData({
      expandedQuestionId: expandedQuestionId === id ? null : id
    })
  },

  onOnlineService() {
    // TODO: 跳转到在线客服页面或打开客服会话
    wx.showToast({
      title: '在线客服功能开发中',
      icon: 'none'
    })
  },

  onPhoneService() {
    config.getConfig('shop_phone').then(phone => {
      if (phone) {
        wx.makePhoneCall({
          phoneNumber: phone
        })
      } else {
        wx.makePhoneCall({
          phoneNumber: '400-123-4567'
        })
      }
    }).catch(() => {
      wx.makePhoneCall({
        phoneNumber: '400-123-4567'
      })
    })
  }
})
