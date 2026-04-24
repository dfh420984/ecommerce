# 电商系统 V2.0 升级指南

## 📋 升级内容概览

本次升级包含以下核心功能：

### ✅ 已实现功能

1. **订单超时自动取消** 
   - 订单创建时设置30分钟过期时间
   - 定时任务每5分钟检查并取消超时订单
   - 自动释放预占库存和优惠券

2. **库存管理系统**
   - 预占库存机制（stock_reserved字段）
   - 支付成功后正式扣减库存
   - 订单取消/超时自动释放预占库存

3. **完整退款流程**
   - 用户提交退款申请（支持上传凭证）
   - 后台审核退款申请
   - 调用微信/支付宝退款API（支持模拟模式）
   - 退款成功后自动恢复库存（未发货订单）

4. **物流跟踪系统**
   - 对接快递100 API
   - 支持模拟物流数据（开发环境）
   - 查询订单物流轨迹

5. **数据统计Dashboard**
   - 今日/昨日/本月销售额统计
   - 新增用户数统计
   - 订单量统计
   - 热销商品TOP10
   - 销售趋势图表（最近7天）
   - 用户增长趋势（最近7天）

6. **新增数据表**
   - refund_applications - 退款申请表
   - product_reviews - 商品评价表
   - user_favorites - 用户收藏表
   - browse_histories - 浏览历史表

---

## 🚀 部署步骤

### 1. 数据库迁移

执行数据库升级脚本：

```bash
mysql -u root -p shop_db < shop_api/database/migration_v2.sql
```

或者在MySQL客户端中直接执行 `migration_v2.sql` 文件的内容。

**验证迁移是否成功：**

```sql
-- 检查订单表新字段
DESC orders;
-- 应该看到: expire_time, stock_reserved

-- 检查商品表新字段
DESC products;
-- 应该看到: avg_rating, review_count, stock_reserved

-- 检查新表是否创建
SHOW TABLES LIKE 'refund_applications';
SHOW TABLES LIKE 'product_reviews';
SHOW TABLES LIKE 'user_favorites';
SHOW TABLES LIKE 'browse_histories';
```

### 2. 更新配置文件

编辑 `shop_api/config.yaml`，添加物流配置（可选）：

```yaml
logistics:
  provider: kuaidi100  # kuaidi100/aliyun
  api_key: ""  # 快递100 API Key，留空则使用模拟数据
  customer: ""  # 快递100 Customer ID
```

**注意：** 
- 如果不配置 `api_key`，系统将使用模拟物流数据
- 如需使用真实物流查询，需要到快递100官网申请API Key

### 3. 重新编译并启动后端

```bash
cd shop_api

# 下载依赖（如果有新增依赖）
go mod tidy

# 编译
go build -o shop_api main.go

# 启动服务
./shop_api
```

或者直接运行：

```bash
go run main.go
```

**验证启动成功：**

查看日志输出，应该看到：
```
启动定时任务...
定时任务启动成功
Server starting on 0.0.0.0:8686
```

### 4. 测试新功能

#### 测试1：订单超时自动取消

1. 创建一个订单（不支付）
2. 等待30分钟，或手动修改订单的expire_time为过去时间
3. 等待5分钟（定时任务执行周期）
4. 检查订单状态是否变为"已取消"
5. 检查库存是否释放

**快速测试方法：**

```sql
-- 手动设置订单过期时间为过去
UPDATE orders 
SET expire_time = DATE_SUB(NOW(), INTERVAL 1 HOUR)
WHERE order_no = 'YOUR_ORDER_NO';

-- 等待5分钟后检查
SELECT order_status, cancel_time FROM orders WHERE order_no = 'YOUR_ORDER_NO';
```

#### 测试2：退款流程

**用户端：**

```bash
# 申请退款
curl -X POST http://localhost:8686/api/user/refunds/apply \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": 1,
    "refund_type": "refund_only",
    "reason": "商品质量问题",
    "images": ["http://example.com/image1.jpg"],
    "refund_amount": 99.00
  }'

# 查看我的退款申请
curl -X GET http://localhost:8686/api/user/refunds \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**管理端：**

```bash
# 查看退款申请列表
curl -X GET http://localhost:8686/api/admin/refunds \
  -H "Authorization: Bearer ADMIN_TOKEN"

# 审核通过
curl -X POST http://localhost:8686/api/admin/refunds/1/approve \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "handler_reply": "审核通过，正在退款"
  }'

# 拒绝退款
curl -X POST http://localhost:8686/api/admin/refunds/1/reject \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "handler_reply": "不符合退款条件"
  }'
```

#### 测试3：物流查询

```bash
# 查询订单物流
curl -X GET http://localhost:8686/api/user/orders/1/logistics \
  -H "Authorization: Bearer YOUR_TOKEN"
```

返回示例（模拟数据）：
```json
{
  "code": 200,
  "data": {
    "express_company": "SF Express",
    "express_no": "SF1234567890",
    "tracks": [
      {
        "time": "2026-04-24 10:30:00",
        "status": "已签收",
        "desc": "您的快件已被签收"
      },
      ...
    ]
  }
}
```

#### 测试4：数据统计

```bash
# 获取仪表盘统计数据
curl -X GET http://localhost:8686/api/admin/statistics/dashboard \
  -H "Authorization: Bearer ADMIN_TOKEN"

# 获取销售趋势
curl -X GET http://localhost:8686/api/admin/statistics/sales-trend \
  -H "Authorization: Bearer ADMIN_TOKEN"

# 获取用户增长趋势
curl -X GET http://localhost:8686/api/admin/statistics/users-trend \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

---

## 📊 API接口清单

### 用户端接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | /api/user/refunds/apply | 申请退款 | 需登录 |
| GET | /api/user/refunds | 我的退款列表 | 需登录 |
| GET | /api/user/refunds/:id | 退款详情 | 需登录 |
| GET | /api/user/orders/:id/logistics | 查询订单物流 | 需登录 |

### 管理端接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/admin/refunds | 退款申请列表 | 需管理员 |
| GET | /api/admin/refunds/:id | 退款详情 | 需管理员 |
| POST | /api/admin/refunds/:id/approve | 审核通过退款 | 需管理员 |
| POST | /api/admin/refunds/:id/reject | 拒绝退款申请 | 需管理员 |
| GET | /api/admin/statistics/dashboard | 仪表盘统计 | 需管理员 |
| GET | /api/admin/statistics/sales-trend | 销售趋势 | 需管理员 |
| GET | /api/admin/statistics/users-trend | 用户趋势 | 需管理员 |

---

## 🔧 配置说明

### 快递100 API配置（可选）

如果需要真实的物流查询功能：

1. 访问 [快递100官网](https://www.kuaidi100.com/)
2. 注册账号并申请API Key
3. 在 `config.yaml` 中配置：

```yaml
logistics:
  provider: kuaidi100
  api_key: "your-api-key-here"
  customer: "your-customer-id-here"
```

**常用快递公司编码：**
- SF - 顺丰速运
- YTO - 圆通速递
- STO - 申通快递
- ZTO - 中通快递
- YD - 韵达快递
- EMS - EMS

### 微信支付退款配置（待实现）

目前退款功能使用模拟模式。如需接入真实微信支付退款：

1. 在 `services/pay.go` 中的 `WechatRefund` 方法实现微信退款API调用
2. 配置微信商户证书路径
3. 参考微信支付V3 API文档实现退款接口

---

## ⚠️ 注意事项

### 1. 定时任务

- 订单超时检查任务每5分钟执行一次
- 可以在 `tasks/cron_jobs.go` 中调整执行频率
- 生产环境建议使用专业的任务调度工具（如Cron、Systemd Timer）

### 2. 库存管理

- 订单创建时会预占库存（stock_reserved + quantity）
- 支付成功后正式扣减库存（stock - quantity, stock_reserved - quantity）
- 订单取消/超时会自动释放预占库存

### 3. 退款流程

- 仅已支付的订单可以申请退款
- 退款审核通过后会自动调用支付平台退款接口
- 如果商品未发货，退款成功后会恢复库存

### 4. 物流查询

- 未配置API Key时使用模拟数据
- 模拟数据仅用于开发和测试
- 生产环境建议配置真实的物流API

---

## 🐛 常见问题

### Q1: 定时任务没有执行？

**A:** 检查以下几点：
1. 确认 `main.go` 中调用了 `tasks.StartCronJobs()`
2. 查看启动日志是否有"定时任务启动成功"
3. 检查服务器时间是否正确

### Q2: 订单超时后没有自动取消？

**A:** 
1. 确认订单的 `expire_time` 字段有值
2. 等待至少5分钟（定时任务执行周期）
3. 查看后端日志是否有错误信息

### Q3: 退款申请失败？

**A:**
1. 确认订单状态是"已支付"
2. 确认订单未退款过
3. 检查请求参数是否正确

### Q4: 物流查询返回空数据？

**A:**
1. 确认订单已发货（有快递公司和单号）
2. 如果未配置API Key，会返回模拟数据
3. 检查快递公司和单号是否正确

---

## 📝 后续优化建议

1. **性能优化**
   - 为统计查询添加缓存（Redis）
   - 优化大数据量下的分页查询

2. **功能增强**
   - 实现商品评价系统
   - 实现商品收藏功能
   - 实现浏览历史记录
   - 完善微信/支付宝真实退款逻辑

3. **监控告警**
   - 添加定时任务执行监控
   - 退款失败告警
   - 库存异常告警

4. **安全性**
   - 退款金额校验（不能超过订单实付金额）
   - 防止重复退款
   - 接口限流

---

## 📞 技术支持

如有问题，请检查：
1. 后端日志输出
2. 数据库数据是否正确
3. API请求参数是否符合规范

祝使用愉快！🎉
