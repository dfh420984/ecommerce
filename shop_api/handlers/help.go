package handlers

import (
	"shop_api/database"
	"shop_api/models"
	"shop_api/utils"

	"github.com/gin-gonic/gin"
)

// ============ 分类管理 ============

type CategoryInput struct {
	Name        string `json:"name" binding:"required"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}

// GetHelpCategories 获取帮助中心分类列表（管理员）
func GetHelpCategories(c *gin.Context) {
	var categories []models.HelpCategory
	if err := database.GetDB().Order("sort ASC, id ASC").Find(&categories).Error; err != nil {
		utils.Fail(c, "获取分类列表失败")
		return
	}

	utils.Success(c, categories)
}

// GetActiveHelpCategories 获取启用的分类列表（小程序使用）
func GetActiveHelpCategories(c *gin.Context) {
	var categories []models.HelpCategory
	if err := database.GetDB().Where("status = ?", 1).Order("sort ASC, id ASC").Find(&categories).Error; err != nil {
		utils.Fail(c, "获取分类列表失败")
		return
	}

	utils.Success(c, categories)
}

// CreateHelpCategory 创建分类
func CreateHelpCategory(c *gin.Context) {
	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	category := models.HelpCategory{
		Name:        input.Name,
		Sort:        input.Sort,
		Status:      input.Status,
		Description: input.Description,
	}

	if err := database.GetDB().Create(&category).Error; err != nil {
		utils.Fail(c, "创建分类失败")
		return
	}

	utils.Success(c, category)
}

// UpdateHelpCategory 更新分类
func UpdateHelpCategory(c *gin.Context) {
	id := c.Param("id")

	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	var category models.HelpCategory
	if err := database.GetDB().First(&category, id).Error; err != nil {
		utils.Fail(c, "分类不存在")
		return
	}

	category.Name = input.Name
	category.Sort = input.Sort
	category.Status = input.Status
	category.Description = input.Description

	if err := database.GetDB().Save(&category).Error; err != nil {
		utils.Fail(c, "更新分类失败")
		return
	}

	utils.Success(c, category)
}

// DeleteHelpCategory 删除分类
func DeleteHelpCategory(c *gin.Context) {
	id := c.Param("id")

	// 检查是否有问题使用该分类
	var count int64
	database.GetDB().Model(&models.HelpQuestion{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		utils.Fail(c, "该分类下还有问题，无法删除")
		return
	}

	if err := database.GetDB().Delete(&models.HelpCategory{}, id).Error; err != nil {
		utils.Fail(c, "删除分类失败")
		return
	}

	utils.Success(c, nil)
}

// ============ 问题管理 ============

type QuestionInput struct {
	CategoryID uint64 `json:"category_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
	Sort       int    `json:"sort"`
	Status     int    `json:"status"`
}

// GetHelpQuestions 获取问题列表（管理员）
func GetHelpQuestions(c *gin.Context) {
	categoryID := c.Query("category_id")

	query := database.GetDB().Preload("Category")

	if categoryID != "" && categoryID != "0" {
		query = query.Where("category_id = ?", categoryID)
	}

	var questions []models.HelpQuestion
	if err := query.Order("sort ASC, id ASC").Find(&questions).Error; err != nil {
		utils.Fail(c, "获取问题列表失败")
		return
	}

	utils.Success(c, questions)
}

// GetHelpQuestionsByCategory 根据分类获取问题列表（小程序使用）
func GetHelpQuestionsByCategory(c *gin.Context) {
	categoryID := c.Param("category_id")

	query := database.GetDB().Where("status = ?", 1)

	// 只有 categoryID > 0 时才加入分类过滤条件
	if categoryID != "" && categoryID != "0" {
		query = query.Where("category_id = ?", categoryID)
	}

	var questions []models.HelpQuestion
	if err := query.Order("sort ASC, id ASC").Find(&questions).Error; err != nil {
		utils.Fail(c, "获取问题列表失败")
		return
	}

	utils.Success(c, questions)
}

// GetHelpQuestionDetail 获取问题详情
func GetHelpQuestionDetail(c *gin.Context) {
	id := c.Param("id")

	var question models.HelpQuestion
	if err := database.GetDB().Preload("Category").First(&question, id).Error; err != nil {
		utils.Fail(c, "问题不存在")
		return
	}

	utils.Success(c, question)
}

// CreateHelpQuestion 创建问题
func CreateHelpQuestion(c *gin.Context) {
	var input QuestionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 检查分类是否存在
	var category models.HelpCategory
	if err := database.GetDB().First(&category, input.CategoryID).Error; err != nil {
		utils.Fail(c, "分类不存在")
		return
	}

	question := models.HelpQuestion{
		CategoryID: input.CategoryID,
		Title:      input.Title,
		Answer:     input.Answer,
		Sort:       input.Sort,
		Status:     input.Status,
	}

	if err := database.GetDB().Create(&question).Error; err != nil {
		utils.Fail(c, "创建问题失败")
		return
	}

	// 重新加载关联数据
	database.GetDB().Preload("Category").First(&question, question.ID)

	utils.Success(c, question)
}

// UpdateHelpQuestion 更新问题
func UpdateHelpQuestion(c *gin.Context) {
	id := c.Param("id")

	var input QuestionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	// 检查分类是否存在
	var category models.HelpCategory
	if err := database.GetDB().First(&category, input.CategoryID).Error; err != nil {
		utils.Fail(c, "分类不存在")
		return
	}

	var question models.HelpQuestion
	if err := database.GetDB().First(&question, id).Error; err != nil {
		utils.Fail(c, "问题不存在")
		return
	}

	question.CategoryID = input.CategoryID
	question.Title = input.Title
	question.Answer = input.Answer
	question.Sort = input.Sort
	question.Status = input.Status

	if err := database.GetDB().Save(&question).Error; err != nil {
		utils.Fail(c, "更新问题失败")
		return
	}

	// 重新加载关联数据
	database.GetDB().Preload("Category").First(&question, question.ID)

	utils.Success(c, question)
}

// DeleteHelpQuestion 删除问题
func DeleteHelpQuestion(c *gin.Context) {
	id := c.Param("id")

	if err := database.GetDB().Delete(&models.HelpQuestion{}, id).Error; err != nil {
		utils.Fail(c, "删除问题失败")
		return
	}

	utils.Success(c, nil)
}

// SearchHelpQuestions 搜索问题（小程序使用）
func SearchHelpQuestions(c *gin.Context) {
	keyword := c.Query("keyword")
	categoryID := c.Query("category_id")

	if keyword == "" {
		utils.Fail(c, "请输入搜索关键词")
		return
	}

	query := database.GetDB().Preload("Category").
		Where("status = ?", 1).
		Where("title LIKE ? OR answer LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	if categoryID != "" && categoryID != "0" {
		query = query.Where("category_id = ?", categoryID)
	}

	var questions []models.HelpQuestion
	if err := query.Order("sort ASC, id ASC").Find(&questions).Error; err != nil {
		utils.Fail(c, "搜索失败")
		return
	}

	utils.Success(c, questions)
}

// GetHelpConfig 获取帮助中心配置（服务时间等）
func GetHelpConfig(c *gin.Context) {
	configNames := []string{"service_time", "customer_service_phone", "online_service_url"}

	var configs []models.Config
	if err := database.GetDB().Where("name IN ?", configNames).Find(&configs).Error; err != nil {
		utils.Fail(c, "获取配置失败")
		return
	}

	result := make(map[string]string)
	for _, config := range configs {
		result[config.Name] = config.Value
	}

	// 设置默认值
	if _, ok := result["service_time"]; !ok {
		result["service_time"] = "09:00-21:00"
	}
	if _, ok := result["customer_service_phone"]; !ok {
		result["customer_service_phone"] = "400-123-4567"
	}

	utils.Success(c, result)
}

// BatchUpdateQuestionsStatus 批量更新问题状态
func BatchUpdateQuestionsStatus(c *gin.Context) {
	type Request struct {
		IDs    []uint64 `json:"ids" binding:"required"`
		Status int      `json:"status" binding:"required"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	if err := database.GetDB().Model(&models.HelpQuestion{}).
		Where("id IN ?", req.IDs).
		Update("status", req.Status).Error; err != nil {
		utils.Fail(c, "更新失败")
		return
	}

	utils.Success(c, nil)
}

// GetHelpStatistics 获取帮助中心统计信息
func GetHelpStatistics(c *gin.Context) {
	var categoryCount, questionCount, activeQuestionCount int64

	database.GetDB().Model(&models.HelpCategory{}).Count(&categoryCount)
	database.GetDB().Model(&models.HelpQuestion{}).Count(&questionCount)
	database.GetDB().Model(&models.HelpQuestion{}).Where("status = ?", 1).Count(&activeQuestionCount)

	stats := map[string]interface{}{
		"category_count":        categoryCount,
		"question_count":        questionCount,
		"active_question_count": activeQuestionCount,
	}

	utils.Success(c, stats)
}
