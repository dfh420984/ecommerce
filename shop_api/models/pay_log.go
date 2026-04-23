package models

import "shop_api/types"

type PayLog struct {
	ID         uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID    uint64          `gorm:"not null;index" json:"order_id"`
	OrderNo    string          `gorm:"size:32;not null;index" json:"order_no"`
	TradeNo    string          `gorm:"size:64;default:'';index" json:"trade_no"`
	PayType    int8            `gorm:"not null" json:"pay_type"`
	PayStatus  int8            `gorm:"not null" json:"pay_status"`
	PayAmount  float64         `gorm:"type:decimal(10,2);not null" json:"pay_amount"`
	NotifyData string          `gorm:"type:text" json:"notify_data"`
	CreatedAt  types.LocalTime `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  types.LocalTime `gorm:"autoUpdateTime" json:"updated_at"`
}

func (PayLog) TableName() string {
	return "pay_logs"
}
