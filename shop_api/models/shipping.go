package models

import "shop_api/types"

// ShippingTemplate 运费模板
type ShippingTemplate struct {
	ID               uint64                   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string                   `gorm:"size:100;not null" json:"name"`                            // 模板名称
	IsDefault        int8                     `gorm:"not null;default:0" json:"is_default"`                     // 是否默认模板:0否,1是
	FreeShippingType int8                     `gorm:"not null;default:1" json:"free_shipping_type"`             // 包邮类型:1满额包邮,2满件包邮,3满额或满件,4不包邮
	FreeAmount       float64                  `gorm:"type:decimal(10,2);not null;default:0" json:"free_amount"` // 包邮金额阈值
	FreeQuantity     int                      `gorm:"not null;default:0" json:"free_quantity"`                  // 包邮数量阈值
	BaseFee          float64                  `gorm:"type:decimal(10,2);not null;default:0" json:"base_fee"`    // 基础运费
	Status           int8                     `gorm:"not null;default:1;index" json:"status"`                   // 状态:0禁用,1启用
	CreatedAt        types.LocalTime          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        types.LocalTime          `gorm:"autoUpdateTime" json:"updated_at"`
	Regions          []ShippingTemplateRegion `gorm:"foreignKey:TemplateID" json:"regions,omitempty"` // 区域配置
}

func (ShippingTemplate) TableName() string {
	return "shipping_templates"
}

// ShippingTemplateRegion 运费模板区域配置
type ShippingTemplateRegion struct {
	ID           uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TemplateID   uint64          `gorm:"not null;index" json:"template_id"`                        // 模板ID
	Province     string          `gorm:"size:50;not null" json:"province"`                         // 省份
	City         string          `gorm:"size:50;not null" json:"city"`                             // 城市
	District     string          `gorm:"size:50;default:''" json:"district"`                       // 区县（可选）
	Fee          float64         `gorm:"type:decimal(10,2);not null;default:0" json:"fee"`         // 该区域运费
	FreeAmount   float64         `gorm:"type:decimal(10,2);not null;default:0" json:"free_amount"` // 该区域包邮金额
	FreeQuantity int             `gorm:"not null;default:0" json:"free_quantity"`                  // 该区域包邮数量
	CreatedAt    types.LocalTime `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    types.LocalTime `gorm:"autoUpdateTime" json:"updated_at"`
}

func (ShippingTemplateRegion) TableName() string {
	return "shipping_template_regions"
}

const (
	FreeShippingTypeAmount   int8 = 1 // 满额包邮
	FreeShippingTypeQuantity int8 = 2 // 满件包邮
	FreeShippingTypeEither   int8 = 3 // 满额或满件
	FreeShippingTypeNone     int8 = 4 // 不包邮
)
