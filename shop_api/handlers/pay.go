package handlers

import (
	"io"
	"net/http"
	"shop_api/database"
	"shop_api/models"
	"shop_api/services"
	"shop_api/types"
	"shop_api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PayInput struct {
	OrderID uint64 `json:"order_id" binding:"required"`
	PayType int8   `json:"pay_type" binding:"required"`
}

func GetPayURL(c *gin.Context) {
	userID := utils.GetUserID(c)

	var input PayInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	// 使用事务和行锁防止并发支付
	tx := database.GetDB().Begin()

	var order models.Order
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ? AND user_id = ?", input.OrderID, userID).First(&order).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "订单不存在")
		return
	}

	if order.PayStatus != models.PayStatusUnpaid {
		tx.Rollback()
		utils.Fail(c, "订单已支付，请勿重复提交")
		return
	}

	// 如果是模拟支付，直接在事务中处理
	if input.PayType == models.PayTypeWechat || input.PayType == models.PayTypeAlipay {
		// 检查是否配置了真实支付
		payService := services.GetPayService()
		isConfigured := false

		if input.PayType == models.PayTypeWechat {
			isConfigured = payService.IsWechatConfigured()
		} else {
			isConfigured = payService.IsAlipayConfigured()
		}

		if !isConfigured {
			// 未配置，使用模拟支付，在事务中处理
			tx.Rollback() // 释放锁
			mockPayAndRespond(c, &order, input.PayType)
			return
		}
	}

	tx.Rollback() // 释放锁

	var payURL string
	var err error

	payService := services.GetPayService()

	switch input.PayType {
	case models.PayTypeWechat:
		payURL, err = payService.GetWechatPayURL(&order)
	case models.PayTypeAlipay:
		payURL, err = payService.GetAlipayURL(&order)
	default:
		utils.Fail(c, "不支持的支付方式")
		return
	}

	if err != nil {
		utils.Error("获取支付链接失败: %v", err)
		utils.Fail(c, "获取支付链接失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"pay_url":  payURL,
		"order_no": order.OrderNo,
	})
}

// mockPayAndRespond 模拟支付并返回结果
func mockPayAndRespond(c *gin.Context, order *models.Order, payType int8) {
	userID := utils.GetUserID(c)

	// 重新加载订单（包含订单项）
	var fullOrder models.Order
	if err := database.GetDB().Preload("Items").Where("id = ? AND user_id = ?", order.ID, userID).First(&fullOrder).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	tx := database.GetDB().Begin()
	now := types.Now()

	fullOrder.PayStatus = models.PayStatusPaid
	fullOrder.PayType = payType
	fullOrder.PayTime = &now
	fullOrder.OrderStatus = models.OrderStatusPaid

	// 扣减库存并增加销量
	for _, item := range fullOrder.Items {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			utils.Fail(c, "商品不存在")
			return
		}
		if product.Stock < item.Quantity {
			tx.Rollback()
			utils.Fail(c, "商品["+item.ProductName+"]库存不足")
			return
		}
		if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			utils.Fail(c, "扣减库存失败")
			return
		}
		if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).UpdateColumn("sales", gorm.Expr("sales + ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			utils.Fail(c, "更新销量失败")
			return
		}
	}

	if err := tx.Save(&fullOrder).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "更新订单失败")
		return
	}

	// 记录支付日志
	payLog := models.PayLog{
		OrderID:    fullOrder.ID,
		OrderNo:    fullOrder.OrderNo,
		TradeNo:    "mock_" + fullOrder.OrderNo,
		PayType:    payType,
		PayStatus:  models.PayStatusPaid,
		PayAmount:  fullOrder.PayAmount,
		NotifyData: `{"mock":true,"message":"模拟支付成功"}`,
	}
	if err := tx.Create(&payLog).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "记录支付日志失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Fail(c, "模拟支付失败")
		return
	}

	utils.Info("模拟支付成功: 订单号=%s, 支付方式=%d", fullOrder.OrderNo, payType)

	utils.Success(c, gin.H{
		"order_id":     fullOrder.ID,
		"order_no":     fullOrder.OrderNo,
		"pay_status":   fullOrder.PayStatus,
		"order_status": fullOrder.OrderStatus,
		"mock_pay":     true,
		"message":      "模拟支付成功（开发模式）",
	})
}

func WechatNotify(c *gin.Context) {
	data, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	payService := services.GetPayService()
	if err := payService.WechatNotify(data); err != nil {
		c.String(http.StatusBadRequest, "FAIL")
		return
	}

	c.String(http.StatusOK, "SUCCESS")
}

func AlipayNotify(c *gin.Context) {
	data := make(map[string]string)
	for key, values := range c.Request.PostForm {
		data[key] = values[0]
	}

	payService := services.GetPayService()
	if err := payService.AlipayNotify(data); err != nil {
		c.String(http.StatusBadRequest, "FAIL")
		return
	}

	c.String(http.StatusOK, "SUCCESS")
}

func QueryPayStatus(c *gin.Context) {
	userID := utils.GetUserID(c)
	orderID := c.Param("id")

	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	utils.Success(c, gin.H{
		"pay_status":   order.PayStatus,
		"order_status": order.OrderStatus,
	})
}

func MockPaySuccess(c *gin.Context) {
	userID := utils.GetUserID(c)
	orderID := c.Param("id")

	var order models.Order
	if err := database.GetDB().Preload("Items").Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	if order.PayStatus != models.PayStatusUnpaid {
		utils.Fail(c, "订单已支付")
		return
	}

	tx := database.GetDB().Begin()
	now := types.Now()

	order.PayStatus = models.PayStatusPaid
	order.PayType = models.PayTypeWechat
	order.PayTime = &now
	order.OrderStatus = models.OrderStatusPaid

	for _, item := range order.Items {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			utils.Fail(c, "商品不存在")
			return
		}
		if product.Stock < item.Quantity {
			tx.Rollback()
			utils.Fail(c, "商品["+item.ProductName+"]库存不足")
			return
		}
		if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).UpdateColumn("stock", tx.Raw("stock - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			utils.Fail(c, "扣减库存失败")
			return
		}
		if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).UpdateColumn("sales", tx.Raw("sales + ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			utils.Fail(c, "更新销量失败")
			return
		}
	}

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "更新订单失败")
		return
	}

	payLog := models.PayLog{
		OrderID:    order.ID,
		OrderNo:    order.OrderNo,
		TradeNo:    "mock_" + order.OrderNo,
		PayType:    models.PayTypeWechat,
		PayStatus:  models.PayStatusPaid,
		PayAmount:  order.PayAmount,
		NotifyData: `{"mock":true,"message":"dev pay success"}`,
	}
	if err := tx.Create(&payLog).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "记录支付日志失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Fail(c, "模拟支付失败")
		return
	}

	utils.Success(c, gin.H{
		"order_id":     order.ID,
		"order_no":     order.OrderNo,
		"pay_status":   order.PayStatus,
		"order_status": order.OrderStatus,
	})
}
