package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"shop_api/database"
	"shop_api/models"
	"shop_api/types"
	"shop_api/utils"
)

// CreateReview 发表评论
func CreateReview(c *gin.Context) {
	userID := utils.GetUserID(c)

	var input struct {
		OrderID     uint64   `json:"order_id" binding:"required"`
		ProductID   uint64   `json:"product_id" binding:"required"`
		OrderItemID uint64   `json:"order_item_id" binding:"required"`
		Rating      int8     `json:"rating" binding:"required,min=1,max=5"`
		Content     string   `json:"content"`
		Images      []string `json:"images"`
		IsAnonymous int8     `json:"is_anonymous"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 验证订单是否属于该用户且已完成
	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", input.OrderID, userID).First(&order).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	if order.OrderStatus < models.OrderStatusReceived {
		utils.Fail(c, "订单未完成，无法评价")
		return
	}

	// 检查是否已评价
	var existReview models.ProductReview
	err := database.GetDB().Where("order_item_id = ?", input.OrderItemID).First(&existReview).Error
	if err == nil {
		utils.Fail(c, "该商品已评价")
		return
	}

	// 创建评价
	review := models.ProductReview{
		OrderID:     input.OrderID,
		UserID:      userID,
		ProductID:   input.ProductID,
		OrderItemID: input.OrderItemID,
		Rating:      input.Rating,
		Content:     input.Content,
		Images:      input.Images,
		IsAnonymous: input.IsAnonymous,
		Status:      1,
	}

	tx := database.GetDB().Begin()

	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "评价失败")
		return
	}

	// 更新商品评分统计
	var avgRating float64
	var reviewCount int64
	tx.Model(&models.ProductReview{}).
		Where("product_id = ? AND status = 1", input.ProductID).
		Select("AVG(rating)").
		Scan(&avgRating)
	tx.Model(&models.ProductReview{}).
		Where("product_id = ? AND status = 1", input.ProductID).
		Count(&reviewCount)

	tx.Model(&models.Product{}).Where("id = ?", input.ProductID).Updates(map[string]interface{}{
		"avg_rating":   avgRating,
		"review_count": reviewCount,
	})

	tx.Commit()

	utils.Success(c, gin.H{
		"message":   "评价成功",
		"review_id": review.ID,
	})
}

// GetProductReviews 获取商品评论列表
func GetProductReviews(c *gin.Context) {
	productID := c.Param("id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	withImages := c.Query("with_images")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := database.GetDB().Model(&models.ProductReview{}).
		Where("product_id = ? AND status = 1", productID)

	if withImages == "1" {
		query = query.Where("images IS NOT NULL AND images != '[]'")
	}

	var total int64
	query.Count(&total)

	var reviews []models.ProductReview
	if err := query.Preload("User").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reviews).Error; err != nil {
		utils.Fail(c, "获取评论失败")
		return
	}

	utils.PageSuccess(c, reviews, total, page, pageSize)
}

// GetMyReviews 获取我的评论列表
func GetMyReviews(c *gin.Context) {
	userID := utils.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	database.GetDB().Model(&models.ProductReview{}).Where("user_id = ?", userID).Count(&total)

	var reviews []models.ProductReview
	if err := database.GetDB().Preload("Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reviews).Error; err != nil {
		utils.Fail(c, "获取评论失败")
		return
	}

	utils.PageSuccess(c, reviews, total, page, pageSize)
}

// GetReviewStats 获取商品评论统计
func GetReviewStats(c *gin.Context) {
	productID := c.Param("id")

	var stats models.ReviewStats

	// 获取平均评分和总数
	database.GetDB().Model(&models.ProductReview{}).
		Where("product_id = ? AND status = 1", productID).
		Select("AVG(rating) as average_rating, COUNT(*) as total_count").
		Scan(&stats)

	// 获取各星级数量
	var starCounts []struct {
		Rating int8
		Count  int64
	}
	database.GetDB().Model(&models.ProductReview{}).
		Where("product_id = ? AND status = 1", productID).
		Select("rating, COUNT(*) as count").
		Group("rating").
		Scan(&starCounts)

	for _, sc := range starCounts {
		switch sc.Rating {
		case 5:
			stats.FiveStar = sc.Count
		case 4:
			stats.FourStar = sc.Count
		case 3:
			stats.ThreeStar = sc.Count
		case 2:
			stats.TwoStar = sc.Count
		case 1:
			stats.OneStar = sc.Count
		}
	}

	utils.Success(c, stats)
}

// ============ 后台管理接口 ============

// AdminGetReviews 获取评论列表（后台）
func AdminGetReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	productName := c.Query("product_name")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := database.GetDB().Model(&models.ProductReview{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if productName != "" {
		query = query.Joins("JOIN products ON products.id = product_reviews.product_id").
			Where("products.name LIKE ?", "%"+productName+"%")
	}

	var total int64
	query.Count(&total)

	var reviews []models.ProductReview
	if err := query.Preload("User").Preload("Product").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reviews).Error; err != nil {
		utils.Fail(c, "获取评论失败")
		return
	}

	utils.PageSuccess(c, reviews, total, page, pageSize)
}

// AdminUpdateReviewStatus 更新评论状态
func AdminUpdateReviewStatus(c *gin.Context) {
	reviewID := c.Param("id")

	var input struct {
		Status int8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	var review models.ProductReview
	if err := database.GetDB().First(&review, reviewID).Error; err != nil {
		utils.Fail(c, "评论不存在")
		return
	}

	oldStatus := review.Status
	review.Status = input.Status

	if err := database.GetDB().Save(&review).Error; err != nil {
		utils.Fail(c, "更新失败")
		return
	}

	// 如果状态改变，更新商品评分统计
	if oldStatus != input.Status {
		var avgRating float64
		var reviewCount int64
		database.GetDB().Model(&models.ProductReview{}).
			Where("product_id = ? AND status = 1", review.ProductID).
			Select("AVG(rating)").
			Scan(&avgRating)
		database.GetDB().Model(&models.ProductReview{}).
			Where("product_id = ? AND status = 1", review.ProductID).
			Count(&reviewCount)

		database.GetDB().Model(&models.Product{}).Where("id = ?", review.ProductID).Updates(map[string]interface{}{
			"avg_rating":   avgRating,
			"review_count": reviewCount,
		})
	}

	utils.Success(c, gin.H{
		"message": "更新成功",
	})
}

// AdminReplyReview 商家回复评论
func AdminReplyReview(c *gin.Context) {
	reviewID := c.Param("id")

	var input struct {
		ReplyContent string `json:"reply_content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "请输入回复内容")
		return
	}

	var review models.ProductReview
	if err := database.GetDB().First(&review, reviewID).Error; err != nil {
		utils.Fail(c, "评论不存在")
		return
	}

	now := types.Now()
	review.ReplyContent = input.ReplyContent
	review.ReplyTime = &now

	if err := database.GetDB().Save(&review).Error; err != nil {
		utils.Fail(c, "回复失败")
		return
	}

	utils.Success(c, gin.H{
		"message": "回复成功",
	})
}

// AdminDeleteReview 删除评论
func AdminDeleteReview(c *gin.Context) {
	reviewID := c.Param("id")

	var review models.ProductReview
	if err := database.GetDB().First(&review, reviewID).Error; err != nil {
		utils.Fail(c, "评论不存在")
		return
	}

	tx := database.GetDB().Begin()

	if err := tx.Delete(&review).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "删除失败")
		return
	}

	// 更新商品评分统计
	var avgRating float64
	var reviewCount int64
	tx.Model(&models.ProductReview{}).
		Where("product_id = ? AND status = 1", review.ProductID).
		Select("AVG(rating)").
		Scan(&avgRating)
	tx.Model(&models.ProductReview{}).
		Where("product_id = ? AND status = 1", review.ProductID).
		Count(&reviewCount)

	tx.Model(&models.Product{}).Where("id = ?", review.ProductID).Updates(map[string]interface{}{
		"avg_rating":   avgRating,
		"review_count": reviewCount,
	})

	tx.Commit()

	utils.Success(c, gin.H{
		"message": "删除成功",
	})
}

// GetCanReviewOrders 获取可评价的订单列表
func GetCanReviewOrders(c *gin.Context) {
	userID := utils.GetUserID(c)

	// 查询已完成但未评价的订单
	var orders []models.Order
	database.GetDB().Where("user_id = ? AND order_status >= ?", userID, models.OrderStatusReceived).
		Order("created_at DESC").
		Find(&orders)

	var canReviewItems []gin.H

	for _, order := range orders {
		var items []models.OrderItem
		database.GetDB().Where("order_id = ?", order.ID).Find(&items)

		for _, item := range items {
			// 检查是否已评价
			var count int64
			database.GetDB().Model(&models.ProductReview{}).
				Where("order_item_id = ?", item.ID).
				Count(&count)

			if count == 0 {
				canReviewItems = append(canReviewItems, gin.H{
					"order_id":      order.ID,
					"order_no":      order.OrderNo,
					"order_item_id": item.ID,
					"product_id":    item.ProductID,
					"product_name":  item.ProductName,
					"product_image": item.ProductImage,
				})
			}
		}
	}

	utils.Success(c, canReviewItems)
}
