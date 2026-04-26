package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"shop_api/database"
	"shop_api/models"
	"shop_api/utils"
)

// LikeProduct 点赞商品
func LikeProduct(c *gin.Context) {
	userID := utils.GetUserID(c)
	productID := c.Param("id")

	// 检查商品是否存在
	var product models.Product
	if err := database.GetDB().First(&product, productID).Error; err != nil {
		utils.Fail(c, "商品不存在")
		return
	}

	// 检查是否已点赞
	var like models.ProductLike
	err := database.GetDB().Where("product_id = ? AND user_id = ?", productID, userID).First(&like).Error
	if err == nil {
		utils.Fail(c, "已经点赞过该商品")
		return
	}

	// 创建点赞记录
	like = models.ProductLike{
		ProductID: product.ID,
		UserID:    userID,
	}

	if err := database.GetDB().Create(&like).Error; err != nil {
		utils.Fail(c, "点赞失败")
		return
	}

	// 更新商品点赞数
	database.GetDB().Model(&product).UpdateColumn("like_count", database.GetDB().Raw("like_count + 1"))

	utils.Success(c, gin.H{
		"message": "点赞成功",
		"like_id": like.ID,
	})
}

// UnlikeProduct 取消点赞
func UnlikeProduct(c *gin.Context) {
	userID := utils.GetUserID(c)
	productID := c.Param("id")

	// 删除点赞记录
	result := database.GetDB().Where("product_id = ? AND user_id = ?", productID, userID).Delete(&models.ProductLike{})
	if result.RowsAffected == 0 {
		utils.Fail(c, "未找到点赞记录")
		return
	}

	// 更新商品点赞数
	var product models.Product
	database.GetDB().First(&product, productID)
	database.GetDB().Model(&product).UpdateColumn("like_count", database.GetDB().Raw("GREATEST(like_count - 1, 0)"))

	utils.Success(c, gin.H{
		"message": "取消点赞成功",
	})
}

// CheckLikeStatus 检查点赞状态
func CheckLikeStatus(c *gin.Context) {
	userID := utils.GetUserID(c)
	productID := c.Param("id")

	var like models.ProductLike
	err := database.GetDB().Where("product_id = ? AND user_id = ?", productID, userID).First(&like).Error

	utils.Success(c, gin.H{
		"is_liked": err == nil,
	})
}

// FavoriteProduct 收藏商品
func FavoriteProduct(c *gin.Context) {
	userID := utils.GetUserID(c)
	productID := c.Param("id")

	// 检查商品是否存在
	var product models.Product
	if err := database.GetDB().First(&product, productID).Error; err != nil {
		utils.Fail(c, "商品不存在")
		return
	}

	// 检查是否已收藏
	var favorite models.UserFavorite
	err := database.GetDB().Where("product_id = ? AND user_id = ?", productID, userID).First(&favorite).Error
	if err == nil {
		utils.Fail(c, "已经收藏过该商品")
		return
	}

	// 创建收藏记录
	favorite = models.UserFavorite{
		ProductID: product.ID,
		UserID:    userID,
	}

	if err := database.GetDB().Create(&favorite).Error; err != nil {
		utils.Fail(c, "收藏失败")
		return
	}

	utils.Success(c, gin.H{
		"message":     "收藏成功",
		"favorite_id": favorite.ID,
	})
}

// UnfavoriteProduct 取消收藏
func UnfavoriteProduct(c *gin.Context) {
	userID := utils.GetUserID(c)
	productID := c.Param("id")

	// 删除收藏记录
	result := database.GetDB().Where("product_id = ? AND user_id = ?", productID, userID).Delete(&models.UserFavorite{})
	if result.RowsAffected == 0 {
		utils.Fail(c, "未找到收藏记录")
		return
	}

	utils.Success(c, gin.H{
		"message": "取消收藏成功",
	})
}

// CheckFavoriteStatus 检查收藏状态
func CheckFavoriteStatus(c *gin.Context) {
	userID := utils.GetUserID(c)
	productID := c.Param("id")

	var favorite models.UserFavorite
	err := database.GetDB().Where("product_id = ? AND user_id = ?", productID, userID).First(&favorite).Error

	utils.Success(c, gin.H{
		"is_favorited": err == nil,
	})
}

// GetMyFavorites 获取我的收藏列表
func GetMyFavorites(c *gin.Context) {
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
	database.GetDB().Model(&models.UserFavorite{}).Where("user_id = ?", userID).Count(&total)

	var favorites []models.UserFavorite
	if err := database.GetDB().Preload("Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&favorites).Error; err != nil {
		utils.Fail(c, "获取收藏列表失败")
		return
	}

	utils.PageSuccess(c, favorites, total, page, pageSize)
}
