package handlers

import (
	"shop_api/database"
	"shop_api/models"
	"shop_api/services"
	"shop_api/types"
	"shop_api/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ApplyRefund 申请退款
func ApplyRefund(c *gin.Context) {
	userID := utils.GetUserID(c)

	var input struct {
		OrderID      uint64   `json:"order_id" binding:"required"`
		RefundType   string   `json:"refund_type" binding:"required,oneof=refund_only exchange"`
		Reason       string   `json:"reason" binding:"required"`
		Images       []string `json:"images"`
		RefundAmount float64  `json:"refund_amount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 验证订单
	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", input.OrderID, userID).First(&order).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	// 检查订单状态是否可以退款
	if order.PayStatus != models.PayStatusPaid {
		utils.Fail(c, "订单未支付，无法申请退款")
		return
	}

	if order.OrderStatus >= models.OrderStatusRefunded {
		utils.Fail(c, "订单已退款")
		return
	}

	// 如果未指定退款金额，默认为订单实付金额
	refundAmount := input.RefundAmount
	if refundAmount <= 0 {
		refundAmount = order.PayAmount
	}

	// 创建退款申请
	refund := models.RefundApplication{
		OrderID:      input.OrderID,
		UserID:       userID,
		RefundType:   input.RefundType,
		Reason:       input.Reason,
		Images:       input.Images,
		RefundAmount: refundAmount,
		Status:       models.RefundStatusPending,
	}

	if err := database.GetDB().Create(&refund).Error; err != nil {
		utils.Fail(c, "申请退款失败")
		return
	}

	utils.Success(c, gin.H{
		"refund_id": refund.ID,
		"message":   "退款申请已提交，请等待审核",
	})
}

// GetMyRefunds 获取我的退款申请列表
func GetMyRefunds(c *gin.Context) {
	userID := utils.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := database.GetDB().Model(&models.RefundApplication{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var refunds []models.RefundApplication
	if err := query.Preload("Order").Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&refunds).Error; err != nil {
		utils.Fail(c, "获取退款列表失败")
		return
	}

	utils.PageSuccess(c, refunds, total, page, pageSize)
}

// GetRefundDetail 获取退款详情
func GetRefundDetail(c *gin.Context) {
	userID := utils.GetUserID(c)
	refundID := c.Param("id")

	var refund models.RefundApplication
	if err := database.GetDB().Preload("Order").Preload("User").Preload("Handler").
		First(&refund, refundID).Error; err != nil {
		utils.Fail(c, "退款申请不存在")
		return
	}

	// 权限检查
	if refund.UserID != userID {
		utils.Forbidden(c, "无权查看此退款申请")
		return
	}

	utils.Success(c, refund)
}

// AdminGetRefunds 后台获取退款申请列表
func AdminGetRefunds(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	orderNo := c.Query("order_no")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := database.GetDB().Model(&models.RefundApplication{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if orderNo != "" {
		query = query.Joins("JOIN orders ON orders.id = refund_applications.order_id").
			Where("orders.order_no LIKE ?", "%"+orderNo+"%")
	}

	var total int64
	query.Count(&total)

	var refunds []models.RefundApplication
	if err := query.Preload("Order").Preload("User").Preload("Handler").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&refunds).Error; err != nil {
		utils.Fail(c, "获取退款列表失败")
		return
	}

	utils.PageSuccess(c, refunds, total, page, pageSize)
}

// AdminGetRefundDetail 后台获取退款详情
func AdminGetRefundDetail(c *gin.Context) {
	refundID := c.Param("id")

	var refund models.RefundApplication
	if err := database.GetDB().Preload("Order").Preload("User").Preload("Handler").
		First(&refund, refundID).Error; err != nil {
		utils.Fail(c, "退款申请不存在")
		return
	}

	utils.Success(c, refund)
}

// AdminApproveRefund 审核通过退款
func AdminApproveRefund(c *gin.Context) {
	adminID := utils.GetUserID(c)
	refundID := c.Param("id")

	var input struct {
		HandlerReply string `json:"handler_reply"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		input.HandlerReply = "审核通过"
	}

	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var refund models.RefundApplication
	if err := tx.Preload("Order").First(&refund, refundID).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "退款申请不存在")
		return
	}

	if refund.Status != models.RefundStatusPending {
		tx.Rollback()
		utils.Fail(c, "退款申请状态不正确")
		return
	}

	// 更新退款状态为退款中
	now := types.LocalTime(time.Now())
	refund.Status = models.RefundStatusRefunding
	refund.HandlerID = &adminID
	refund.HandlerReply = input.HandlerReply

	if err := tx.Save(&refund).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "更新退款状态失败")
		return
	}

	// 调用支付服务进行退款
	payService := services.GetPayService()
	var tradeNo string
	var err error

	if refund.Order.PayType == models.PayTypeWechat {
		tradeNo, err = payService.WechatRefund(refund.Order, refund.RefundAmount)
	} else if refund.Order.PayType == models.PayTypeAlipay {
		tradeNo, err = payService.AlipayRefund(refund.Order, refund.RefundAmount)
	} else {
		// 模拟支付，直接成功
		tradeNo = "mock_refund_" + refund.Order.OrderNo
		err = nil
	}

	if err != nil {
		tx.Rollback()
		utils.Fail(c, "退款失败: "+err.Error())
		return
	}

	// 更新退款状态为已退款
	refund.Status = models.RefundStatusRefunded
	refund.RefundTradeNo = tradeNo
	refund.RefundedAt = &now

	if err := tx.Save(&refund).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "更新退款状态失败")
		return
	}

	// 更新订单状态
	if err := tx.Model(refund.Order).Update("order_status", models.OrderStatusRefunded).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "更新订单状态失败")
		return
	}

	// 如果商品未发货，恢复库存
	if refund.Order.OrderStatus < models.OrderStatusShipped {
		var items []models.OrderItem
		tx.Where("order_id = ?", refund.Order.ID).Find(&items)
		for _, item := range items {
			tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
				UpdateColumn("stock", gorm.Expr("stock + ?", item.Quantity))
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.Fail(c, "退款处理失败")
		return
	}

	utils.Success(c, gin.H{
		"message": "退款成功",
	})
}

// AdminRejectRefund 拒绝退款申请
func AdminRejectRefund(c *gin.Context) {
	adminID := utils.GetUserID(c)
	refundID := c.Param("id")

	var input struct {
		HandlerReply string `json:"handler_reply" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "请输入拒绝原因")
		return
	}

	var refund models.RefundApplication
	if err := database.GetDB().First(&refund, refundID).Error; err != nil {
		utils.Fail(c, "退款申请不存在")
		return
	}

	if refund.Status != models.RefundStatusPending {
		utils.Fail(c, "退款申请状态不正确")
		return
	}

	now := types.Now()
	database.GetDB().Model(&refund).Updates(map[string]interface{}{
		"status":        models.RefundStatusRejected,
		"handler_id":    adminID,
		"handler_reply": input.HandlerReply,
		"updated_at":    now,
	})

	utils.Success(c, gin.H{
		"message": "已拒绝退款申请",
	})
}
