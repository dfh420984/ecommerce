package models

import "time"

type Order struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo        string     `gorm:"size:32;not null;uniqueIndex" json:"order_no"`
	UserID         uint64     `gorm:"not null;index" json:"user_id"`
	OrderStatus    int8       `gorm:"not null;default:1;index" json:"order_status"`
	PayStatus      int8       `gorm:"not null;default:0;index" json:"pay_status"`
	PayType        int8       `gorm:"not null;default:0" json:"pay_type"`
	TotalAmount    float64    `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	DiscountAmount float64    `gorm:"type:decimal(10,2);not null;default:0" json:"discount_amount"`
	FreightAmount  float64    `gorm:"type:decimal(10,2);not null;default:0" json:"freight_amount"`
	PayAmount      float64    `gorm:"type:decimal(10,2);not null" json:"pay_amount"`
	PayTime        *time.Time `json:"pay_time"`
	Consignee      string     `gorm:"size:50;not null" json:"consignee"`
	Phone          string     `gorm:"size:20;not null" json:"phone"`
	Province       string     `gorm:"size:50;not null" json:"province"`
	City           string     `gorm:"size:50;not null" json:"city"`
	District       string     `gorm:"size:50;not null" json:"district"`
	Address        string     `gorm:"size:255;not null" json:"address"`
	Remark         string     `gorm:"size:500;default:''" json:"remark"`
	CancelTime     *time.Time `json:"cancel_time"`
	CompleteTime   *time.Time `json:"complete_time"`
	CreatedAt      time.Time  `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	Items          []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

const (
	OrderStatusPending   int8 = 1
	OrderStatusPaid      int8 = 2
	OrderStatusShipped   int8 = 3
	OrderStatusReceived  int8 = 4
	OrderStatusCompleted int8 = 5
	OrderStatusCancelled int8 = 6
	OrderStatusRefund    int8 = 7
	OrderStatusRefunded  int8 = 8
)

const (
	PayStatusUnpaid  int8 = 0
	PayStatusPaid    int8 = 1
	PayStatusRefunded int8 = 2
)

const (
	PayTypeNone     int8 = 0
	PayTypeWechat   int8 = 1
	PayTypeAlipay   int8 = 2
	PayTypeBankCard int8 = 3
)
