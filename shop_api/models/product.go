package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"shop_api/types"
)

type Product struct {
	ID            uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string          `gorm:"size:200;not null" json:"name"`
	CategoryID    uint64          `gorm:"not null;index" json:"category_id"`
	Price         float64         `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPrice float64         `gorm:"type:decimal(10,2);default:0" json:"original_price"`
	Cost          float64         `gorm:"type:decimal(10,2);default:0" json:"cost"`
	Stock         int             `gorm:"not null;default:0" json:"stock"`
	Sales         int             `gorm:"not null;default:0" json:"sales"`
	Images        StringArray     `gorm:"type:text" json:"images"`
	Description   string          `gorm:"type:text" json:"description"`
	Content       string          `gorm:"type:text" json:"content"`
	Specs         StringArray     `gorm:"type:text" json:"specs"`
	IsOnline      int8            `gorm:"not null;default:1;index" json:"is_online"`
	IsRecommend   int8            `gorm:"not null;default:0;index" json:"is_recommend"`
	IsNew         int8            `gorm:"not null;default:0" json:"is_new"`
	Sort          int             `gorm:"not null;default:0" json:"sort"`
	Status        int8            `gorm:"not null;default:1;index" json:"status"`
	CreatedAt     types.LocalTime `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     types.LocalTime `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}

type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal StringArray value")
	}
	return json.Unmarshal(bytes, s)
}
