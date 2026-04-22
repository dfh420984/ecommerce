package handlers

import (
	"github.com/gin-gonic/gin"
	"shop_api/database"
	"shop_api/models"
	"shop_api/utils"
)

type CategoryCreateInput struct {
	Name     string `json:"name" binding:"required"`
	ParentID uint64 `json:"parent_id"`
	Sort     int    `json:"sort"`
	Icon     string `json:"icon"`
	Image    string `json:"image"`
	Status   int8   `json:"status"`
}

func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := database.GetDB().Where("status = ?", 1).Order("sort ASC, id ASC").Find(&categories).Error; err != nil {
		utils.Fail(c, "获取分类失败")
		return
	}

	tree := buildCategoryTree(categories, 0)
	utils.Success(c, tree)
}

func buildCategoryTree(categories []models.Category, parentID uint64) []models.Category {
	var result []models.Category
	for _, cat := range categories {
		if cat.ParentID == parentID {
			cat.Children = buildCategoryTree(categories, cat.ID)
			result = append(result, cat)
		}
	}
	return result
}

func GetCategory(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	if err := database.GetDB().First(&category, id).Error; err != nil {
		utils.Fail(c, "分类不存在")
		return
	}

	utils.Success(c, category)
}

func CreateCategory(c *gin.Context) {
	var input CategoryCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	level := 1
	if input.ParentID > 0 {
		var parent models.Category
		if err := database.GetDB().First(&parent, input.ParentID).Error; err != nil {
			utils.Fail(c, "父分类不存在")
			return
		}
		level = parent.Level + 1
	}

	category := models.Category{
		Name:     input.Name,
		ParentID: input.ParentID,
		Level:    level,
		Sort:     input.Sort,
		Icon:     input.Icon,
		Image:    input.Image,
		Status:   input.Status,
	}

	if category.Status == 0 {
		category.Status = 1
	}

	if err := database.GetDB().Create(&category).Error; err != nil {
		utils.Fail(c, "创建失败")
		return
	}

	utils.Success(c, category)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var input CategoryCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	var category models.Category
	if err := database.GetDB().First(&category, id).Error; err != nil {
		utils.Fail(c, "分类不存在")
		return
	}

	level := 1
	if input.ParentID > 0 {
		var parent models.Category
		if err := database.GetDB().First(&parent, input.ParentID).Error; err != nil {
			utils.Fail(c, "父分类不存在")
			return
		}
		level = parent.Level + 1
	}

	category.Name = input.Name
	category.ParentID = input.ParentID
	category.Level = level
	category.Sort = input.Sort
	category.Icon = input.Icon
	category.Image = input.Image
	category.Status = input.Status

	database.GetDB().Save(&category)

	utils.Success(c, category)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	var count int64
	database.GetDB().Model(&models.Category{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		utils.Fail(c, "该分类下有子分类，无法删除")
		return
	}

	if err := database.GetDB().Delete(&models.Category{}, id).Error; err != nil {
		utils.Fail(c, "删除失败")
		return
	}

	utils.Success(c, nil)
}

type ProductListInput struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	CategoryID uint64 `form:"category_id"`
	Keyword    string `form:"keyword"`
	IsOnline   int    `form:"is_online"`
}

func GetProducts(c *gin.Context) {
	var input ProductListInput
	if err := c.ShouldBindQuery(&input); err != nil {
		input.Page = 1
		input.PageSize = 10
	}

	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 10
	}

	query := database.GetDB().Model(&models.Product{}).Where("status = ?", 1)

	if input.CategoryID > 0 {
		query = query.Where("category_id = ?", input.CategoryID)
	}
	if input.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+input.Keyword+"%")
	}
	if input.IsOnline > 0 {
		query = query.Where("is_online = ?", input.IsOnline)
	}

	var total int64
	query.Count(&total)

	var products []models.Product
	if err := query.Order("sort ASC, id DESC").
		Offset((input.Page - 1) * input.PageSize).
		Limit(input.PageSize).
		Find(&products).Error; err != nil {
		utils.Fail(c, "获取商品失败")
		return
	}

	utils.PageSuccess(c, products, total, input.Page, input.PageSize)
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := database.GetDB().Where("status = ?", 1).First(&product, id).Error; err != nil {
		utils.Fail(c, "商品不存在")
		return
	}

	utils.Success(c, product)
}

func GetRecommendProducts(c *gin.Context) {
	var products []models.Product
	if err := database.GetDB().Where("status = ? AND is_online = ? AND is_recommend = ?", 1, 1, 1).
		Order("sort ASC, id DESC").
		Limit(10).
		Find(&products).Error; err != nil {
		utils.Fail(c, "获取推荐商品失败")
		return
	}

	utils.Success(c, products)
}

func GetNewProducts(c *gin.Context) {
	var products []models.Product
	if err := database.GetDB().Where("status = ? AND is_online = ? AND is_new = ?", 1, 1, 1).
		Order("id DESC").
		Limit(10).
		Find(&products).Error; err != nil {
		utils.Fail(c, "获取新品失败")
		return
	}

	utils.Success(c, products)
}

type ProductCreateInput struct {
	Name          string   `json:"name" binding:"required"`
	CategoryID    uint64   `json:"category_id" binding:"required"`
	Price         float64  `json:"price" binding:"required"`
	OriginalPrice float64  `json:"original_price"`
	Cost          float64  `json:"cost"`
	Stock         int      `json:"stock"`
	Images        []string `json:"images"`
	Description   string   `json:"description"`
	Content       string   `json:"content"`
	IsOnline      int8     `json:"is_online"`
	IsRecommend   int8     `json:"is_recommend"`
	IsNew         int8     `json:"is_new"`
	Sort          int      `json:"sort"`
}

func CreateProduct(c *gin.Context) {
	var input ProductCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	product := models.Product{
		Name:          input.Name,
		CategoryID:    input.CategoryID,
		Price:         input.Price,
		OriginalPrice: input.OriginalPrice,
		Cost:          input.Cost,
		Stock:         input.Stock,
		Images:        input.Images,
		Description:   input.Description,
		Content:       input.Content,
		IsOnline:      input.IsOnline,
		IsRecommend:   input.IsRecommend,
		IsNew:         input.IsNew,
		Sort:          input.Sort,
		Status:        1,
	}

	if err := database.GetDB().Create(&product).Error; err != nil {
		utils.Fail(c, "创建失败")
		return
	}

	utils.Success(c, product)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var input ProductCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	var product models.Product
	if err := database.GetDB().First(&product, id).Error; err != nil {
		utils.Fail(c, "商品不存在")
		return
	}

	product.Name = input.Name
	product.CategoryID = input.CategoryID
	product.Price = input.Price
	product.OriginalPrice = input.OriginalPrice
	product.Cost = input.Cost
	product.Stock = input.Stock
	product.Images = input.Images
	product.Description = input.Description
	product.Content = input.Content
	product.IsOnline = input.IsOnline
	product.IsRecommend = input.IsRecommend
	product.IsNew = input.IsNew
	product.Sort = input.Sort

	database.GetDB().Save(&product)

	utils.Success(c, product)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := database.GetDB().Model(&models.Product{}).Where("id = ?", id).Update("status", 0).Error; err != nil {
		utils.Fail(c, "删除失败")
		return
	}

	utils.Success(c, nil)
}

func GetBanners(c *gin.Context) {
	var banners []models.Banner
	if err := database.GetDB().Where("status = ?", 1).Order("sort ASC, id DESC").Find(&banners).Error; err != nil {
		utils.Fail(c, "获取轮播图失败")
		return
	}

	utils.Success(c, banners)
}

type BannerCreateInput struct {
	Title    string `json:"title"`
	Image    string `json:"image" binding:"required"`
	Link     string `json:"link"`
	LinkType int8   `json:"link_type"`
	TargetID uint64 `json:"target_id"`
	Sort     int    `json:"sort"`
	Status   int8   `json:"status"`
}

func CreateBanner(c *gin.Context) {
	var input BannerCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	banner := models.Banner{
		Title:    input.Title,
		Image:    input.Image,
		Link:     input.Link,
		LinkType: input.LinkType,
		TargetID: input.TargetID,
		Sort:     input.Sort,
		Status:   input.Status,
	}

	if banner.Status == 0 {
		banner.Status = 1
	}

	if err := database.GetDB().Create(&banner).Error; err != nil {
		utils.Fail(c, "创建失败")
		return
	}

	utils.Success(c, banner)
}

func UpdateBanner(c *gin.Context) {
	id := c.Param("id")

	var input BannerCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	var banner models.Banner
	if err := database.GetDB().First(&banner, id).Error; err != nil {
		utils.Fail(c, "轮播图不存在")
		return
	}

	banner.Title = input.Title
	banner.Image = input.Image
	banner.Link = input.Link
	banner.LinkType = input.LinkType
	banner.TargetID = input.TargetID
	banner.Sort = input.Sort
	banner.Status = input.Status

	database.GetDB().Save(&banner)

	utils.Success(c, banner)
}

func DeleteBanner(c *gin.Context) {
	id := c.Param("id")

	if err := database.GetDB().Delete(&models.Banner{}, id).Error; err != nil {
		utils.Fail(c, "删除失败")
		return
	}

	utils.Success(c, nil)
}
