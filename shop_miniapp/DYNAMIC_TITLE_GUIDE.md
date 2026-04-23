# 小程序动态导航栏标题使用指南

## 概述

虽然 `app.json` 中的 `navigationBarTitleText` 不能直接从配置接口获取，但我们可以通过 `wx.setNavigationBarTitle()` API 在页面加载时动态设置导航栏标题。

## 实现方案

### 1. 配置工具类

已创建 `utils/config.js` 工具类，提供以下功能：

- ✅ 获取配置值（带缓存）
- ✅ 批量获取配置（带缓存）
- ✅ 设置导航栏标题
- ✅ 清除配置缓存

### 2. 使用方法

#### 方法一：使用工具类（推荐）

```javascript
const config = require('../../utils/config')

Page({
  async onLoad() {
    // 从配置接口获取并设置导航栏标题
    // 第一个参数：配置名称
    // 第二个参数：默认标题（配置获取失败时使用）
    await config.setNavigationBarTitle('shop_name', '电商小程序')
    
    // 其他初始化代码...
  }
})
```

#### 方法二：直接调用 API

```javascript
const api = require('../../services/api')

Page({
  async onLoad() {
    try {
      const res = await api.getConfig('shop_name')
      if (res.code === 0 && res.data.value) {
        wx.setNavigationBarTitle({
          title: res.data.value
        })
      }
    } catch (err) {
      console.error('获取配置失败', err)
    }
  }
})
```

## 实际应用示例

### 示例1：首页动态标题

```javascript
// pages/index/index.js
const config = require('../../utils/config')

Page({
  async onLoad() {
    // 设置首页标题为店铺名称
    await config.setNavigationBarTitle('shop_name', '电商小程序')
  }
})
```

### 示例2：关于我们页面

```javascript
// pages/about/about.js
const config = require('../../utils/config')

Page({
  async onLoad() {
    // 设置标题为"关于XXX"
    const shopName = await config.getConfig('shop_name')
    const title = shopName ? `关于${shopName}` : '关于我们'
    
    wx.setNavigationBarTitle({
      title: title
    })
  }
})
```

### 示例3：商品列表页（根据分类动态标题）

```javascript
// pages/product-list/product-list.js
const config = require('../../utils/config')

Page({
  data: {
    categoryId: null,
    categoryName: ''
  },
  
  async onLoad(options) {
    if (options.categoryId) {
      this.setData({ categoryId: options.categoryId })
      await this.loadCategoryInfo()
    }
  },
  
  async loadCategoryInfo() {
    // 假设有一个获取分类详情的API
    const category = await getCategoryDetail(this.data.categoryId)
    
    if (category) {
      wx.setNavigationBarTitle({
        title: category.name
      })
    }
  }
})
```

### 示例4：订单详情页（根据订单状态动态标题）

```javascript
// pages/order-detail/order-detail.js
Page({
  data: {
    orderId: null
  },
  
  async onLoad(options) {
    if (options.id) {
      this.setData({ orderId: options.id })
      await this.loadOrderDetail()
    }
  },
  
  async loadOrderDetail() {
    const order = await getOrderDetail(this.data.orderId)
    
    if (order) {
      // 根据订单状态设置不同标题
      const statusMap = {
        0: '待付款',
        1: '待发货',
        2: '待收货',
        3: '已完成',
        4: '已取消'
      }
      
      const title = `${statusMap[order.status] || '订单详情'} - ${order.order_no}`
      
      wx.setNavigationBarTitle({
        title: title
      })
    }
  }
})
```

## 高级用法

### 1. 批量获取多个配置

```javascript
const config = require('../../utils/config')

Page({
  async onLoad() {
    // 批量获取多个配置
    const configs = await config.getConfigs([
      'shop_name',
      'shop_phone',
      'business_hours'
    ])
    
    console.log(configs.shop_name)      // 店铺名称
    console.log(configs.shop_phone)     // 客服电话
    console.log(configs.business_hours) // 营业时间
    
    // 设置标题
    wx.setNavigationBarTitle({
      title: configs.shop_name || '电商小程序'
    })
  }
})
```

### 2. 不使用缓存（获取最新配置）

```javascript
// 第二个参数设为 false，强制从服务器获取
const value = await config.getConfig('shop_name', false)
```

### 3. 清除配置缓存

```javascript
// 当管理员更新配置后，可以清除缓存
config.clearConfigCache()
```

## 性能优化建议

### 1. 利用缓存

工具类默认启用缓存（1小时），减少重复请求：

```javascript
// 首次调用：从服务器获取
const name1 = await config.getConfig('shop_name')

// 1小时内再次调用：从缓存读取，不发起网络请求
const name2 = await config.getConfig('shop_name')
```

### 2. 预加载配置

在 `app.js` 中预加载常用配置：

```javascript
// app.js
const config = require('./utils/config')

App({
  async onLaunch() {
    // 预加载常用配置
    await config.getConfigs([
      'shop_name',
      'shop_logo',
      'shop_phone'
    ])
  }
})
```

### 3. 并行加载

```javascript
async onLoad() {
  // 并行执行：设置标题和加载数据
  await Promise.all([
    config.setNavigationBarTitle('shop_name', '电商小程序'),
    this.loadData()
  ])
}
```

## 注意事项

### 1. 异步处理

`setNavigationBarTitle` 是异步操作，如果需要确保标题设置完成后再执行其他操作，使用 `await`：

```javascript
// ✅ 正确
await config.setNavigationBarTitle('shop_name', '电商小程序')
console.log('标题已设置')

// ❌ 可能有问题
config.setNavigationBarTitle('shop_name', '电商小程序')
console.log('标题可能还未设置完成')
```

### 2. 默认标题

始终提供默认标题，防止配置获取失败时标题为空：

```javascript
// ✅ 推荐
await config.setNavigationBarTitle('shop_name', '电商小程序')

// ❌ 不推荐（配置不存在时标题为空）
await config.setNavigationBarTitle('shop_name')
```

### 3. TabBar 页面限制

TabBar 页面不能使用 `wx.navigateTo`，但可以使用 `wx.setNavigationBarTitle`：

```javascript
// pages/cart/cart.js (TabBar 页面)
Page({
  async onLoad() {
    // 可以正常设置标题
    await config.setNavigationBarTitle('cart_title', '购物车')
  }
})
```

### 4. 配置不存在的情况

如果配置不存在，会返回 `null`，此时会使用默认标题：

```javascript
// 如果 'non_existent_config' 不存在，标题将设置为'默认标题'
await config.setNavigationBarTitle('non_existent_config', '默认标题')
```

## 常见问题

### Q1: 为什么不在 app.json 中直接写死标题？

**A:** 使用配置接口的好处：
- ✅ 可以随时修改标题，无需重新发布小程序
- ✅ 可以根据运营活动动态调整
- ✅ 支持多语言、个性化等场景

### Q2: 每次页面加载都会请求配置接口吗？

**A:** 不会。工具类有缓存机制，1小时内相同配置只会请求一次。

### Q3: 如何立即看到配置变更的效果？

**A:** 
1. 清除缓存：`config.clearConfigCache()`
2. 或者等待1小时缓存过期
3. 或者重启小程序

### Q4: 能否在 JSON 配置文件中引用配置接口？

**A:** 不能。JSON 配置文件是静态的，无法执行 JavaScript 代码。必须通过 `wx.setNavigationBarTitle()` API 动态设置。

## 完整示例

查看以下文件的完整实现：
- `utils/config.js` - 配置工具类
- `pages/index/index.js` - 首页使用示例

## 总结

虽然不能在 `app.json` 中直接从配置接口获取标题，但通过以下方式可以实现同样的效果：

1. ✅ 使用 `wx.setNavigationBarTitle()` API 动态设置
2. ✅ 创建配置工具类统一管理
3. ✅ 利用缓存提高性能
4. ✅ 提供默认标题保证用户体验

这种方式更加灵活，可以实现动态化的标题管理！🎉
