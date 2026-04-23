# 系统配置功能使用说明

## 功能概述

系统配置功能允许管理员在后台动态管理系统参数，小程序端可以实时获取这些配置信息，实现灵活的运营配置。

## 后端接口

### 1. 管理端接口（需要管理员权限）

#### 获取配置列表
```
GET /api/admin/configs
```

#### 获取单个配置
```
GET /api/admin/configs/:id
```

#### 创建配置
```
POST /api/admin/configs
Content-Type: application/json

{
  "name": "配置名称",
  "value": "配置值",
  "description": "配置描述"
}
```

#### 更新配置
```
PUT /api/admin/configs/:id
Content-Type: application/json

{
  "name": "配置名称",
  "value": "配置值",
  "description": "配置描述"
}
```

#### 删除配置
```
DELETE /api/admin/configs/:id
```

### 2. 小程序端接口（公开接口）

#### 根据名称获取单个配置
```
GET /api/miniapp/config/:name
```

响应示例：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "name": "shop_name",
    "value": "电商小程序",
    "description": "店铺名称"
  }
}
```

#### 批量获取多个配置
```
POST /api/miniapp/configs/batch
Content-Type: application/json

{
  "names": ["shop_name", "shop_phone", "shop_address"]
}
```

响应示例：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "shop_name": "电商小程序",
    "shop_phone": "400-123-4567",
    "shop_address": "北京市朝阳区xxx路xxx号"
  }
}
```

## 前端使用示例

### 管理后台使用

1. 访问系统配置页面：`/configs`
2. 点击"添加配置"按钮
3. 填写配置信息：
   - **配置名称**：英文标识符，如 `shop_name`、`min_order_amount`
   - **配置值**：具体的配置值
   - **描述**：配置的说明信息
4. 保存后即可在小程序端使用

### 小程序端使用

#### 方式一：获取单个配置

```javascript
const api = require('../../services/api')

// 在页面的 onLoad 或 onShow 中获取配置
Page({
  data: {
    shopName: ''
  },
  
  onLoad() {
    this.loadConfig()
  },
  
  async loadConfig() {
    try {
      const res = await api.getConfig('shop_name')
      if (res.code === 0) {
        this.setData({
          shopName: res.data.value
        })
      }
    } catch (error) {
      console.error('获取配置失败', error)
    }
  }
})
```

#### 方式二：批量获取多个配置（推荐）

```javascript
const api = require('../../services/api')

Page({
  data: {
    config: {}
  },
  
  onLoad() {
    this.loadConfigs()
  },
  
  async loadConfigs() {
    try {
      // 一次性获取多个配置
      const res = await api.getConfigsByNames([
        'shop_name',
        'shop_logo',
        'shop_phone',
        'shop_address',
        'min_order_amount',
        'free_shipping_amount'
      ])
      
      if (res.code === 0) {
        this.setData({
          config: res.data
        })
        
        // 使用配置
        console.log('店铺名称:', res.data.shop_name)
        console.log('客服电话:', res.data.shop_phone)
        console.log('最小订单金额:', res.data.min_order_amount)
      }
    } catch (error) {
      console.error('获取配置失败', error)
    }
  }
})
```

#### 在 WXML 中使用

```html
<view class="shop-info">
  <text>{{config.shop_name}}</text>
  <text>客服电话：{{config.shop_phone}}</text>
  <text>营业时间：{{config.business_hours}}</text>
</view>

<!-- 判断是否满足包邮条件 -->
<view wx:if="{{cartAmount >= config.free_shipping_amount}}">
  已满足包邮条件
</view>
```

## 常用配置项示例

| 配置名称 | 配置值示例 | 描述 |
|---------|-----------|------|
| shop_name | 电商小程序 | 店铺名称 |
| shop_logo | /uploads/logo.png | 店铺Logo |
| shop_phone | 400-123-4567 | 客服电话 |
| shop_address | 北京市朝阳区xxx路xxx号 | 店铺地址 |
| min_order_amount | 99 | 最小订单金额 |
| free_shipping_amount | 199 | 包邮金额 |
| return_policy | 7天无理由退换货 | 退换货政策 |
| business_hours | 09:00-21:00 | 营业时间 |

## 注意事项

1. **配置名称规范**：
   - 只能包含字母、数字和下划线
   - 建议使用小写字母
   - 命名要有意义，如 `shop_name`、`min_order_amount`

2. **配置值类型**：
   - 所有配置值都以字符串形式存储
   - 如需数值类型，使用时需要转换：`parseInt(config.min_order_amount)`

3. **性能优化**：
   - 建议在小程序启动时批量获取所需配置
   - 可以将配置缓存到本地存储，减少请求次数
   - 配置变更频率较低，可以适当延长缓存时间

4. **安全性**：
   - 敏感配置（如密钥）不应存储在 configs 表中
   - 管理端接口有权限验证，只有管理员可以修改配置
   - 小程序端接口为公开接口，不要存储敏感信息

## 应用场景

1. **店铺信息管理**：店铺名称、Logo、联系方式等
2. **运营策略配置**：包邮金额、最小订单金额、优惠规则等
3. **页面内容配置**：公告信息、活动说明、服务协议等
4. **功能开关控制**：是否开启某项功能、显示隐藏模块等
