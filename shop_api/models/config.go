package models

import "shop_api/types"

type Config struct {
	ID          uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string          `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Value       string          `gorm:"type:text" json:"value"`
	Description string          `gorm:"size:500;default:''" json:"description"`
	CreatedAt   types.LocalTime `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   types.LocalTime `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Config) TableName() string {
	return "configs"
}

type OperationLog struct {
	ID        uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64          `gorm:"default:0;index" json:"user_id"`
	Username  string          `gorm:"size:50;default:''" json:"username"`
	Action    string          `gorm:"size:100;not null;index" json:"action"`
	Target    string          `gorm:"size:100;default:''" json:"target"`
	TargetID  uint64          `gorm:"default:0" json:"target_id"`
	Content   string          `gorm:"type:text" json:"content"`
	IP        string          `gorm:"size:50;default:''" json:"ip"`
	UserAgent string          `gorm:"size:500;default:''" json:"user_agent"`
	CreatedAt types.LocalTime `gorm:"autoCreateTime;index" json:"created_at"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}
