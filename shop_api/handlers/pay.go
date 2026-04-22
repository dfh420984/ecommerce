package handlers

import (
	"io"
	"net/http"
	"shop_api/database"
	"shop_api/models"
	"shop_api/services"
	"shop_api/utils"

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
		"pay_url": payURL,
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
