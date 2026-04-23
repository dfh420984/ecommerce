package handlers

import (
	"shop_api/database"
	"shop_api/models"
	"shop_api/utils"

	"github.com/gin-gonic/gin"
)

type ConfigCreateInput struct {
	Name        string `json:"name" binding:"required"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

// GetConfigs 获取系统配置列表（管理员）
func GetConfigs(c *gin.Context) {
	var configs []models.Config
	if err := database.GetDB().Order("id ASC").Find(&configs).Error; err != nil {
		utils.Fail(c, "获取配置列表失败")
		return
	}

	utils.Success(c, configs)
}

// GetConfigByName 根据名称获取配置值（公开接口，小程序使用）
func GetConfigByName(c *gin.Context) {
	name := c.Param("name")

	var config models.Config
	if err := database.GetDB().Where("name = ?", name).First(&config).Error; err != nil {
		utils.Fail(c, "配置不存在")
		return
	}

	utils.Success(c, config)
}

// GetConfig 获取单个配置详情
func GetConfig(c *gin.Context) {
	id := c.Param("id")

	var config models.Config
	if err := database.GetDB().First(&config, id).Error; err != nil {
		utils.Fail(c, "配置不存在")
		return
	}

	utils.Success(c, config)
}

// CreateConfig 创建系统配置
func CreateConfig(c *gin.Context) {
	var input ConfigCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	var existingConfig models.Config
	if err := database.GetDB().Where("name = ?", input.Name).First(&existingConfig).Error; err == nil {
		utils.Fail(c, "配置名称已存在")
		return
	}

	config := models.Config{
		Name:        input.Name,
		Value:       input.Value,
		Description: input.Description,
	}

	if err := database.GetDB().Create(&config).Error; err != nil {
		utils.Fail(c, "创建配置失败")
		return
	}

	utils.Success(c, config)
}

// UpdateConfig 更新系统配置
func UpdateConfig(c *gin.Context) {
	id := c.Param("id")

	var input ConfigCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	var config models.Config
	if err := database.GetDB().First(&config, id).Error; err != nil {
		utils.Fail(c, "配置不存在")
		return
	}

	config.Name = input.Name
	config.Value = input.Value
	config.Description = input.Description

	if err := database.GetDB().Save(&config).Error; err != nil {
		utils.Fail(c, "更新配置失败")
		return
	}

	utils.Success(c, config)
}

// DeleteConfig 删除系统配置
func DeleteConfig(c *gin.Context) {
	id := c.Param("id")

	if err := database.GetDB().Delete(&models.Config{}, id).Error; err != nil {
		utils.Fail(c, "删除配置失败")
		return
	}

	utils.Success(c, nil)
}

// GetConfigsByNames 批量获取多个配置值（公开接口，小程序使用）
func GetConfigsByNames(c *gin.Context) {
	type Request struct {
		Names []string `json:"names" binding:"required"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	var configs []models.Config
	if err := database.GetDB().Where("name IN ?", req.Names).Find(&configs).Error; err != nil {
		utils.Fail(c, "获取配置失败")
		return
	}

	result := make(map[string]string)
	for _, config := range configs {
		result[config.Name] = config.Value
	}

	utils.Success(c, result)
}
