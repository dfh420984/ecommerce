package handlers

import (
	"io"
	"net/http"
	"shop_api/database"
	"shop_api/models"
	"shop_api/services"
	"shop_api/utils"
	"time"

	"github.com/gin-gonic/gin"
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

	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", input.OrderID, userID).First(&order).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	if order.PayStatus != models.PayStatusUnpaid {
		utils.Fail(c, "订单已支付")
		return
	}

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
		utils.Fail(c, "获取支付链接失败")
		return
	}

	utils.Success(c, gin.H{
		"pay_url":  payURL,
		"order_no": order.OrderNo,
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
	now := time.Now()

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

func ApplyRefund(c *gin.Context) {
	userID := utils.GetUserID(c)
	orderID := c.Param("id")

	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	if order.PayStatus != models.PayStatusPaid {
		utils.Fail(c, "该订单无法申请退款")
		return
	}

	order.OrderStatus = models.OrderStatusRefund
	database.GetDB().Save(&order)

	utils.Success(c, nil)
}
