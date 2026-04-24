package main

import (
	"fmt"
	"log"
	"shop_api/config"
	"shop_api/database"
	"shop_api/routes"
	"shop_api/tasks"
	"shop_api/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := database.Init(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	if err := database.InitRedis(&cfg.Redis); err != nil {
		log.Printf("Warning: Failed to connect Redis: %v", err)
	}

	utils.SetJWTSecret(cfg.App.JWTSecret)

	if cfg.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := routes.SetupRouter()

	// 启动定时任务
	tasks.StartCronJobs()

	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	log.Printf("Server starting on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
