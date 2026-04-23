# "关于我们"页面配置指南

## 📋 常见配置项

"关于我们"页面通常包含以下信息，这些都可以从配置接口动态获取：

### 1. **基础信息**
- ✅ **店铺名称** (`shop_name`) - 显示在页面顶部
- ✅ **应用Logo** (`shop_logo`) - 店铺/应用图标
- ✅ **版本号** (`app_version`) - 当前应用版本

### 2. **公司介绍**
- ✅ **公司简介** (`about_intro`) - 公司/产品介绍
- ✅ **使命愿景** (`about_mission`) - 公司使命和愿景（可选）
- ✅ **发展历程** (`about_history`) - 公司发展历程（可选）

### 3. **技术信息**
- ✅ **技术栈** (`about_tech_stack`) - 使用的技术框架
- ✅ **开源协议** (`about_license`) - 开源许可证信息（可选）

### 4. **联系信息**
- ✅ **联系方式** (`about_contact`) - 邮箱、电话、地址等
- ✅ **客服电话** (`shop_phone`) - 客服热线
- ✅ **公司地址** (`shop_address`) - 办公地址
- ✅ **社交媒体** (`social_media`) - 微信、微博等（可选）

### 5. **服务信息**
- ✅ **营业时间** (`business_hours`) - 客服工作时间
- ✅ **退换货政策** (`return_policy`) - 售后政策
- ✅ **隐私政策** (`privacy_policy`) - 隐私条款链接（可选）
- ✅ **用户协议** (`user_agreement`) - 用户协议链接（可选）

### 6. **其他信息**
- ✅ **备案号** (`icp_number`) - ICP备案号（可选）
- ✅ ** copyright** (`copyright_info`) - 版权信息（可选）

---

## 🎯 已实现的配置

当前系统已实现以下配置项：

| 配置名称 | 说明 | 示例值 |
|---------|------|--------|
| `shop_name` | 店铺名称 | 电商小程序 |
| `shop_logo` | 店铺Logo | /uploads/logo.png |
| `app_version` | 应用版本 | 1.0.0 |
| `about_intro` | 公司简介 | 这是一款基于 Go + Vue3 + 微信小程序的电商系统... |
| `about_tech_stack` | 技术栈 | 后端：Go + Gin + GORM...\n前端：Vue3 + Element Plus... |
| `about_contact` | 联系方式 | 邮箱：contact@example.com\n电话：400-123-4567... |
| `business_hours` | 营业时间 | 09:00-21:00 |
| `return_policy` | 退换货政策 | 7天无理由退换货 |

---

## 💡 使用示例

### 在管理后台添加配置

1. 登录管理后台 (http://localhost:8000)
2. 进入"系统配置"页面
3. 点击"添加配置"
4. 填写配置信息：

**示例1：添加公司简介**
```
配置名称: about_intro
配置值: XXX公司成立于2020年，专注于为用户提供优质的购物体验...
描述: 公司简介
```

**示例2：添加联系方式**
```
配置名称: about_contact
配置值: 客服电话：400-123-4567
       邮箱：service@example.com
       地址：北京市朝阳区XXX大厦
描述: 联系方式
```

**示例3：添加隐私政策链接**
```
配置名称: privacy_policy_url
配置值: https://example.com/privacy
描述: 隐私政策链接
```

### 在小程序中使用

```javascript
const config = require('../../utils/config')

Page({
  async onLoad() {
    // 获取单个配置
    const intro = await config.getConfig('about_intro')
    
    // 或批量获取
    const configs = await config.getConfigs([
      'about_intro',
      'about_contact',
      'privacy_policy_url'
    ])
    
    this.setData({
      intro: configs.about_intro,
      contact: configs.about_contact,
      privacyUrl: configs.privacy_policy_url
    })
  }
})
```

---

## 🔧 扩展建议

### 可以添加的配置项

#### 1. 多语言支持
```
about_intro_zh - 中文简介
about_intro_en - 英文简介
```

#### 2. 团队成员
```
team_members - JSON格式的团队信息
[
  {"name": "张三", "role": "CEO", "avatar": "..."},
  {"name": "李四", "role": "CTO", "avatar": "..."}
]
```

#### 3. 荣誉资质
```
certificates - JSON格式的证书列表
[
  {"name": "高新技术企业", "image": "..."},
  {"name": "ISO认证", "image": "..."}
]
```

#### 4. 合作伙伴
```
partners - JSON格式的合作伙伴
[
  {"name": "阿里云", "logo": "..."},
  {"name": "腾讯云", "logo": "..."}
]
```

#### 5. 下载地址
```
ios_download_url - iOS下载地址
android_download_url - Android下载地址
```

---

## 📱 最佳实践

### 1. 配置命名规范
```
✅ 推荐：
- about_intro (关于我们简介)
- about_contact (联系方式)
- privacy_policy_url (隐私政策链接)

❌ 避免：
- intro1, intro2 (无意义)
- myAboutInfo (不要使用my前缀)
```

### 2. 配置值格式

**简单文本：**
```
配置值: 这是一段简单的文本介绍
```

**多行文本（使用 \n 换行）：**
```
配置值: 第一行内容\n第二行内容\n第三行内容
```

**JSON格式（复杂数据）：**
```
配置值: [{"name":"张三","role":"CEO"},{"name":"李四","role":"CTO"}]
```

**URL链接：**
```
配置值: https://example.com/privacy
```

### 3. 性能优化

```javascript
// ✅ 推荐：批量获取，减少请求次数
const configs = await config.getConfigs([
  'about_intro',
  'about_contact',
  'business_hours'
])

// ❌ 不推荐：多次单独请求
const intro = await config.getConfig('about_intro')
const contact = await config.getConfig('about_contact')
const hours = await config.getConfig('business_hours')
```

### 4. 缓存策略

配置工具类已内置1小时缓存，无需额外处理：

```javascript
// 首次调用：从服务器获取
const intro = await config.getConfig('about_intro')

// 1小时内再次调用：从缓存读取
const intro2 = await config.getConfig('about_intro') // 不发起网络请求
```

---

## 🎨 UI设计建议

### 1. 信息分组
将相关信息分组展示，例如：
- 公司信息组（简介、使命、愿景）
- 联系信息组（电话、邮箱、地址）
- 服务信息组（营业时间、退换货政策）

### 2. 图标装饰
为每个信息项添加图标，提升视觉效果：
```html
<view class="info-item">
  <icon type="phone" />
  <text class="label">客服电话</text>
  <text class="value">{{phone}}</text>
</view>
```

### 3. 可点击链接
对于电话、邮箱、网址等，添加点击功能：
```javascript
// 拨打电话
makePhoneCall() {
  wx.makePhoneCall({
    phoneNumber: this.data.phone
  })
}

// 发送邮件
sendEmail() {
  wx.openSystemEmail({
    toEmail: this.data.email
  })
}

// 打开网页
openPrivacyPolicy() {
  wx.navigateTo({
    url: `/pages/webview/webview?url=${encodeURIComponent(this.data.privacyUrl)}`
  })
}
```

---

## 🔍 常见问题

### Q1: 配置更新后多久生效？
**A:** 立即生效。小程序端有1小时缓存，如需立即看到更新，可以：
- 清除缓存：`config.clearConfigCache()`
- 或者重启小程序

### Q2: 配置值太长怎么办？
**A:** 
- 可以使用多行文本（用 `\n` 分隔）
- 或者存储为 JSON 格式
- 也可以在 WXML 中使用 `wx:if` 条件渲染

### Q3: 如何添加富文本内容？
**A:** 
- 配置值中可以使用 HTML 标签
- 使用 `<rich-text>` 组件渲染
```html
<rich-text nodes="{{intro}}"></rich-text>
```

### Q4: 配置不存在时如何处理？
**A:** 提供默认值：
```javascript
this.setData({
  intro: configs.about_intro || '默认简介内容'
})
```

---

## 📊 配置管理建议

### 1. 定期审查
- 每季度检查配置是否过时
- 更新版本号、联系方式等信息
- 删除不再使用的配置

### 2. 备份配置
- 导出重要配置作为备份
- 记录配置变更历史
- 测试环境验证后再更新生产环境

### 3. 权限控制
- 只有管理员可以修改配置
- 敏感信息（如密钥）不要存储在配置表中
- 定期检查配置访问日志

---

## 🎉 总结

通过配置接口管理"关于我们"页面的优势：

✅ **灵活性强** - 随时修改，无需重新发布小程序  
✅ **统一管理** - 所有配置在后台集中管理  
✅ **实时生效** - 修改后立即生效（缓存过期后）  
✅ **易于维护** - 非技术人员也可以更新内容  
✅ **支持多端** - 同一套配置可用于小程序、H5等  

建议根据实际业务需求，合理规划和组织配置项，保持配置的清晰和可维护性。
