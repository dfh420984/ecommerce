package database

import (
	"fmt"
	"log"
	"shop_api/config"
	"shop_api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.DatabaseConfig) error {
	var err error
	DB, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB failed: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	log.Println("Database connection established")
	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Banner{},
		&models.Address{},
		&models.Cart{},
		&models.Order{},
		&models.OrderItem{},
		&models.PayLog{},
		&models.Config{},
		&models.OperationLog{},
	)
}

func GetDB() *gorm.DB {
	return DB
}
