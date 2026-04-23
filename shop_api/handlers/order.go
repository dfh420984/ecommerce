package handlers

import (
	"gorm.io/gorm"
	"shop_api/database"
	"shop_api/models"
	"shop_api/types"
	"shop_api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateOrderInput struct {
	AddressID uint64   `json:"address_id" binding:"required"`
	Remark    string   `json:"remark"`
	CartIDs   []uint64 `json:"cart_ids"`
	ProductID uint64   `json:"product_id"`
	Quantity  int      `json:"quantity"`
}

func CreateOrder(c *gin.Context) {
	userID := utils.GetUserID(c)

	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	var address models.Address
	if err := database.GetDB().Where("id = ? AND user_id = ?", input.AddressID, userID).First(&address).Error; err != nil {
		utils.Fail(c, "收货地址不存在")
		return
	}

	var items []models.OrderItem
	var totalAmount float64
	var carts []models.Cart // 提前声明，用于后续减少库存

	if len(input.CartIDs) > 0 {
		if err := database.GetDB().Preload("Product").Where("id IN ? AND user_id = ? AND selected = ?", input.CartIDs, userID, 1).Find(&carts).Error; err != nil {
			utils.Fail(c, "购物车商品不存在")
			return
		}

		if len(carts) == 0 {
			utils.Fail(c, "请选择购物车商品")
			return
		}

		for _, cart := range carts {
			if cart.Product == nil {
				continue
			}
			if cart.Product.Stock < cart.Quantity {
				utils.Fail(c, "商品["+cart.Product.Name+"]库存不足")
				return
			}

			item := models.OrderItem{
				ProductID:    cart.ProductID,
				ProductName:  cart.Product.Name,
				ProductImage: "",
				Price:        cart.Product.Price,
				Quantity:     cart.Quantity,
				Subtotal:     cart.Product.Price * float64(cart.Quantity),
			}
			if len(cart.Product.Images) > 0 {
				item.ProductImage = cart.Product.Images[0]
			}
			items = append(items, item)
			totalAmount += item.Subtotal
		}
	} else if input.ProductID > 0 && input.Quantity > 0 {
		var product models.Product
		if err := database.GetDB().First(&product, input.ProductID).Error; err != nil {
			utils.Fail(c, "商品不存在")
			return
		}

		if product.Stock < input.Quantity {
			utils.Fail(c, "库存不足")
			return
		}

		item := models.OrderItem{
			ProductID:    product.ID,
			ProductName:  product.Name,
			ProductImage: "",
			Price:        product.Price,
			Quantity:     input.Quantity,
			Subtotal:     product.Price * float64(input.Quantity),
		}
		if len(product.Images) > 0 {
			item.ProductImage = product.Images[0]
		}
		items = append(items, item)
		totalAmount = item.Subtotal
	} else {
		utils.Fail(c, "参数错误")
		return
	}

	order := models.Order{
		OrderNo:        utils.GenerateOrderNo(),
		UserID:         userID,
		OrderStatus:    models.OrderStatusPending,
		PayStatus:      models.PayStatusUnpaid,
		TotalAmount:    totalAmount,
		DiscountAmount: 0,
		FreightAmount:  0,
		PayAmount:      totalAmount,
		Consignee:      address.Consignee,
		Phone:          address.Phone,
		Province:       address.Province,
		City:           address.City,
		District:       address.District,
		Address:        address.Address,
		Remark:         input.Remark,
	}

	tx := database.GetDB().Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		utils.Error("创建订单失败: %v", err)
		utils.Fail(c, "创建订单失败: "+err.Error())
		return
	}

	for i := range items {
		items[i].OrderID = order.ID
	}

	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		utils.Error("创建订单明细失败: %v", err)
		utils.Fail(c, "创建订单明细失败: "+err.Error())
		return
	}

	// 减少商品库存
	for _, cart := range carts {
		if cart.Product != nil {
			if err := tx.Model(&models.Product{}).Where("id = ?", cart.ProductID).UpdateColumn("stock", gorm.Expr("stock - ?", cart.Quantity)).Error; err != nil {
				tx.Rollback()
				utils.Error("更新库存失败: %v", err)
				utils.Fail(c, "更新库存失败")
				return
			}
		}
	}

	// 使用事务删除购物车
	if len(input.CartIDs) > 0 {
		if err := tx.Where("id IN ? AND user_id = ?", input.CartIDs, userID).Delete(&models.Cart{}).Error; err != nil {
			tx.Rollback()
			utils.Error("删除购物车失败: %v", err)
			// 删除购物车失败不影响订单创建，只记录日志
			utils.Info("删除购物车失败，但订单已创建: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error("事务提交失败: %v", err)
		utils.Fail(c, "订单创建失败，请稍后重试")
		return
	}

	utils.Success(c, gin.H{
		"order_id":   order.ID,
		"order_no":   order.OrderNo,
		"pay_amount": order.PayAmount,
	})
}

func GetOrders(c *gin.Context) {
	userID := utils.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := database.GetDB().Model(&models.Order{}).Where("user_id = ?", userID)

	if status > 0 {
		query = query.Where("order_status = ?", status)
	}

	var total int64
	query.Count(&total)

	var orders []models.Order
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders).Error; err != nil {
		utils.Fail(c, "获取订单失败")
		return
	}

	for i := range orders {
		var items []models.OrderItem
		database.GetDB().Where("order_id = ?", orders[i].ID).Find(&items)
		orders[i].Items = items
	}

	utils.PageSuccess(c, orders, total, page, pageSize)
}

func GetOrder(c *gin.Context) {
	userID := utils.GetUserID(c)
	orderID := c.Param("id")

	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	var items []models.OrderItem
	database.GetDB().Where("order_id = ?", order.ID).Find(&items)
	order.Items = items

	utils.Success(c, order)
}

func CancelOrder(c *gin.Context) {
	userID := utils.GetUserID(c)
	orderID := c.Param("id")

	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	if order.OrderStatus != models.OrderStatusPending {
		utils.Fail(c, "当前状态无法取消订单")
		return
	}

	order.OrderStatus = models.OrderStatusCancelled
	now := types.Now()
	order.CancelTime = &now

	database.GetDB().Save(&order)

	utils.Success(c, nil)
}

func ConfirmReceive(c *gin.Context) {
	userID := utils.GetUserID(c)
	orderID := c.Param("id")

	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	if order.OrderStatus != models.OrderStatusShipped {
		utils.Fail(c, "当前状态无法确认收货")
		return
	}

	order.OrderStatus = models.OrderStatusReceived
	now := types.Now()
	order.CompleteTime = &now

	database.GetDB().Save(&order)

	utils.Success(c, nil)
}

func DeleteOrder(c *gin.Context) {
	userID := utils.GetUserID(c)
	orderID := c.Param("id")

	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	if order.OrderStatus != models.OrderStatusCompleted && order.OrderStatus != models.OrderStatusCancelled {
		utils.Fail(c, "当前状态无法删除订单")
		return
	}

	database.GetDB().Delete(&order)
	database.GetDB().Where("order_id = ?", order.ID).Delete(&models.OrderItem{})

	utils.Success(c, nil)
}

// 以下后台管理接口

// AdminGetOrders 后台获取订单列表
func AdminGetOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	orderNo := c.Query("order_no")
	orderStatus, _ := strconv.Atoi(c.DefaultQuery("order_status", "0"))
	payStatus, _ := strconv.Atoi(c.DefaultQuery("pay_status", "0"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := database.GetDB().Model(&models.Order{})

	// 订单号搜索
	if orderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+orderNo+"%")
	}

	// 订单状态筛选
	if orderStatus > 0 {
		query = query.Where("order_status = ?", orderStatus)
	}

	// 支付状态筛选
	if payStatus > 0 {
		query = query.Where("pay_status = ?", payStatus)
	}

	var total int64
	query.Count(&total)

	var orders []models.Order
	if err := query.Preload("User").Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders).Error; err != nil {
		utils.Fail(c, "获取订单失败")
		return
	}

	// 加载订单项
	for i := range orders {
		var items []models.OrderItem
		database.GetDB().Where("order_id = ?", orders[i].ID).Find(&items)
		orders[i].Items = items
	}

	utils.PageSuccess(c, orders, total, page, pageSize)
}

// AdminGetOrder 后台获取订单详情
func AdminGetOrder(c *gin.Context) {
	orderID := c.Param("id")

	var order models.Order
	if err := database.GetDB().Preload("User").First(&order, orderID).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	var items []models.OrderItem
	database.GetDB().Where("order_id = ?", order.ID).Find(&items)
	order.Items = items

	utils.Success(c, order)
}

// ShipOrder 发货
func ShipOrder(c *gin.Context) {
	orderID := c.Param("id")

	var input struct {
		ExpressCompany string `json:"express_company"` // 快递公司
		ExpressNo      string `json:"express_no"`      // 快递单号
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	var order models.Order
	if err := database.GetDB().First(&order, orderID).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	if order.OrderStatus != models.OrderStatusPaid {
		utils.Fail(c, "当前状态无法发货")
		return
	}

	order.OrderStatus = models.OrderStatusShipped
	order.ExpressCompany = input.ExpressCompany
	order.ExpressNo = input.ExpressNo

	database.GetDB().Save(&order)

	utils.Success(c, order)
}

// UpdateOrderStatus 更新订单状态
func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var input struct {
		OrderStatus int8 `json:"order_status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	var order models.Order
	if err := database.GetDB().First(&order, orderID).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	order.OrderStatus = input.OrderStatus
	now := types.Now()

	// 如果状态变为已完成，记录完成时间
	if input.OrderStatus == models.OrderStatusCompleted {
		order.CompleteTime = &now
	}

	database.GetDB().Save(&order)

	utils.Success(c, order)
}
