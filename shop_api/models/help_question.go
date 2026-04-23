package models

import "shop_api/types"

type HelpCategory struct {
	ID          uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string          `gorm:"size:50;not null" json:"name"`
	Sort        int             `gorm:"default:0" json:"sort"`
	Status      int             `gorm:"default:1" json:"status"` // 1:启用 0:禁用
	Description string          `gorm:"size:200;default:''" json:"description"`
	CreatedAt   types.LocalTime `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   types.LocalTime `gorm:"autoUpdateTime" json:"updated_at"`
}

func (HelpCategory) TableName() string {
	return "help_categories"
}

type HelpQuestion struct {
	ID         uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID uint64          `gorm:"not null;index" json:"category_id"`
	Title      string          `gorm:"size:200;not null" json:"title"`
	Answer     string          `gorm:"type:text;not null" json:"answer"`
	Sort       int             `gorm:"default:0" json:"sort"`
	Status     int             `gorm:"default:1" json:"status"` // 1:启用 0:禁用
	CreatedAt  types.LocalTime `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  types.LocalTime `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联分类
	Category HelpCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (HelpQuestion) TableName() string {
	return "help_questions"
}
