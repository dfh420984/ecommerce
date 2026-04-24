package services

import (
	"log"
	"shop_api/database"
	"shop_api/models"
	"shop_api/types"
	"time"

	"gorm.io/gorm"
)

// OrderTimeoutService 订单超时服务
type OrderTimeoutService struct{}

var orderTimeoutService *OrderTimeoutService

func GetOrderTimeoutService() *OrderTimeoutService {
	if orderTimeoutService == nil {
		orderTimeoutService = &OrderTimeoutService{}
	}
	return orderTimeoutService
}

// CheckAndCancelTimeoutOrders 检查并取消超时订单
func (s *OrderTimeoutService) CheckAndCancelTimeoutOrders() {
	log.Println("开始检查超时订单...")

	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("检查超时订单失败: %v", r)
		}
	}()

	// 查找超时未支付的订单
	var orders []models.Order
	now := time.Now()
	err := tx.Where("order_status = ? AND pay_status = 0 AND expire_time < ?",
		models.OrderStatusPending, now).Find(&orders).Error

	if err != nil {
		tx.Rollback()
		log.Printf("查询超时订单失败: %v", err)
		return
	}

	if len(orders) == 0 {
		tx.Rollback()
		log.Println("没有超时订单")
		return
	}

	log.Printf("找到 %d 个超时订单，开始处理...", len(orders))

	cancelledCount := 0
	for _, order := range orders {
		if err := s.cancelOrder(tx, &order); err != nil {
			log.Printf("取消订单 %d 失败: %v", order.ID, err)
			continue
		}
		cancelledCount++
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		return
	}

	log.Printf("成功取消 %d 个超时订单", cancelledCount)
}

// cancelOrder 取消单个订单并释放库存
func (s *OrderTimeoutService) cancelOrder(tx *gorm.DB, order *models.Order) error {
	// 更新订单状态
	now := types.Now()
	if err := tx.Model(order).Updates(map[string]interface{}{
		"order_status": models.OrderStatusCancelled,
		"cancel_time":  &now,
	}).Error; err != nil {
		return err
	}

	// 释放预占库存
	var items []models.OrderItem
	if err := tx.Where("order_id = ?", order.ID).Find(&items).Error; err != nil {
		return err
	}

	for _, item := range items {
		// 减少预占库存
		if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
			UpdateColumn("stock_reserved", gorm.Expr("GREATEST(stock_reserved - ?, 0)", item.Quantity)).Error; err != nil {
			log.Printf("释放商品 %d 预占库存失败: %v", item.ProductID, err)
			continue
		}
	}

	// 如果使用了优惠券，返还优惠券
	if order.CouponID != nil && *order.CouponID > 0 {
		if err := tx.Model(&models.UserCoupon{}).Where("id = ?", *order.CouponID).
			Updates(map[string]interface{}{
				"status":   1, // 未使用
				"use_time": nil,
				"order_id": nil,
			}).Error; err != nil {
			log.Printf("返还优惠券失败: %v", err)
		}
	}

	return nil
}

// CreateOrderWithReserveStock 创建订单并预占库存
func (s *OrderTimeoutService) CreateOrderWithReserveStock(order *models.Order, items []models.OrderItem) error {
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 设置订单过期时间（30分钟后）
	expireTime := types.LocalTime(time.Now().Add(30 * time.Minute))
	order.ExpireTime = &expireTime

	// 计算预占库存总数
	totalReserved := 0
	for _, item := range items {
		totalReserved += item.Quantity
	}
	order.StockReserved = totalReserved

	// 创建订单
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建订单项
	for i := range items {
		items[i].OrderID = order.ID
	}
	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 预占库存
	for _, item := range items {
		if err := tx.Model(&models.Product{}).Where("id = ? AND stock_reserved >= 0", item.ProductID).
			UpdateColumn("stock_reserved", gorm.Expr("stock_reserved + ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// ConfirmPayAndDeductStock 支付成功后正式扣减库存
func (s *OrderTimeoutService) ConfirmPayAndDeductStock(orderID uint64) error {
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var order models.Order
	if err := tx.Preload("Items").First(&order, orderID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 正式扣减库存并增加销量
	for _, item := range order.Items {
		// 扣减实际库存
		if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
			UpdateColumn("stock", gorm.Expr("GREATEST(stock - ?, 0)", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 减少预占库存
		if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
			UpdateColumn("stock_reserved", gorm.Expr("GREATEST(stock_reserved - ?, 0)", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 增加销量
		if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
			UpdateColumn("sales", gorm.Expr("sales + ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 清除订单的预占库存标记
	if err := tx.Model(&order).Update("stock_reserved", 0).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
