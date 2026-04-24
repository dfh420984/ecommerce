package services

import (
	"shop_api/database"
	"shop_api/models"
	"time"
)

// StatisticsService 统计服务
type StatisticsService struct{}

var statisticsService *StatisticsService

func GetStatisticsService() *StatisticsService {
	if statisticsService == nil {
		statisticsService = &StatisticsService{}
	}
	return statisticsService
}

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	TodaySales     float64          `json:"today_sales"`     // 今日销售额
	YesterdaySales float64          `json:"yesterday_sales"` // 昨日销售额
	MonthSales     float64          `json:"month_sales"`     // 本月销售额
	TodayUsers     int64            `json:"today_users"`     // 今日新增用户
	TodayOrders    int64            `json:"today_orders"`    // 今日订单数
	PendingOrders  int64            `json:"pending_orders"`  // 待处理订单
	HotProducts    []models.Product `json:"hot_products"`    // 热销商品
}

// SalesTrendItem 销售趋势项
type SalesTrendItem struct {
	Date   string  `json:"date"`
	Sales  float64 `json:"sales"`
	Orders int64   `json:"orders"`
}

// UsersTrendItem 用户趋势项
type UsersTrendItem struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// GetDashboardStats 获取仪表盘统计数据
func (s *StatisticsService) GetDashboardStats() (*DashboardStats, error) {
	db := database.GetDB()
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	firstDayOfMonth := time.Now().Format("2006-01-01")

	stats := &DashboardStats{}

	// 今日销售额
	db.Model(&models.Order{}).
		Where("pay_status = ? AND DATE(pay_time) = ?", models.PayStatusPaid, today).
		Select("COALESCE(SUM(pay_amount), 0)").
		Scan(&stats.TodaySales)

	// 昨日销售额
	db.Model(&models.Order{}).
		Where("pay_status = ? AND DATE(pay_time) = ?", models.PayStatusPaid, yesterday).
		Select("COALESCE(SUM(pay_amount), 0)").
		Scan(&stats.YesterdaySales)

	// 本月销售额
	db.Model(&models.Order{}).
		Where("pay_status = ? AND DATE(pay_time) >= ?", models.PayStatusPaid, firstDayOfMonth).
		Select("COALESCE(SUM(pay_amount), 0)").
		Scan(&stats.MonthSales)

	// 今日新增用户
	db.Model(&models.User{}).
		Where("DATE(created_at) = ?", today).
		Count(&stats.TodayUsers)

	// 今日订单数
	db.Model(&models.Order{}).
		Where("DATE(created_at) = ?", today).
		Count(&stats.TodayOrders)

	// 待处理订单（待支付+已支付待发货）
	db.Model(&models.Order{}).
		Where("order_status IN (?)", []int8{
			models.OrderStatusPending,
			models.OrderStatusPaid,
		}).
		Count(&stats.PendingOrders)

	// 热销商品TOP10
	db.Order("sales DESC").Limit(10).Find(&stats.HotProducts)

	return stats, nil
}

// GetSalesTrend 获取销售趋势（最近7天）
func (s *StatisticsService) GetSalesTrend() ([]SalesTrendItem, error) {
	db := database.GetDB()
	var trend []SalesTrendItem

	db.Raw(`
		SELECT 
			DATE(pay_time) as date,
			COALESCE(SUM(pay_amount), 0) as sales,
			COUNT(*) as orders
		FROM orders
		WHERE pay_status = ? 
		  AND pay_time >= DATE_SUB(NOW(), INTERVAL 7 DAY)
		GROUP BY DATE(pay_time)
		ORDER BY date ASC
	`, models.PayStatusPaid).Scan(&trend)

	// 如果某些日期没有数据，补充0
	return s.fillMissingDates(trend, 7), nil
}

// GetUsersTrend 获取用户增长趋势（最近7天）
func (s *StatisticsService) GetUsersTrend() ([]UsersTrendItem, error) {
	db := database.GetDB()
	var trend []UsersTrendItem

	db.Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM users
		WHERE created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`).Scan(&trend)

	return s.fillMissingUserDates(trend, 7), nil
}

// fillMissingDates 填充缺失的日期数据
func (s *StatisticsService) fillMissingDates(trend []SalesTrendItem, days int) []SalesTrendItem {
	result := make([]SalesTrendItem, 0, days)
	dateMap := make(map[string]*SalesTrendItem)

	for i := range trend {
		dateMap[trend[i].Date] = &trend[i]
	}

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		if item, ok := dateMap[date]; ok {
			result = append(result, *item)
		} else {
			result = append(result, SalesTrendItem{
				Date:   date,
				Sales:  0,
				Orders: 0,
			})
		}
	}

	return result
}

// fillMissingUserDates 填充缺失的用户日期数据
func (s *StatisticsService) fillMissingUserDates(trend []UsersTrendItem, days int) []UsersTrendItem {
	result := make([]UsersTrendItem, 0, days)
	dateMap := make(map[string]*UsersTrendItem)

	for i := range trend {
		dateMap[trend[i].Date] = &trend[i]
	}

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		if item, ok := dateMap[date]; ok {
			result = append(result, *item)
		} else {
			result = append(result, UsersTrendItem{
				Date:  date,
				Count: 0,
			})
		}
	}

	return result
}
