package handlers

import (
	"shop_api/database"
	"shop_api/models"
	"shop_api/utils"

	"github.com/gin-gonic/gin"
)

type AddressInput struct {
	Consignee  string `json:"consignee" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	Province   string `json:"province" binding:"required"`
	City       string `json:"city" binding:"required"`
	District   string `json:"district" binding:"required"`
	Address    string `json:"address" binding:"required"`
	PostalCode string `json:"postal_code"`
	IsDefault  int8   `json:"is_default"`
}

func GetAddresses(c *gin.Context) {
	userID := utils.GetUserID(c)

	var addresses []models.Address
	if err := database.GetDB().Where("user_id = ?", userID).Order("is_default DESC, id DESC").Find(&addresses).Error; err != nil {
		utils.Fail(c, "获取地址失败")
		return
	}

	utils.Success(c, addresses)
}

func GetAddress(c *gin.Context) {
	id := c.Param("id")
	userID := utils.GetUserID(c)

	var address models.Address
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil {
		utils.Fail(c, "地址不存在")
		return
	}

	utils.Success(c, address)
}

func CreateAddress(c *gin.Context) {
	userID := utils.GetUserID(c)

	var input AddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	if input.IsDefault == 1 {
		database.GetDB().Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", 0)
	}

	address := models.Address{
		UserID:     userID,
		Consignee:  input.Consignee,
		Phone:      input.Phone,
		Province:   input.Province,
		City:       input.City,
		District:   input.District,
		Address:    input.Address,
		PostalCode: input.PostalCode,
		IsDefault:  input.IsDefault,
	}

	if err := database.GetDB().Create(&address).Error; err != nil {
		utils.Fail(c, "创建失败")
		return
	}

	utils.Success(c, address)
}

func UpdateAddress(c *gin.Context) {
	id := c.Param("id")
	userID := utils.GetUserID(c)

	var input AddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	var address models.Address
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil {
		utils.Fail(c, "地址不存在")
		return
	}

	if input.IsDefault == 1 {
		database.GetDB().Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", 0)
	}

	address.Consignee = input.Consignee
	address.Phone = input.Phone
	address.Province = input.Province
	address.City = input.City
	address.District = input.District
	address.Address = input.Address
	address.PostalCode = input.PostalCode
	address.IsDefault = input.IsDefault

	database.GetDB().Save(&address)

	utils.Success(c, address)
}

func DeleteAddress(c *gin.Context) {
	id := c.Param("id")
	userID := utils.GetUserID(c)

	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).Delete(&models.Address{}).Error; err != nil {
		utils.Fail(c, "删除失败")
		return
	}

	utils.Success(c, nil)
}

func SetDefaultAddress(c *gin.Context) {
	id := c.Param("id")
	userID := utils.GetUserID(c)

	database.GetDB().Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", 0)
	database.GetDB().Model(&models.Address{}).Where("id = ? AND user_id = ?", id, userID).Update("is_default", 1)

	utils.Success(c, nil)
}
