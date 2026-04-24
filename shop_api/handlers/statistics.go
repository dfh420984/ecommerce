package handlers

import (
	"shop_api/database"
	"shop_api/models"
	"shop_api/services"
	"shop_api/utils"

	"github.com/gin-gonic/gin"
)

// GetDashboardStats 获取仪表盘统计数据
func GetDashboardStats(c *gin.Context) {
	stats, err := services.GetStatisticsService().GetDashboardStats()
	if err != nil {
		utils.Fail(c, "获取统计数据失败: "+err.Error())
		return
	}

	utils.Success(c, stats)
}

// GetSalesTrend 获取销售趋势
func GetSalesTrend(c *gin.Context) {
	trend, err := services.GetStatisticsService().GetSalesTrend()
	if err != nil {
		utils.Fail(c, "获取销售趋势失败: "+err.Error())
		return
	}

	utils.Success(c, trend)
}

// GetUsersTrend 获取用户增长趋势
func GetUsersTrend(c *gin.Context) {
	trend, err := services.GetStatisticsService().GetUsersTrend()
	if err != nil {
		utils.Fail(c, "获取用户趋势失败: "+err.Error())
		return
	}

	utils.Success(c, trend)
}

// GetOrderLogistics 查询订单物流
func GetOrderLogistics(c *gin.Context) {
	userID := utils.GetUserID(c)
	orderID := c.Param("id")

	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		utils.Fail(c, "订单不存在")
		return
	}

	if order.ExpressCompany == "" || order.ExpressNo == "" {
		utils.Fail(c, "暂无物流信息")
		return
	}

	tracks, err := services.GetLogisticsService().QueryTrack(order.ExpressCompany, order.ExpressNo)
	if err != nil {
		utils.Fail(c, "查询物流失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"express_company": order.ExpressCompany,
		"express_no":      order.ExpressNo,
		"tracks":          tracks,
	})
}
