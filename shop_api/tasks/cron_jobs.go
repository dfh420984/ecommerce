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

	log.Println("定时任务启动成功")
}
