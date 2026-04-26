# 订单自动完成功能说明

## 📋 功能概述

系统已实现订单自动完成功能，当用户确认收货后，如果在指定天数内没有申请退款，订单将自动从"已收货"状态变更为"已完成"状态。

## ⚙️ 配置说明

### 配置项名称
`auto_complete_days`

### 默认值
`7` 天

### 配置位置
后台管理系统 → 系统配置

### 如何修改
1. 登录后台管理系统
2. 进入"系统配置"页面
3. 找到 `auto_complete_days` 配置项
4. 点击"编辑"按钮
5. 修改配置值为需要的天数（建议 3-15 天）
6. 保存配置

### 配置示例
```
配置名称: auto_complete_days
配置值: 7
描述: 订单自动完成天数（用户确认收货后多少天自动完成）
```

## 🔄 工作流程

### 定时任务执行时间
- **执行频率**：每天凌晨 2:00
- **首次执行**：系统启动后的下一个凌晨 2:00

### 自动完成条件
订单必须同时满足以下条件：
1. 订单状态为 `4`（已收货）
2. 订单最后更新时间超过配置的天数（默认 7 天）
3. 订单未申请退款

### 状态变更流程
```
用户确认收货 (order_status = 4)
    ↓
等待 N 天（可配置，默认 7 天）
    ↓
定时任务检查（每天凌晨 2:00）
    ↓
自动完成订单 (order_status = 5, complete_time 记录)
    ↓
用户可以评价商品
```

## 💻 技术实现

### 后端文件

#### 1. 服务层
**文件**: `shop_api/services/order_timeout.go`

**主要方法**:
- `CheckAndAutoCompleteOrders()` - 检查并自动完成订单
- `completeOrder()` - 完成单个订单

**核心逻辑**:
```go
// 从配置中读取自动完成天数
autoCompleteDays := 7
var config models.Config
if err := tx.Where("name = ?", "auto_complete_days").First(&config).Error; err == nil {
    if days, err := strconv.Atoi(config.Value); err == nil && days > 0 {
        autoCompleteDays = days
    }
}

// 查找已收货超过指定天数的订单
thresholdTime := now.AddDate(0, 0, -autoCompleteDays)
err := tx.Where("order_status = ? AND updated_at < ?",
    models.OrderStatusReceived, thresholdTime).Find(&orders).Error
```

#### 2. 定时任务
**文件**: `shop_api/tasks/cron_jobs.go`

**任务配置**:
```go
// 每天凌晨2点检查并自动完成已收货超时的订单
go func() {
    for {
        now := time.Now()
        // 计算到下一个凌晨2点的时间
        next := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
        if now.After(next) {
            next = next.Add(24 * time.Hour)
        }
        duration := next.Sub(now)
        
        log.Printf("下次自动完成订单检查将在 %v 后执行", duration)
        time.Sleep(duration)
        
        log.Println("执行定时任务：自动完成订单")
        services.GetOrderTimeoutService().CheckAndAutoCompleteOrders()
    }
}()
```

### 数据库配置

**表**: `configs`

**SQL 插入语句**:
```sql
INSERT INTO `configs` (`name`, `value`, `description`) VALUES
('auto_complete_days', '7', '订单自动完成天数（用户确认收货后多少天自动完成）');
```

## 📊 日志输出

### 正常日志
```
开始检查自动完成订单...
找到 5 个需要自动完成的订单（已收货超过7天），开始处理...
订单 123 已自动完成
订单 124 已自动完成
订单 125 已自动完成
订单 126 已自动完成
订单 127 已自动完成
成功自动完成 5 个订单
```

### 无订单日志
```
开始检查自动完成订单...
没有需要自动完成的订单
```

### 错误日志
```
开始检查自动完成订单...
查询自动完成订单失败: [错误信息]
```

## 🔧 自定义配置

### 修改自动完成天数

#### 方式1：通过后台管理界面（推荐）
1. 登录后台管理系统
2. 进入"系统配置"
3. 编辑 `auto_complete_days` 配置项
4. 修改值为所需天数
5. 保存即可生效（下次定时任务执行时生效）

#### 方式2：直接修改数据库
```sql
UPDATE configs SET value = '10' WHERE name = 'auto_complete_days';
```

#### 方式3：修改代码默认值
在 `shop_api/services/order_timeout.go` 中修改：
```go
autoCompleteDays := 7  // 修改这里的数字
```

### 修改定时任务执行时间

如果需要修改定时任务的执行时间，编辑 `shop_api/tasks/cron_jobs.go`：

```go
// 当前是凌晨2点执行
next := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())

// 改为凌晨3点执行
next := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
```

### 修改执行频率

如果需要更频繁地检查（例如每小时检查一次）：

```go
// 每小时检查一次
go func() {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        log.Println("执行定时任务：自动完成订单")
        services.GetOrderTimeoutService().CheckAndAutoCompleteOrders()
    }
}()
```

## ✅ 测试验证

### 手动测试步骤

1. **创建测试订单**
   ```
   小程序端：选择商品 → 下单 → 支付
   ```

2. **模拟发货**
   ```
   后台管理：订单管理 → 找到订单 → 点击"发货" → 填写快递信息
   ```

3. **确认收货**
   ```
   小程序端：我的订单 → 找到订单 → 点击"确认收货"
   ```
   
   此时订单状态变为 `4`（已收货）

4. **等待自动完成**
   - 方式1：等待 7 天，定时任务会自动完成
   - 方式2：临时修改配置为 1 天，第二天查看结果
   - 方式3：手动调用接口测试（见下方）

5. **验证结果**
   ```
   小程序端：我的订单 → 查看订单状态应为"已完成"
   后台管理：订单管理 → 查看订单状态应为"已完成"
   ```

### 手动触发测试（开发环境）

可以创建一个临时接口来手动触发自动完成检查：

```go
// 在 routes/router.go 中添加测试路由
adminAuth.POST("/orders/auto-complete", handlers.ManualAutoCompleteOrders)

// 在 handlers/order.go 中添加处理方法
func ManualAutoCompleteOrders(c *gin.Context) {
    services.GetOrderTimeoutService().CheckAndAutoCompleteOrders()
    utils.Success(c, gin.H{
        "message": "自动完成检查已执行，请查看日志",
    })
}
```

然后调用：
```bash
curl -X POST http://localhost:8080/api/admin/orders/auto-complete \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 🎯 业务场景

### 适用场景
1. **标准电商流程**：给用户留出售后时间窗口
2. **自动化运营**：减少人工操作，提高效率
3. **评价引导**：订单完成后自动开放评价入口

### 优势
- ✅ 自动化处理，无需人工干预
- ✅ 可配置天数，灵活适应不同业务需求
- ✅ 定时执行，系统负载低
- ✅ 事务保证，数据一致性
- ✅ 详细日志，便于排查问题

### 注意事项
- ⚠️ 确保用户有足够的时间申请售后（建议至少 7 天）
- ⚠️ 自动完成后用户仍可以评价商品
- ⚠️ 如果用户申请退款，订单不会自动完成
- ⚠️ 定时任务只在凌晨 2 点执行，不是实时完成

## 🐛 常见问题

### Q1: 为什么订单没有自动完成？
**A**: 检查以下几点：
1. 订单状态是否为 4（已收货）
2. 订单最后更新时间是否超过配置天数
3. 定时任务是否正常运行（查看日志）
4. 配置项 `auto_complete_days` 是否正确设置

### Q2: 如何查看定时任务是否运行？
**A**: 查看后端日志：
```
grep "自动完成订单" logs/app.log
```

应该能看到类似这样的日志：
```
2026/04/26 02:00:00 执行定时任务：自动完成订单
2026/04/26 02:00:01 开始检查自动完成订单...
```

### Q3: 可以立即完成订单吗？
**A**: 目前只能等待定时任务执行。如果需要立即完成，可以：
1. 在后台管理系统手动点击"完成"按钮
2. 或者添加一个手动触发的接口（见上方测试部分）

### Q4: 自动完成会影响退款吗？
**A**: 不会。即使用户申请退款，订单也不会自动完成。只有状态为"已收货"且未申请退款的订单才会被自动完成。

### Q5: 如何调整自动完成的时间？
**A**: 
- 短期调整：修改数据库配置 `UPDATE configs SET value = 'X' WHERE name = 'auto_complete_days';`
- 长期调整：在后台管理系统中修改配置

## 📝 更新日志

### v1.0.0 (2026-04-26)
- ✅ 实现订单自动完成功能
- ✅ 添加定时任务（每天凌晨 2 点执行）
- ✅ 支持配置自动完成天数
- ✅ 添加详细日志输出
- ✅ 事务保证数据一致性

## 🔗 相关文件

- 服务层: `shop_api/services/order_timeout.go`
- 定时任务: `shop_api/tasks/cron_jobs.go`
- 路由配置: `shop_api/routes/router.go`
- 数据库脚本: `shop_api/database/shop.sql`
- 后台管理: `shop_admin/src/views/Configs.vue`

---

**最后更新**: 2026-04-26  
**维护者**: 开发团队
