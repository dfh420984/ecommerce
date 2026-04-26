package models

import "shop_api/types"

// ProductLike 商品点赞
type ProductLike struct {
	ID        uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint64          `gorm:"not null;uniqueIndex:uk_product_user" json:"product_id"`
	UserID    uint64          `gorm:"not null;uniqueIndex:uk_product_user" json:"user_id"`
	CreatedAt types.LocalTime `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt types.LocalTime `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ProductLike) TableName() string {
	return "product_likes"
}
