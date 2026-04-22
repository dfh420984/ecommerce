package handlers

import (
	"shop_api/database"
	"shop_api/models"
	"shop_api/utils"
	"strconv"
	"time"

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

	if len(input.CartIDs) > 0 {
		var carts []models.Cart
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
		utils.Fail(c, "创建订单失败")
		return
	}

	for i := range items {
		items[i].OrderID = order.ID
	}

	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "创建订单明细失败")
		return
	}

	if len(input.CartIDs) > 0 {
		database.GetDB().Where("id IN ? AND user_id = ?", input.CartIDs, userID).Delete(&models.Cart{})
	}

	tx.Commit()

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
	now := time.Now()
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
	now := time.Now()
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
