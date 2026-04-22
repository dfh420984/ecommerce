package models

import "time"

type Cart struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"not null;index" json:"user_id"`
	ProductID uint64    `gorm:"not null;index" json:"product_id"`
	SkuID     uint64    `gorm:"default:0" json:"sku_id"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	Selected  int8      `gorm:"not null;default:1" json:"selected"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (Cart) TableName() string {
	return "cart"
}
