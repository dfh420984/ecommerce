package models

import (
	"shop_api/types"
)

type User struct {
	ID        uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string          `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password  string          `gorm:"size:255;not null" json:"-"`
	Nickname  string          `gorm:"size:100;default:''" json:"nickname"`
	Avatar    string          `gorm:"size:500;default:''" json:"avatar"`
	Phone     string          `gorm:"size:20;default:''" json:"phone"`
	Email     string          `gorm:"size:100;default:''" json:"email"`
	Status    int8            `gorm:"not null;default:1" json:"status"`
	UserType  int8            `gorm:"not null;default:1" json:"user_type"`
	OpenID    string          `gorm:"column:openid;size:100;default:'';index" json:"openid"`
	CreatedAt types.LocalTime `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt types.LocalTime `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

const (
	UserStatusDisabled int8 = 0
	UserStatusEnabled  int8 = 1
	UserTypeNormal     int8 = 1
	UserTypeAdmin      int8 = 2
)
