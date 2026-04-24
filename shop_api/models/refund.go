package models

import "shop_api/types"

// RefundApplication 退款申请
type RefundApplication struct {
	ID            uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID       uint64           `gorm:"not null;index" json:"order_id"`
	UserID        uint64           `gorm:"not null;index" json:"user_id"`
	RefundType    string           `gorm:"size:50;not null;default:'refund_only'" json:"refund_type"` // refund_only/exchange
	Reason        string           `gorm:"size:500;not null" json:"reason"`
	Images        StringArray      `gorm:"type:json" json:"images"`
	RefundAmount  float64          `gorm:"type:decimal(10,2);not null" json:"refund_amount"`
	Status        int8             `gorm:"not null;default:1;index" json:"status"` // 1待审核 2已通过 3已拒绝 4退款中 5已退款
	HandlerID     *uint64          `gorm:"default:null" json:"handler_id"`
	HandlerReply  string           `gorm:"type:text" json:"handler_reply"`
	RefundTradeNo string           `gorm:"size:100" json:"refund_trade_no"`
	RefundedAt    *types.LocalTime `json:"refunded_at"`
	CreatedAt     types.LocalTime  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     types.LocalTime  `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Order   *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	User    *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Handler *User  `gorm:"foreignKey:HandlerID" json:"handler,omitempty"`
}

func (RefundApplication) TableName() string {
	return "refund_applications"
}

// 退款状态常量
const (
	RefundStatusPending   int8 = 1 // 待审核
	RefundStatusApproved  int8 = 2 // 已通过
	RefundStatusRejected  int8 = 3 // 已拒绝
	RefundStatusRefunding int8 = 4 // 退款中
	RefundStatusRefunded  int8 = 5 // 已退款
)

// 退款类型常量
const (
	RefundTypeOnly     = "refund_only" // 仅退款
	RefundTypeExchange = "exchange"    // 换货
)
