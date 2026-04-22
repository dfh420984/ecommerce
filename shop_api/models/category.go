package models

import "time"

type Category struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"size:100;not null" json:"name"`
	ParentID  uint64     `gorm:"not null;default:0;index" json:"parent_id"`
	Level     int        `gorm:"not null;default:1" json:"level"`
	Sort      int        `gorm:"not null;default:0" json:"sort"`
	Icon      string     `gorm:"size:500;default:''" json:"icon"`
	Image     string     `gorm:"size:500;default:''" json:"image"`
	Status    int8       `gorm:"not null;default:1;index" json:"status"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	Children  []Category `gorm:"-" json:"children,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
