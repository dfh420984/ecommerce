package models

import "time"

type OrderItem struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID      uint64    `gorm:"not null;index" json:"order_id"`
	ProductID    uint64    `gorm:"not null;index" json:"product_id"`
	ProductName  string    `gorm:"size:200;not null" json:"product_name"`
	ProductImage string    `gorm:"size:500;not null" json:"product_image"`
	Price        float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Quantity     int       `gorm:"not null" json:"quantity"`
	SkuID        uint64    `gorm:"default:0" json:"sku_id"`
	SkuName      string    `gorm:"size:200;default:''" json:"sku_name"`
	Subtotal     float64   `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
