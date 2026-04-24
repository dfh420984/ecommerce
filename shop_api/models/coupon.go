package models

import (
	"shop_api/types"
	"time"
)

// Coupon 优惠券模板
type Coupon struct {
	ID            uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string           `gorm:"size:100;not null" json:"name"`
	Type          int8             `gorm:"not null;default:1" json:"type"` // 1-满减券,2-折扣券,3-无门槛券
	DiscountValue float64          `gorm:"type:decimal(10,2);not null" json:"discount_value"`
	MinAmount     float64          `gorm:"type:decimal(10,2);not null;default:0" json:"min_amount"`
	MaxDiscount   float64          `gorm:"type:decimal(10,2);not null;default:0" json:"max_discount"`
	TotalCount    int              `gorm:"not null;default:0" json:"total_count"` // 0表示不限
	ReceivedCount int              `gorm:"not null;default:0" json:"received_count"`
	UsedCount     int              `gorm:"not null;default:0" json:"used_count"`
	PerUserLimit  int              `gorm:"not null;default:1" json:"per_user_limit"`
	ValidType     int8             `gorm:"not null;default:1" json:"valid_type"` // 1-固定时间,2-领取后N天
	StartTime     *types.LocalTime `gorm:"type:timestamp;null" json:"start_time"`
	EndTime       *types.LocalTime `gorm:"type:timestamp;null" json:"end_time"`
	ValidDays     int              `gorm:"not null;default:0" json:"valid_days"`
	Status        int8             `gorm:"not null;default:1;index" json:"status"`      // 0-禁用,1-启用
	IsNewUser     int8             `gorm:"not null;default:0;index" json:"is_new_user"` // 是否新人券
	Description   string           `gorm:"type:text" json:"description"`
	CreatedAt     types.LocalTime  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     types.LocalTime  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Coupon) TableName() string {
	return "coupons"
}

// IsAvailable 检查优惠券是否可用
func (c *Coupon) IsAvailable() bool {
	if c.Status != 1 {
		return false
	}

	now := time.Now()

	// 固定时间类型
	if c.ValidType == 1 {
		if c.StartTime != nil && now.Before(time.Time(*c.StartTime)) {
			return false
		}
		if c.EndTime != nil && now.After(time.Time(*c.EndTime)) {
			return false
		}
	}

	// 检查库存
	if c.TotalCount > 0 && c.ReceivedCount >= c.TotalCount {
		return false
	}

	return true
}

// UserCoupon 用户优惠券
type UserCoupon struct {
	ID          uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	CouponID    uint64           `gorm:"not null;index" json:"coupon_id"`
	UserID      uint64           `gorm:"not null;index" json:"user_id"`
	Status      int8             `gorm:"not null;default:1;index" json:"status"` // 1-未使用,2-已使用,3-已过期
	ReceiveTime types.LocalTime  `gorm:"autoCreateTime" json:"receive_time"`
	UseTime     *types.LocalTime `gorm:"type:timestamp;null" json:"use_time"`
	OrderID     *uint64          `gorm:"default:null" json:"order_id"`
	ExpireTime  types.LocalTime  `gorm:"type:timestamp;not null" json:"expire_time"`
	CreatedAt   types.LocalTime  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   types.LocalTime  `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联查询
	Coupon *Coupon `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
}

func (UserCoupon) TableName() string {
	return "user_coupons"
}

// IsUsable 检查用户优惠券是否可用
func (uc *UserCoupon) IsUsable() bool {
	if uc.Status != 1 {
		return false
	}
	if time.Now().After(time.Time(uc.ExpireTime)) {
		return false
	}
	return true
}
