package handlers

import (
	"shop_api/database"
	"shop_api/models"
	"shop_api/utils"

	"github.com/gin-gonic/gin"
)

// ShippingTemplateInput 运费模板输入
type ShippingTemplateInput struct {
	Name             string                      `json:"name" binding:"required"`
	IsDefault        int8                        `json:"is_default"`
	FreeShippingType int8                        `json:"free_shipping_type" binding:"required,min=1,max=4"`
	FreeAmount       float64                     `json:"free_amount"`
	FreeQuantity     int                         `json:"free_quantity"`
	BaseFee          float64                     `json:"base_fee"`
	Status           int8                        `json:"status"`
	Regions          []ShippingTemplateRegionReq `json:"regions"` // 区域配置
}

// ShippingTemplateRegionReq 区域配置请求
type ShippingTemplateRegionReq struct {
	Province     string  `json:"province" binding:"required"`
	City         string  `json:"city" binding:"required"`
	District     string  `json:"district"`
	Fee          float64 `json:"fee"`
	FreeAmount   float64 `json:"free_amount"`
	FreeQuantity int     `json:"free_quantity"`
}

// GetShippingTemplates 获取运费模板列表
func GetShippingTemplates(c *gin.Context) {
	var templates []models.ShippingTemplate
	if err := database.GetDB().Order("is_default DESC, id ASC").Find(&templates).Error; err != nil {
		utils.Fail(c, "获取运费模板失败")
		return
	}

	// 加载每个模板的区域配置
	for i := range templates {
		var regions []models.ShippingTemplateRegion
		database.GetDB().Where("template_id = ?", templates[i].ID).Find(&regions)
		templates[i].Regions = regions
	}

	utils.Success(c, templates)
}

// GetShippingTemplate 获取单个运费模板详情
func GetShippingTemplate(c *gin.Context) {
	id := c.Param("id")

	var template models.ShippingTemplate
	if err := database.GetDB().First(&template, id).Error; err != nil {
		utils.Fail(c, "运费模板不存在")
		return
	}

	// 加载区域配置
	var regions []models.ShippingTemplateRegion
	database.GetDB().Where("template_id = ?", template.ID).Find(&regions)
	template.Regions = regions

	utils.Success(c, template)
}

// CreateShippingTemplate 创建运费模板
func CreateShippingTemplate(c *gin.Context) {
	var input ShippingTemplateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 如果设置为默认模板，取消其他默认模板
	if input.IsDefault == 1 {
		database.GetDB().Model(&models.ShippingTemplate{}).Where("is_default = ?", 1).Update("is_default", 0)
	}

	tx := database.GetDB().Begin()

	template := models.ShippingTemplate{
		Name:             input.Name,
		IsDefault:        input.IsDefault,
		FreeShippingType: input.FreeShippingType,
		FreeAmount:       input.FreeAmount,
		FreeQuantity:     input.FreeQuantity,
		BaseFee:          input.BaseFee,
		Status:           input.Status,
	}

	if err := tx.Create(&template).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "创建运费模板失败")
		return
	}

	// 创建区域配置
	for _, regionReq := range input.Regions {
		region := models.ShippingTemplateRegion{
			TemplateID:   template.ID,
			Province:     regionReq.Province,
			City:         regionReq.City,
			District:     regionReq.District,
			Fee:          regionReq.Fee,
			FreeAmount:   regionReq.FreeAmount,
			FreeQuantity: regionReq.FreeQuantity,
		}
		if err := tx.Create(&region).Error; err != nil {
			tx.Rollback()
			utils.Fail(c, "创建区域配置失败")
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.Fail(c, "提交事务失败")
		return
	}

	utils.Success(c, template)
}

// UpdateShippingTemplate 更新运费模板
func UpdateShippingTemplate(c *gin.Context) {
	id := c.Param("id")

	var input ShippingTemplateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	var template models.ShippingTemplate
	if err := database.GetDB().First(&template, id).Error; err != nil {
		utils.Fail(c, "运费模板不存在")
		return
	}

	// 如果设置为默认模板，取消其他默认模板
	if input.IsDefault == 1 && template.IsDefault != 1 {
		database.GetDB().Model(&models.ShippingTemplate{}).Where("is_default = ?", 1).Update("is_default", 0)
	}

	tx := database.GetDB().Begin()

	template.Name = input.Name
	template.IsDefault = input.IsDefault
	template.FreeShippingType = input.FreeShippingType
	template.FreeAmount = input.FreeAmount
	template.FreeQuantity = input.FreeQuantity
	template.BaseFee = input.BaseFee
	template.Status = input.Status

	if err := tx.Save(&template).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "更新运费模板失败")
		return
	}

	// 删除旧的区域配置
	if err := tx.Where("template_id = ?", template.ID).Delete(&models.ShippingTemplateRegion{}).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "删除区域配置失败")
		return
	}

	// 创建新的区域配置
	for _, regionReq := range input.Regions {
		region := models.ShippingTemplateRegion{
			TemplateID:   template.ID,
			Province:     regionReq.Province,
			City:         regionReq.City,
			District:     regionReq.District,
			Fee:          regionReq.Fee,
			FreeAmount:   regionReq.FreeAmount,
			FreeQuantity: regionReq.FreeQuantity,
		}
		if err := tx.Create(&region).Error; err != nil {
			tx.Rollback()
			utils.Fail(c, "创建区域配置失败")
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.Fail(c, "提交事务失败")
		return
	}

	utils.Success(c, template)
}

// DeleteShippingTemplate 删除运费模板
func DeleteShippingTemplate(c *gin.Context) {
	id := c.Param("id")

	var template models.ShippingTemplate
	if err := database.GetDB().First(&template, id).Error; err != nil {
		utils.Fail(c, "运费模板不存在")
		return
	}

	// 不允许删除默认模板
	if template.IsDefault == 1 {
		utils.Fail(c, "不能删除默认运费模板")
		return
	}

	tx := database.GetDB().Begin()

	// 删除区域配置
	if err := tx.Where("template_id = ?", id).Delete(&models.ShippingTemplateRegion{}).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "删除区域配置失败")
		return
	}

	// 删除模板
	if err := tx.Delete(&template).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "删除运费模板失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Fail(c, "提交事务失败")
		return
	}

	utils.Success(c, nil)
}

// SetDefaultTemplate 设置默认运费模板
func SetDefaultTemplate(c *gin.Context) {
	id := c.Param("id")

	var template models.ShippingTemplate
	if err := database.GetDB().First(&template, id).Error; err != nil {
		utils.Fail(c, "运费模板不存在")
		return
	}

	tx := database.GetDB().Begin()

	// 取消其他默认模板
	if err := tx.Model(&models.ShippingTemplate{}).Where("is_default = ?", 1).Update("is_default", 0).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "设置默认模板失败")
		return
	}

	// 设置当前模板为默认
	if err := tx.Model(&template).Update("is_default", 1).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "设置默认模板失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Fail(c, "提交事务失败")
		return
	}

	utils.Success(c, nil)
}

// CalculateShippingFeeRequest 计算运费请求
type CalculateShippingFeeRequest struct {
	Province string  `json:"province" binding:"required"`
	City     string  `json:"city" binding:"required"`
	District string  `json:"district"`
	Amount   float64 `json:"amount" binding:"required,min=0"`
	Quantity int     `json:"quantity" binding:"required,min=1"`
}

// CalculateShippingFeeResponse 计算运费响应
type CalculateShippingFeeResponse struct {
	FreightAmount float64 `json:"freight_amount"` // 运费金额
	IsFree        bool    `json:"is_free"`        // 是否包邮
	Message       string  `json:"message"`        // 提示信息
}

// CalculateShippingFee 计算运费（小程序端调用）
func CalculateShippingFee(c *gin.Context) {
	var req CalculateShippingFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	freightAmount, isFree, message := calculateFreight(req.Province, req.City, req.District, req.Amount, req.Quantity)

	utils.Success(c, CalculateShippingFeeResponse{
		FreightAmount: freightAmount,
		IsFree:        isFree,
		Message:       message,
	})
}

// calculateFreight 计算运费的核心逻辑
func calculateFreight(province, city, district string, amount float64, quantity int) (float64, bool, string) {
	// 获取默认运费模板
	var template models.ShippingTemplate
	if err := database.GetDB().Where("is_default = ? AND status = ?", 1, 1).First(&template).Error; err != nil {
		// 如果没有默认模板，返回0运费
		return 0, true, "暂无运费配置"
	}

	// 检查是否满足全局包邮条件
	isFreeByGlobal := checkGlobalFreeShipping(template, amount, quantity)

	// 查找该地区的特殊配置
	var region models.ShippingTemplateRegion
	err := database.GetDB().Where("template_id = ? AND province = ? AND city = ?", template.ID, province, city).
		First(&region).Error

	if err == nil {
		// 找到地区特殊配置
		// 检查该地区是否满足包邮条件
		isFreeByRegion := checkRegionFreeShipping(region, amount, quantity)

		if isFreeByRegion || isFreeByGlobal {
			return 0, true, "满足包邮条件"
		}

		return region.Fee, false, ""
	}

	// 没有特殊配置，使用基础运费
	if isFreeByGlobal {
		return 0, true, "满足包邮条件"
	}

	return template.BaseFee, false, ""
}

// checkGlobalFreeShipping 检查是否满足全局包邮条件
func checkGlobalFreeShipping(template models.ShippingTemplate, amount float64, quantity int) bool {
	switch template.FreeShippingType {
	case models.FreeShippingTypeAmount:
		// 满额包邮
		return template.FreeAmount > 0 && amount >= template.FreeAmount
	case models.FreeShippingTypeQuantity:
		// 满件包邮
		return template.FreeQuantity > 0 && quantity >= template.FreeQuantity
	case models.FreeShippingTypeEither:
		// 满额或满件
		return (template.FreeAmount > 0 && amount >= template.FreeAmount) ||
			(template.FreeQuantity > 0 && quantity >= template.FreeQuantity)
	case models.FreeShippingTypeNone:
		// 不包邮
		return false
	default:
		return false
	}
}

// checkRegionFreeShipping 检查是否满足区域包邮条件
func checkRegionFreeShipping(region models.ShippingTemplateRegion, amount float64, quantity int) bool {
	// 优先使用区域的包邮条件，如果区域未设置则使用0（表示不包邮）
	if region.FreeAmount > 0 && amount >= region.FreeAmount {
		return true
	}
	if region.FreeQuantity > 0 && quantity >= region.FreeQuantity {
		return true
	}
	return false
}
