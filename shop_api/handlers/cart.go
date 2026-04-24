package handlers

import (
	"shop_api/database"
	"shop_api/models"
	"shop_api/utils"

	"github.com/gin-gonic/gin"
)

type CartInput struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	SkuID     uint64 `json:"sku_id"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

func GetCart(c *gin.Context) {
	userID := utils.GetUserID(c)

	var carts []models.Cart
	if err := database.GetDB().Preload("Product").Where("user_id = ?", userID).Find(&carts).Error; err != nil {
		utils.Fail(c, "获取购物车失败")
		return
	}

	var total float64
	for _, cart := range carts {
		if cart.Product != nil && cart.Selected == 1 {
			total += cart.Product.Price * float64(cart.Quantity)
		}
	}

	// 保留2位小数，避免浮点数精度问题
	total = float64(int(total*100+0.5)) / 100

	utils.Success(c, gin.H{
		"list":  carts,
		"total": total,
		"count": len(carts),
	})
}

func AddCart(c *gin.Context) {
	userID := utils.GetUserID(c)

	var input CartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	var product models.Product
	if err := database.GetDB().First(&product, input.ProductID).Error; err != nil {
		utils.Fail(c, "商品不存在")
		return
	}

	if product.Stock < input.Quantity {
		utils.Fail(c, "库存不足")
		return
	}

	var cart models.Cart
	err := database.GetDB().Where("user_id = ? AND product_id = ? AND sku_id = ?", userID, input.ProductID, input.SkuID).First(&cart).Error

	if err == nil {
		cart.Quantity += input.Quantity
		if cart.Quantity > product.Stock {
			cart.Quantity = product.Stock
		}
		database.GetDB().Save(&cart)
	} else {
		cart = models.Cart{
			UserID:    userID,
			ProductID: input.ProductID,
			SkuID:     input.SkuID,
			Quantity:  input.Quantity,
			Selected:  1,
		}
		database.GetDB().Create(&cart)
	}

	utils.Success(c, cart)
}

func UpdateCart(c *gin.Context) {
	id := c.Param("id")
	userID := utils.GetUserID(c)

	var input struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	var cart models.Cart
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&cart).Error; err != nil {
		utils.Fail(c, "购物车商品不存在")
		return
	}

	var product models.Product
	if err := database.GetDB().First(&product, cart.ProductID).Error; err != nil {
		utils.Fail(c, "商品不存在")
		return
	}

	if input.Quantity > product.Stock {
		utils.Fail(c, "库存不足")
		return
	}

	cart.Quantity = input.Quantity
	database.GetDB().Save(&cart)

	utils.Success(c, cart)
}

func SelectCart(c *gin.Context) {
	id := c.Param("id")
	userID := utils.GetUserID(c)

	var input struct {
		Selected int8 `json:"selected"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	database.GetDB().Model(&models.Cart{}).Where("id = ? AND user_id = ?", id, userID).Update("selected", input.Selected)

	utils.Success(c, nil)
}

func SelectAllCart(c *gin.Context) {
	userID := utils.GetUserID(c)

	var input struct {
		Selected int8 `json:"selected"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	database.GetDB().Model(&models.Cart{}).Where("user_id = ?", userID).Update("selected", input.Selected)

	utils.Success(c, nil)
}

func DeleteCart(c *gin.Context) {
	id := c.Param("id")
	userID := utils.GetUserID(c)

	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).Delete(&models.Cart{}).Error; err != nil {
		utils.Fail(c, "删除失败")
		return
	}

	utils.Success(c, nil)
}

func ClearCart(c *gin.Context) {
	userID := utils.GetUserID(c)

	database.GetDB().Where("user_id = ?", userID).Delete(&models.Cart{})

	utils.Success(c, nil)
}

func GetCartCount(c *gin.Context) {
	userID := utils.GetUserID(c)

	var count int64
	database.GetDB().Model(&models.Cart{}).Where("user_id = ?", userID).Count(&count)

	utils.Success(c, gin.H{"count": count})
}
