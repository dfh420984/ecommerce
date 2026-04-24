package models

import "shop_api/types"

// UserFavorite 用户收藏
type UserFavorite struct {
	ID        uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64          `gorm:"not null;uniqueIndex:uk_user_product" json:"user_id"`
	ProductID uint64          `gorm:"not null;uniqueIndex:uk_user_product" json:"product_id"`
	CreatedAt types.LocalTime `gorm:"autoCreateTime" json:"created_at"`

	// 关联
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (UserFavorite) TableName() string {
	return "user_favorites"
}

// BrowseHistory 浏览历史
type BrowseHistory struct {
	ID        uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64          `gorm:"not null;index:idx_user_time" json:"user_id"`
	ProductID uint64          `gorm:"not null" json:"product_id"`
	ViewedAt  types.LocalTime `gorm:"autoCreateTime;index:idx_user_time" json:"viewed_at"`

	// 关联
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (BrowseHistory) TableName() string {
	return "browse_histories"
}
