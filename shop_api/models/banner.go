package models

import "time"

type Banner struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"size:200;default:''" json:"title"`
	Image     string    `gorm:"size:500;not null" json:"image"`
	Link      string    `gorm:"size:500;default:''" json:"link"`
	LinkType  int8      `gorm:"not null;default:1" json:"link_type"`
	TargetID  uint64    `gorm:"default:0" json:"target_id"`
	Sort      int       `gorm:"not null;default:0" json:"sort"`
	Status    int8      `gorm:"not null;default:1" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Banner) TableName() string {
	return "banners"
}
