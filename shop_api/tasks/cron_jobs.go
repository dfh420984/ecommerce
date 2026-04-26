package tasks

import (
	"log"
	"shop_api/services"
	"time"
)

// StartCronJobs 启动定时任务
func StartCronJobs() {
	log.Println("启动定时任务...")

	// 每5分钟检查一次超时订单
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("执行定时任务：检查超时订单")
			services.GetOrderTimeoutService().CheckAndCancelTimeoutOrders()
		}
	}()

	// 每天凌晨2点检查并自动完成已收货超时的订单
	go func() {
		for {
			now := time.Now()
			// 计算到下一个凌晨2点的时间
			next := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}
			duration := next.Sub(now)

			log.Printf("下次自动完成订单检查将在 %v 后执行", duration)
			time.Sleep(duration)

			log.Println("执行定时任务：自动完成订单")
			services.GetOrderTimeoutService().CheckAndAutoCompleteOrders()
		}
	}()

	log.Println("定时任务启动成功")
}
