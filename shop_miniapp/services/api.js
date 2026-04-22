const { get, post } = require('../utils/request')

// 小程序 API
module.exports = {
  // 用户相关
  register: (data) => post('/miniapp/register', data),
  login: (data) => post('/miniapp/login', data),
  wechatLogin: (data) => post('/miniapp/wechat_login', data),
  getUserInfo: () => get('/user/info'),
  updateUserInfo: (data) => post('/user/info', data),

  // 商品相关
  getBanners: () => get('/miniapp/banners'),
  getCategories: () => get('/miniapp/categories'),
  getSubCategories: (id) => get(`/miniapp/categories/${id}/sub`),
  getProducts: (params) => get('/miniapp/products', params),
  getRecommendProducts: () => get('/miniapp/products/recommend'),
  getNewProducts: () => get('/miniapp/products/new'),
  getProduct: (id) => get(`/miniapp/products/${id}`),

  // 购物车相关
  getCart: () => get('/user/cart'),
  getCartCount: () => get('/user/cart/count'),
  addCart: (data) => post('/user/cart', data),
  updateCart: (id, data) => post(`/user/cart/${id}`, data),
  selectCart: (id, data) => post(`/user/cart/${id}/select`, data),
  selectAllCart: (data) => post('/user/cart/select_all', data),
  deleteCart: (id) => post(`/user/cart/${id}`),
  clearCart: () => post('/user/cart'),

  // 地址相关
  getAddresses: () => get('/user/addresses'),
  getAddress: (id) => get(`/user/addresses/${id}`),
  createAddress: (data) => post('/user/addresses', data),
  updateAddress: (id, data) => post(`/user/addresses/${id}`, data),
  deleteAddress: (id) => post(`/user/addresses/${id}`),
  setDefaultAddress: (id) => post(`/user/addresses/${id}/default`),

  // 订单相关
  getOrders: (params) => get('/user/orders', params),
  getOrder: (id) => get(`/user/orders/${id}`),
  createOrder: (data) => post('/user/orders', data),
  cancelOrder: (id) => post(`/user/orders/${id}/cancel`),
  confirmReceive: (id) => post(`/user/orders/${id}/confirm`),
  deleteOrder: (id) => post(`/user/orders/${id}`),

  // 支付相关
  getPayURL: (data) => post('/user/pay', data),
  queryPayStatus: (id) => get(`/user/pay/status/${id}`),
  applyRefund: (id, data) => post(`/user/pay/refund/${id}`, data)
}
