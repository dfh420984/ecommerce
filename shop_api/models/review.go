package models

import "shop_api/types"

// ProductReview 商品评价
type ProductReview struct {
	ID           uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID      uint64           `gorm:"not null;index" json:"order_id"`
	UserID       uint64           `gorm:"not null;index" json:"user_id"`
	ProductID    uint64           `gorm:"not null;index" json:"product_id"`
	OrderItemID  uint64           `gorm:"not null;uniqueIndex" json:"order_item_id"`
	Rating       int8             `gorm:"not null" json:"rating"` // 1-5星
	Content      string           `gorm:"type:text" json:"content"`
	Images       StringArray      `gorm:"type:json" json:"images"`
	ReplyContent string           `gorm:"type:text" json:"reply_content"`
	ReplyTime    *types.LocalTime `json:"reply_time"`
	IsAnonymous  int8             `gorm:"default:0" json:"is_anonymous"`
	Status       int8             `gorm:"default:1;index" json:"status"` // 1正常 0隐藏
	HelpfulCount int              `gorm:"default:0" json:"helpful_count"`
	CreatedAt    types.LocalTime  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    types.LocalTime  `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (ProductReview) TableName() string {
	return "product_reviews"
}

// ReviewStats 评价统计
type ReviewStats struct {
	AverageRating float64 `json:"average_rating"`
	TotalCount    int64   `json:"total_count"`
	FiveStar      int64   `json:"five_star"`
	FourStar      int64   `json:"four_star"`
	ThreeStar     int64   `json:"three_star"`
	TwoStar       int64   `json:"two_star"`
	OneStar       int64   `json:"one_star"`
}
