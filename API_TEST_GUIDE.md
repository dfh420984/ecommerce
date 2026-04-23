# 系统配置功能 API 测试指南

## 前置条件

1. 后端服务已启动（默认端口：8686）
2. 前端管理后台已启动（默认端口：8000）
3. 数据库中已有示例配置数据

## 测试步骤

### 1. 获取管理员Token

首先需要登录获取token，使用以下任一方式：

#### 方式一：通过管理后台界面
1. 访问 http://localhost:8000/login
2. 使用管理员账号登录（默认用户名：admin，密码需要查看数据库或自行设置）
3. 登录后token会自动保存在localStorage中

#### 方式二：通过API直接登录
```bash
curl -X POST http://localhost:8686/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "your_password"
  }'
```

响应示例：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员"
    }
  }
}
```

### 2. 测试管理端接口

将 `<YOUR_TOKEN>` 替换为实际的token值。

#### 获取配置列表
```bash
curl -X GET http://localhost:8686/api/admin/configs \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

预期响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "name": "shop_name",
      "value": "电商小程序",
      "description": "店铺名称",
      "created_at": "2026-04-23T16:00:00+08:00",
      "updated_at": "2026-04-23T16:00:00+08:00"
    },
    ...
  ]
}
```

#### 获取单个配置
```bash
curl -X GET http://localhost:8686/api/admin/configs/1 \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

#### 创建新配置
```bash
curl -X POST http://localhost:8686/api/admin/configs \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test_config",
    "value": "test_value",
    "description": "测试配置"
  }'
```

#### 更新配置
```bash
curl -X PUT http://localhost:8686/api/admin/configs/1 \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "shop_name",
    "value": "新店铺名称",
    "description": "店铺名称（已更新）"
  }'
```

#### 删除配置
```bash
curl -X DELETE http://localhost:8686/api/admin/configs/9 \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

### 3. 测试小程序端接口（无需token）

#### 根据名称获取配置
```bash
curl -X GET http://localhost:8686/api/miniapp/config/shop_name
```

预期响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "name": "shop_name",
    "value": "电商小程序",
    "description": "店铺名称",
    "created_at": "2026-04-23T16:00:00+08:00",
    "updated_at": "2026-04-23T16:00:00+08:00"
  }
}
```

#### 批量获取多个配置
```bash
curl -X POST http://localhost:8686/api/miniapp/configs/batch \
  -H "Content-Type: application/json" \
  -d '{
    "names": ["shop_name", "shop_phone", "shop_address"]
  }'
```

预期响应：
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

### 4. 通过管理后台界面测试

1. 访问 http://localhost:8000
2. 登录管理后台
3. 点击左侧菜单"系统配置"
4. 查看配置列表
5. 点击"添加配置"按钮，填写表单并保存
6. 点击某行的"编辑"按钮，修改配置
7. 点击某行的"删除"按钮，删除配置

### 5. 在小程序中测试

在小程序任意页面的 `.js` 文件中添加：

```javascript
const api = require('../../services/api')

Page({
  onLoad() {
    this.testConfig()
  },
  
  async testConfig() {
    // 测试获取单个配置
    const res1 = await api.getConfig('shop_name')
    console.log('单个配置:', res1)
    
    // 测试批量获取配置
    const res2 = await api.getConfigsByNames([
      'shop_name', 
      'shop_phone', 
      'shop_address'
    ])
    console.log('批量配置:', res2)
  }
})
```

## 常见错误及解决方案

### 1. 401 Unauthorized
**原因**: Token无效或过期  
**解决**: 重新登录获取新的token

### 2. 400 Bad Request
**原因**: 请求参数错误  
**解决**: 检查请求参数格式是否正确

### 3. 配置名称已存在
**原因**: 创建的配置名称与现有配置重复  
**解决**: 使用唯一的配置名称

### 4. 配置不存在
**原因**: 查询的配置ID或名称不存在  
**解决**: 先获取配置列表确认配置是否存在

## 验证清单

- [ ] 可以成功获取配置列表
- [ ] 可以成功获取单个配置
- [ ] 可以成功创建新配置
- [ ] 可以成功更新配置
- [ ] 可以成功删除配置
- [ ] 配置名称唯一性验证生效
- [ ] 配置名称格式验证生效
- [ ] 小程序端可以获取配置
- [ ] 小程序端可以批量获取配置
- [ ] 管理后台界面操作正常

## 性能测试建议

1. **并发测试**: 模拟多个小程序用户同时获取配置
2. **缓存测试**: 验证Redis缓存是否生效（如果实现了）
3. **大数据量测试**: 添加大量配置项，测试列表加载性能

## 下一步

配置功能已经完整实现并可以正常使用。建议：
1. 在实际环境中测试所有功能
2. 根据业务需求添加更多配置项
3. 考虑实现配置缓存以提高性能
4. 添加配置变更日志功能
