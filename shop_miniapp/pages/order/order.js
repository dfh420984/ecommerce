const api = require('../../services/api.js')
const image = require('../../utils/image')
const app = getApp()

Page({
  data: {
    orderId: null,
    orderNo: null,
    order: null,
    orderItems: [],
    selectedAddress: null,
    cart_ids: [],
    productId: null,
    quantity: 1,
    remark: '',
    totalAmount: '0.00',
    freightAmount: '0.00', // 运费
    isFreeShipping: false, // 是否包邮
    usableCoupons: [], // 可用优惠券列表
    selectedCoupon: null, // 选中的优惠券
    couponDiscount: 0, // 优惠金额
    finalAmount: '0.00', // 最终支付金额
    submitting: false // 是否正在提交订单
  },

  onLoad(options) {
    if (options.order_id) {
      this.setData({ orderId: parseInt(options.order_id) })
    }
    if (options.order_no) {
      this.setData({ orderNo: options.order_no })
    }
    if (this.data.orderId) {
      this.loadOrder()
    }
    if (options.cart_ids) {
      this.setData({ cart_ids: options.cart_ids.split(',').map(id => parseInt(id)) })
      this.loadCartPreview()
    }
    if (options.product_id) {
      this.setData({
        productId: parseInt(options.product_id),
        quantity: parseInt(options.quantity || '1')
      })
      this.loadProductPreview()
    }
    this.loadAddresses()
  },

  onShow() {
    if (app.globalData.token) {
      // 保存之前的地址ID
      const oldAddressId = this.data.selectedAddress ? this.data.selectedAddress.id : null
      
      this.loadAddresses().then(() => {
        // 如果地址发生变化，重新计算运费
        const newAddressId = this.data.selectedAddress ? this.data.selectedAddress.id : null
        if (oldAddressId !== newAddressId && this.data.totalAmount > 0) {
          this.calculateFreight()
        }
      })
    }
  },

  // 页面隐藏时保存当前状态
  onHide() {
    // 可以在这里保存一些临时状态
  },

  async loadOrder() {
    try {
      const res = await api.getOrder(this.data.orderId)
      const order = res.data
      // 格式化订单商品图片
      if (order.items) {
        order.items = order.items.map(item => ({
          ...item,
          product_image: image.formatImageUrl(item.product_image)
        }))
      }
      this.setData({
        order,
        orderItems: order.items || [],
        totalAmount: Number(order.pay_amount || 0).toFixed(2)
      })
    } catch (err) {
      console.error(err)
    }
  },

  async loadCartPreview() {
    try {
      const res = await api.getCart()
      const cartIds = this.data.cart_ids
      const selectedList = (res.data.list || []).filter(item => cartIds.includes(item.id))
      const orderItems = selectedList.map(item => {
        const images = image.formatImageUrls(item.product.images || [])
        return {
          id: item.id,
          product_id: item.product_id,
          product_name: item.product.name,
          product_image: images[0] || '',
          price: item.product.price,
          quantity: item.quantity,
          subtotal: Number(item.product.price) * Number(item.quantity)
        }
      })
      const totalAmount = orderItems.reduce((sum, item) => sum + item.subtotal, 0)

      this.setData({
        orderItems,
        totalAmount: totalAmount.toFixed(2),
        finalAmount: totalAmount.toFixed(2)
      })
      
      // 计算运费
      this.calculateFreight()
      
      // 加载可用优惠券
      this.loadUsableCoupons(totalAmount)
    } catch (err) {
      console.error(err)
    }
  },

  async loadProductPreview() {
    try {
      const res = await api.getProduct(this.data.productId)
      const product = res.data
      const images = image.formatImageUrls(product.images || [])
      const quantity = this.data.quantity
      const subtotal = Number(product.price) * Number(quantity)

      this.setData({
        orderItems: [{
          id: product.id,
          product_id: product.id,
          product_name: product.name,
          product_image: images[0] || '',
          price: product.price,
          quantity,
          subtotal
        }],
        totalAmount: subtotal.toFixed(2),
        finalAmount: subtotal.toFixed(2)
      })
      
      // 计算运费
      this.calculateFreight()
      
      // 加载可用优惠券
      this.loadUsableCoupons(subtotal)
    } catch (err) {
      console.error(err)
    }
  },

  async loadAddresses() {
    try {
      const res = await api.getAddresses()
      const addresses = res.data || []
      const defaultAddr = addresses.find(a => a.is_default === 1) || addresses[0]
      this.setData({
        selectedAddress: defaultAddr || addresses[0]
      })
      
      // 加载地址后重新计算运费
      if (this.data.totalAmount > 0) {
        this.calculateFreight()
      }
      
      return Promise.resolve()
    } catch (err) {
      console.error(err)
      return Promise.reject(err)
    }
  },

  // 加载可用优惠券
  async loadUsableCoupons(amount) {
    try {
      const res = await api.getUsableCoupons(amount)
      const coupons = res.data || []
      
      // 格式化优惠券信息
      const formattedCoupons = coupons.map(coupon => {
        const c = coupon.coupon
        let discountText = ''
        let discountAmount = 0
        
        if (c.type === 1) {
          // 满减券
          discountAmount = c.discount_value
          discountText = `减¥${c.discount_value}`
        } else if (c.type === 2) {
          // 折扣券
          const discount = (1 - c.discount_value) * amount
          discountAmount = c.max_discount > 0 && discount > c.max_discount ? c.max_discount : discount
          discountText = `${(c.discount_value * 10).toFixed(1)}折`
        } else if (c.type === 3) {
          // 无门槛券
          discountAmount = c.discount_value
          discountText = `减¥${c.discount_value}`
        }
        
        return {
          ...coupon,
          discountText,
          discountAmount: parseFloat(discountAmount.toFixed(2)),
          expireTimeText: this.formatExpireTime(coupon.expire_time)
        }
      })
      
      // 按优惠金额从高到低排序
      formattedCoupons.sort((a, b) => b.discountAmount - a.discountAmount)
      
      // 默认选择优惠额度最高的一张
      const selectedCoupon = formattedCoupons.length > 0 ? formattedCoupons[0] : null
      const couponDiscount = selectedCoupon ? selectedCoupon.discountAmount : 0
      const freightAmount = parseFloat(this.data.freightAmount)
      const finalAmount = Math.max(0, amount + freightAmount - couponDiscount)
      
      this.setData({
        usableCoupons: formattedCoupons,
        selectedCoupon,
        couponDiscount: couponDiscount.toFixed(2),
        finalAmount: finalAmount.toFixed(2)
      })
    } catch (err) {
      console.error('加载优惠券失败:', err)
    }
  },
  
  // 格式化过期时间
  formatExpireTime(expireTime) {
    if (!expireTime) return ''
    const date = new Date(expireTime)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  },
  
  // 显示优惠券选择弹窗
  onShowCouponSelector() {
    if (this.data.usableCoupons.length === 0) {
      wx.showToast({ title: '暂无可用优惠券', icon: 'none' })
      return
    }
    
    const items = this.data.usableCoupons.map(c => {
      return `${c.coupon.name} - ${c.discountText}${c.coupon.min_amount > 0 ? `(满${c.coupon.min_amount}可用)` : '(无门槛)'}`
    })
    
    wx.showActionSheet({
      itemList: items,
      success: (res) => {
        const selectedCoupon = this.data.usableCoupons[res.tapIndex]
        const couponDiscount = selectedCoupon.discountAmount
        const totalAmount = parseFloat(this.data.totalAmount)
        const freightAmount = parseFloat(this.data.freightAmount)
        const finalAmount = Math.max(0, totalAmount + freightAmount - couponDiscount)
        
        this.setData({
          selectedCoupon,
          couponDiscount: couponDiscount.toFixed(2),
          finalAmount: finalAmount.toFixed(2)
        })
      }
    })
  },
  
  // 取消选择优惠券
  onCancelCoupon() {
    const totalAmount = parseFloat(this.data.totalAmount)
    const freightAmount = parseFloat(this.data.freightAmount)
    const finalAmount = totalAmount + freightAmount
    this.setData({
      selectedCoupon: null,
      couponDiscount: 0,
      finalAmount: finalAmount.toFixed(2)
    })
  },

  // 计算运费
  async calculateFreight() {
    if (!this.data.selectedAddress) {
      return
    }
    
    try {
      const totalAmount = parseFloat(this.data.totalAmount)
      const quantity = this.getOrderQuantity()
      
      const res = await api.calculateShippingFee({
        province: this.data.selectedAddress.province,
        city: this.data.selectedAddress.city,
        district: this.data.selectedAddress.district || '',
        amount: totalAmount,
        quantity: quantity
      })
      
      const freightData = res.data
      const freightAmount = freightData.freight_amount || 0
      const isFreeShipping = freightData.is_free || false
      
      // 更新最终金额
      const couponDiscount = parseFloat(this.data.couponDiscount)
      const finalAmount = Math.max(0, totalAmount + freightAmount - couponDiscount)
      
      this.setData({
        freightAmount: freightAmount.toFixed(2),
        isFreeShipping,
        finalAmount: finalAmount.toFixed(2)
      })
    } catch (err) {
      console.error('计算运费失败:', err)
    }
  },
  
  // 获取订单商品总数量
  getOrderQuantity() {
    return this.data.orderItems.reduce((sum, item) => sum + item.quantity, 0)
  },

  onAddressTap() {
    const currentId = this.data.selectedAddress ? this.data.selectedAddress.id : ''
    wx.navigateTo({ url: `/pages/address/address?mode=select&id=${currentId}` })
  },

  // 从地址选择页面返回时的回调
  onShowAddressSelect(address) {
    if (address) {
      this.setData({
        selectedAddress: address
      })
      // 重新计算运费
      this.calculateFreight()
    }
  },

  onRemarkInput(e) {
    this.setData({ remark: e.detail.value })
  },

  async onSubmit() {
    if (!this.data.selectedAddress) {
      wx.showToast({ title: '请选择收货地址', icon: 'none' })
      return
    }

    // 防止重复提交
    if (this.data.submitting) {
      wx.showToast({ title: '正在提交订单...', icon: 'none' })
      return
    }

    try {
      // 设置提交状态
      this.setData({ submitting: true })

      let res
      const orderData = {
        address_id: this.data.selectedAddress.id,
        remark: this.data.remark
      }
      
      // 添加优惠券ID（如果选择了）
      if (this.data.selectedCoupon) {
        orderData.coupon_id = this.data.selectedCoupon.id
      }
      
      if (this.data.cart_ids.length > 0) {
        orderData.cart_ids = this.data.cart_ids
        res = await api.createOrder(orderData)
      } else if (this.data.productId) {
        orderData.product_id = this.data.productId
        orderData.quantity = this.data.quantity
        res = await api.createOrder(orderData)
      } else {
        wx.showToast({ title: '订单商品不能为空', icon: 'none' })
        this.setData({ submitting: false })
        return
      }

      // 提交成功后跳转到支付页面
      wx.redirectTo({ 
        url: `/pages/pay/pay?order_id=${res.data.order_id}&order_no=${res.data.order_no}` 
      })
    } catch (err) {
      console.error(err)
      // 提交失败，重置状态
      this.setData({ submitting: false })
    }
  },

  // 页面卸载时重置提交状态
  onUnload() {
    this.setData({ submitting: false })
  }
})
