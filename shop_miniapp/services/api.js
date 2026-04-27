const { request, uploadImage } = require('../utils/request')

// 小程序 API（方法需与后端路由保持一致）
module.exports = {
  // 图片上传
  uploadImage,

  // 用户相关
  register: (data) => request({ url: '/miniapp/register', method: 'POST', data }),
  login: (data) => request({ url: '/miniapp/login', method: 'POST', data }),
  wechatLogin: (data) => request({ url: '/miniapp/wechat_login', method: 'POST', data }),
  getUserInfo: () => request({ url: '/user/info', method: 'GET' }),
  updateUserInfo: (data) => request({ url: '/user/info', method: 'PUT', data }),

  // 商品相关
  getBanners: () => request({ url: '/miniapp/banners', method: 'GET' }),
  getCategories: () => request({ url: '/miniapp/categories', method: 'GET' }),
  getSubCategories: (id) => request({ url: `/miniapp/categories/${id}/sub`, method: 'GET' }),
  getProducts: (params) => request({ url: '/miniapp/products', method: 'GET', data: params }),
  getRecommendProducts: () => request({ url: '/miniapp/products/recommend', method: 'GET' }),
  getNewProducts: () => request({ url: '/miniapp/products/new', method: 'GET' }),
  getProduct: (id) => request({ url: `/miniapp/products/${id}`, method: 'GET' }),

  // 购物车相关
  getCart: () => request({ url: '/user/cart', method: 'GET' }),
  getCartCount: () => request({ url: '/user/cart/count', method: 'GET' }),
  addCart: (data) => request({ url: '/user/cart', method: 'POST', data }),
  updateCart: (id, data) => request({ url: `/user/cart/${id}`, method: 'PUT', data }),
  selectCart: (id, data) => request({ url: `/user/cart/${id}/select`, method: 'PUT', data }),
  selectAllCart: (data) => request({ url: '/user/cart/select_all', method: 'PUT', data }),
  deleteCart: (id) => request({ url: `/user/cart/${id}`, method: 'DELETE' }),
  clearCart: () => request({ url: '/user/cart', method: 'DELETE' }),

  // 地址相关
  getAddresses: () => request({ url: '/user/addresses', method: 'GET' }),
  getAddress: (id) => request({ url: `/user/addresses/${id}`, method: 'GET' }),
  createAddress: (data) => request({ url: '/user/addresses', method: 'POST', data }),
  updateAddress: (id, data) => request({ url: `/user/addresses/${id}`, method: 'PUT', data }),
  deleteAddress: (id) => request({ url: `/user/addresses/${id}`, method: 'DELETE' }),
  setDefaultAddress: (id) => request({ url: `/user/addresses/${id}/default`, method: 'PUT' }),

  // 订单相关
  getOrders: (params) => request({ url: '/user/orders', method: 'GET', data: params }),
  getOrder: (id) => request({ url: `/user/orders/${id}`, method: 'GET' }),
  createOrder: (data) => request({ url: '/user/orders', method: 'POST', data }),
  cancelOrder: (id) => request({ url: `/user/orders/${id}/cancel`, method: 'PUT' }),
  confirmReceive: (id) => request({ url: `/user/orders/${id}/confirm`, method: 'PUT' }),
  deleteOrder: (id) => request({ url: `/user/orders/${id}`, method: 'DELETE' }),

  // 支付相关
  getPayURL: (data) => request({ url: '/user/pay', method: 'POST', data }),
  queryPayStatus: (id) => request({ url: `/user/pay/status/${id}`, method: 'GET' }),
  applyRefund: (id) => request({ url: `/user/pay/refund/${id}`, method: 'POST' }),

  // 退款管理
  applyRefundNew: (data) => request({ url: '/user/refunds/apply', method: 'POST', data }),
  getMyRefunds: (params) => request({ url: '/user/refunds', method: 'GET', data: params }),
  getRefundDetail: (id) => request({ url: `/user/refunds/${id}`, method: 'GET' }),

  // 订单物流
  getOrderLogistics: (id) => request({ url: `/user/orders/${id}/logistics`, method: 'GET' }),

  // 系统配置相关
  getConfig: (name) => request({ url: `/miniapp/config/${name}`, method: 'GET' }),
  getConfigsByNames: (names) => request({ url: '/miniapp/configs/batch', method: 'POST', data: { names } }),

  // 优惠券相关
  getAvailableCoupons: () => request({ url: '/user/coupons/available', method: 'GET' }),
  receiveCoupon: (id) => request({ url: `/user/coupons/receive/${id}`, method: 'POST' }),
  getMyCoupons: (status) => request({ url: '/user/coupons/my', method: 'GET', data: status ? { status } : {} }),
  getUsableCoupons: (amount) => request({ url: '/user/coupons/usable', method: 'GET', data: { amount } }),

  // 运费计算
  calculateShippingFee: (data) => request({ url: '/miniapp/shipping/calculate', method: 'POST', data }),

  // 商品互动（点赞、收藏）
  likeProduct: (id) => request({ url: `/user/products/${id}/like`, method: 'POST' }),
  unlikeProduct: (id) => request({ url: `/user/products/${id}/like`, method: 'DELETE' }),
  checkLikeStatus: (id) => request({ url: `/user/products/${id}/like/status`, method: 'GET' }),
  favoriteProduct: (id) => request({ url: `/user/products/${id}/favorite`, method: 'POST' }),
  unfavoriteProduct: (id) => request({ url: `/user/products/${id}/favorite`, method: 'DELETE' }),
  checkFavoriteStatus: (id) => request({ url: `/user/products/${id}/favorite/status`, method: 'GET' }),
  getMyFavorites: (params) => request({ url: '/user/favorites', method: 'GET', data: params }),

  // 评价管理
  createReview: (data) => request({ url: '/user/reviews', method: 'POST', data }),
  getProductReviews: (id, params) => request({ url: `/user/products/${id}/reviews`, method: 'GET', data: params }),
  getReviewStats: (id) => request({ url: `/user/products/${id}/reviews/stats`, method: 'GET' }),
  getMyReviews: (params) => request({ url: '/user/my-reviews', method: 'GET', data: params }),
  getCanReviewOrders: () => request({ url: '/user/can-review-orders', method: 'GET' })
}
