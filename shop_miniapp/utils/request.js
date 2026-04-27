const app = getApp()

const request = (options) => {
  return new Promise((resolve, reject) => {
    const token = app.globalData.token
    const header = options.header || {}

    if (token) {
      header['Authorization'] = `Bearer ${token}`
    }

    wx.request({
      url: app.globalData.apiBase + options.url,
      method: options.method || 'GET',
      data: options.data || {},
      header: header,
      success: (res) => {
        if (res.data.code === 0) {
          resolve(res.data)
        } else if (res.data.code === 401) {
          app.clearSession()
          wx.navigateTo({ url: '/pages/login/login' })
          reject(new Error(res.data.msg || '未授权'))
        } else {
          wx.showToast({ title: res.data.msg || '请求失败', icon: 'none' })
          reject(new Error(res.data.msg || '请求失败'))
        }
      },
      fail: (err) => {
        wx.showToast({ title: '网络错误', icon: 'none' })
        reject(err)
      }
    })
  })
}

const login = (data) => request({ url: '/miniapp/login', method: 'POST', data })
const wechatLogin = (data) => request({ url: '/miniapp/wechat_login', method: 'POST', data })
const register = (data) => request({ url: '/miniapp/register', method: 'POST', data })
const getUserInfo = () => request({ url: '/user/info' })
const updateUserInfo = (data) => request({ url: '/user/info', method: 'PUT', data })

const getBanners = () => request({ url: '/miniapp/banners' })
const getCategories = () => request({ url: '/miniapp/categories' })
const getProducts = (params) => request({ url: '/miniapp/products', data: params })
const getRecommendProducts = () => request({ url: '/miniapp/products/recommend' })
const getNewProducts = () => request({ url: '/miniapp/products/new' })
const getProduct = (id) => request({ url: `/miniapp/products/${id}` })

const getAddresses = () => request({ url: '/user/addresses' })
const getAddress = (id) => request({ url: `/user/addresses/${id}` })
const createAddress = (data) => request({ url: '/user/addresses', method: 'POST', data })
const updateAddress = (id, data) => request({ url: `/user/addresses/${id}`, method: 'PUT', data })
const deleteAddress = (id) => request({ url: `/user/addresses/${id}`, method: 'DELETE' })
const setDefaultAddress = (id) => request({ url: `/user/addresses/${id}/default`, method: 'PUT' })

const getCart = () => request({ url: '/user/cart' })
const getCartCount = () => request({ url: '/user/cart/count' })
const addCart = (data) => request({ url: '/user/cart', method: 'POST', data })
const updateCart = (id, data) => request({ url: `/user/cart/${id}`, method: 'PUT', data })
const selectCart = (id, data) => request({ url: `/user/cart/${id}/select`, method: 'PUT', data })
const selectAllCart = (data) => request({ url: '/user/cart/select_all', method: 'PUT', data })
const deleteCart = (id) => request({ url: `/user/cart/${id}`, method: 'DELETE' })
const clearCart = () => request({ url: '/user/cart', method: 'DELETE' })

const getOrders = (params) => request({ url: '/user/orders', data: params })
const getOrder = (id) => request({ url: `/user/orders/${id}` })
const createOrder = (data) => request({ url: '/user/orders', method: 'POST', data })
const cancelOrder = (id) => request({ url: `/user/orders/${id}/cancel`, method: 'PUT' })
const confirmReceive = (id) => request({ url: `/user/orders/${id}/confirm`, method: 'PUT' })
const deleteOrder = (id) => request({ url: `/user/orders/${id}`, method: 'DELETE' })

const getPayURL = (data) => request({ url: '/user/pay', method: 'POST', data })
const getPayStatus = (id) => request({ url: `/user/pay/status/${id}` })
const mockPaySuccess = (id) => request({ url: `/user/pay/mock_success/${id}`, method: 'POST' })
const applyRefund = (id) => request({ url: `/user/pay/refund/${id}`, method: 'POST' })

// 优惠券相关
const getAvailableCoupons = () => request({ url: '/user/coupons/available' })
const receiveCoupon = (id) => request({ url: `/user/coupons/receive/${id}`, method: 'POST' })
const getMyCoupons = (status) => request({ url: '/user/coupons/my', data: status ? { status } : {} })
const getUsableCoupons = (amount) => request({ url: '/user/coupons/usable', data: { amount } })

// 商品互动（点赞、收藏）
const likeProduct = (id) => request({ url: `/user/products/${id}/like`, method: 'POST' })
const unlikeProduct = (id) => request({ url: `/user/products/${id}/like`, method: 'DELETE' })
const checkLikeStatus = (id) => request({ url: `/user/products/${id}/like/status` })
const favoriteProduct = (id) => request({ url: `/user/products/${id}/favorite`, method: 'POST' })
const unfavoriteProduct = (id) => request({ url: `/user/products/${id}/favorite`, method: 'DELETE' })
const checkFavoriteStatus = (id) => request({ url: `/user/products/${id}/favorite/status` })
const getMyFavorites = (params) => request({ url: '/user/favorites', data: params })

// 图片上传
const uploadImage = (filePath) => {
  return new Promise((resolve, reject) => {
    const token = app.globalData.token
    wx.uploadFile({
      url: app.globalData.apiBase + '/upload',
      filePath: filePath,
      name: 'file',
      header: {
        'Authorization': `Bearer ${token}`
      },
      success: (res) => {
        const data = JSON.parse(res.data)
        if (data.code === 0) {
          resolve(data)
        } else {
          wx.showToast({ title: data.msg || '上传失败', icon: 'none' })
          reject(new Error(data.msg || '上传失败'))
        }
      },
      fail: (err) => {
        wx.showToast({ title: '网络错误', icon: 'none' })
        reject(err)
      }
    })
  })
}

// 通用 GET 和 POST 方法
const get = (url, params) => request({ url, method: 'GET', data: params })
const post = (url, data) => request({ url, method: 'POST', data })

module.exports = {
  request,
  get,
  post,
  uploadImage,
  login,
  wechatLogin,
  register,
  getUserInfo,
  updateUserInfo,
  getBanners,
  getCategories,
  getProducts,
  getRecommendProducts,
  getNewProducts,
  getProduct,
  getAddresses,
  getAddress,
  createAddress,
  updateAddress,
  deleteAddress,
  setDefaultAddress,
  getCart,
  getCartCount,
  addCart,
  updateCart,
  selectCart,
  selectAllCart,
  deleteCart,
  clearCart,
  getOrders,
  getOrder,
  createOrder,
  cancelOrder,
  confirmReceive,
  deleteOrder,
  getPayURL,
  getPayStatus,
  mockPaySuccess,
  applyRefund,
  getAvailableCoupons,
  receiveCoupon,
  getMyCoupons,
  getUsableCoupons,
  likeProduct,
  unlikeProduct,
  checkLikeStatus,
  favoriteProduct,
  unfavoriteProduct,
  checkFavoriteStatus,
  getMyFavorites
}
