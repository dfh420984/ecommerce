# 系统配置功能 - 快速开始

## 🚀 5分钟快速上手

### 第一步：启动服务

确保以下服务正在运行：

1. **后端服务** (端口 8686)
   ```bash
   cd shop_api
   go run main.go
   ```

2. **前端管理后台** (端口 8000)
   ```bash
   cd shop_admin
   npm run dev
   ```

3. **数据库** 
   - 确保MySQL已启动
   - 执行 `shop_api/database/shop.sql` 初始化数据

### 第二步：访问管理后台

1. 打开浏览器访问：http://localhost:8000
2. 使用管理员账号登录
3. 在左侧菜单找到并点击 **"系统配置"**

### 第三步：查看示例配置

系统已预置8条常用配置：

| 配置名称 | 值 | 说明 |
|---------|---|------|
| shop_name | 电商小程序 | 店铺名称 |
| shop_logo | /uploads/logo.png | 店铺Logo |
| shop_phone | 400-123-4567 | 客服电话 |
| shop_address | 北京市朝阳区xxx路xxx号 | 店铺地址 |
| min_order_amount | 99 | 最小订单金额 |
| free_shipping_amount | 199 | 包邮金额 |
| return_policy | 7天无理由退换货 | 退换货政策 |
| business_hours | 09:00-21:00 | 营业时间 |

### 第四步：添加新配置

1. 点击 **"添加配置"** 按钮
2. 填写配置信息：
   - **配置名称**: `welcome_message`（只能包含字母、数字和下划线）
   - **配置值**: `欢迎光临本店！`
   - **描述**: `首页欢迎语`
3. 点击 **"确定"** 保存

### 第五步：在小程序中使用

在小程序页面的 `.js` 文件中：

```javascript
const api = require('../../services/api')

Page({
  data: {
    welcomeMessage: '',
    shopName: ''
  },
  
  onLoad() {
    this.loadConfig()
  },
  
  async loadConfig() {
    // 批量获取配置
    const res = await api.getConfigsByNames([
      'welcome_message',
      'shop_name'
    ])
    
    if (res.code === 0) {
      this.setData({
        welcomeMessage: res.data.welcome_message,
        shopName: res.data.shop_name
      })
    }
  }
})
```

在 `.wxml` 文件中使用：

```html
<view class="header">
  <text>{{welcomeMessage}}</text>
  <text>欢迎来到 {{shopName}}</text>
</view>
```

## 💡 实用场景

### 场景1：动态修改店铺公告

1. 在管理后台添加配置：
   - 名称：`shop_announcement`
   - 值：`春节期间正常发货，祝您新年快乐！`
   
2. 小程序首页显示公告：
```javascript
const res = await api.getConfig('shop_announcement')
this.setData({ announcement: res.data.value })
```

### 场景2：灵活调整包邮策略

1. 在管理后台修改配置：
   - 名称：`free_shipping_amount`
   - 值：从 `199` 改为 `159`（促销活动）

2. 购物车页面判断：
```javascript
const config = await api.getConfig('free_shipping_amount')
const freeShippingAmount = parseInt(config.data.value)

if (cartAmount >= freeShippingAmount) {
  this.setData({ isFreeShipping: true })
}
```

### 场景3：临时关闭某项功能

1. 添加功能开关配置：
   - 名称：`enable_coupon`
   - 值：`false`

2. 控制优惠券入口显示：
```javascript
const config = await api.getConfig('enable_coupon')
if (config.data.value === 'true') {
  this.setData({ showCouponEntry: true })
}
```

## 🎯 最佳实践

### 1. 配置命名规范
```
✅ 推荐：
- shop_name
- min_order_amount
- enable_feature_x

❌ 避免：
- ShopName (不要用驼峰)
- minOrderAmount (不要用驼峰)
- config1 (无意义命名)
```

### 2. 配置缓存策略

```javascript
// 小程序端缓存配置
const CACHE_KEY = 'app_config'
const CACHE_EXPIRE = 3600000 // 1小时

async function getConfigWithCache(names) {
  // 尝试从缓存读取
  const cached = wx.getStorageSync(CACHE_KEY)
  if (cached && Date.now() - cached.timestamp < CACHE_EXPIRE) {
    return cached.data
  }
  
  // 缓存失效，从服务器获取
  const res = await api.getConfigsByNames(names)
  
  // 更新缓存
  if (res.code === 0) {
    wx.setStorageSync(CACHE_KEY, {
      data: res.data,
      timestamp: Date.now()
    })
  }
  
  return res.data
}
```

### 3. 配置分类管理

建议按模块前缀命名：
```
shop_*     - 店铺相关配置
order_*    - 订单相关配置
pay_*      - 支付相关配置
user_*     - 用户相关配置
feature_*  - 功能开关
```

## 🔍 故障排查

### 问题1：无法访问管理后台
**检查**：
- 前端服务是否启动（端口8000）
- 后端服务是否启动（端口8686）
- 浏览器控制台是否有错误

### 问题2：配置保存失败
**检查**：
- 配置名称是否已存在
- 配置名称格式是否正确（只能包含字母、数字、下划线）
- 网络连接是否正常

### 问题3：小程序获取配置失败
**检查**：
- 后端服务是否正常运行
- 小程序request域名配置是否正确
- 配置名称是否存在

## 📚 相关文档

- [详细使用说明](./CONFIG_GUIDE.md)
- [API测试指南](./API_TEST_GUIDE.md)
- [实现总结](./IMPLEMENTATION_SUMMARY.md)

## ✨ 功能亮点

- ✅ 实时生效，无需重启服务
- ✅ 支持任意类型的配置值
- ✅ 批量获取，减少请求次数
- ✅ 权限控制，安全可靠
- ✅ 简单易用，灵活扩展

## 🎉 开始使用

现在你已经了解了系统配置功能的基本用法，可以：

1. 登录管理后台添加更多配置
2. 在小程序中集成配置功能
3. 根据业务需求灵活运用
4. 参考文档深入了解高级用法

祝使用愉快！🚀
