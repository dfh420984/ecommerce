package models

import "shop_api/types"

type Address struct {
	ID         uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64          `gorm:"not null;index" json:"user_id"`
	Consignee  string          `gorm:"size:50;not null" json:"consignee"`
	Phone      string          `gorm:"size:20;not null" json:"phone"`
	Province   string          `gorm:"size:50;not null" json:"province"`
	City       string          `gorm:"size:50;not null" json:"city"`
	District   string          `gorm:"size:50;not null" json:"district"`
	Address    string          `gorm:"size:255;not null" json:"address"`
	PostalCode string          `gorm:"size:20;default:''" json:"postal_code"`
	IsDefault  int8            `gorm:"not null;default:0;index" json:"is_default"`
	CreatedAt  types.LocalTime `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  types.LocalTime `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Address) TableName() string {
	return "addresses"
}
