package models

import "time"

type Item struct {
	ID                    uint        `json:"id" gorm:"primaryKey"`
	ProductID             uint        `json:"product_id" gorm:"not null;index"`
	IsActive              bool        `json:"is_active"`
	ActivateIfPossible    bool        `json:"activate_if_possible"`
	Name                  string      `json:"name"`
	RefID                 string      `json:"ref_id"`
	PackagedHeight        float64     `json:"packaged_height"`
	PackagedLength        float64     `json:"packaged_length"`
	PackagedWidth         float64     `json:"packaged_width"`
	PackagedWeightKg      float64     `json:"packaged_weight_kg"`
	Height                *float64    `json:"height"`
	Length                *float64    `json:"length"`
	Width                 *float64    `json:"width"`
	WeightKg              *float64    `json:"weight_kg"`
	CubicWeight           float64     `json:"cubic_weight"`
	IsKit                 bool        `json:"is_kit"`
	RewardValue           *float64    `json:"reward_value"`
	EstimatedDateArrival  *time.Time  `json:"estimated_date_arrival"`
	ManufacturerCode      string      `json:"manufacturer_code"`
	CommercialConditionID uint        `json:"commercial_condition_id"`
	MeasurementUnit       string      `json:"measurement_unit"`
	UnitMultiplier        float64     `json:"unit_multiplier"`
	ModalType             *string     `json:"modal_type"`
	KitItensSellApart     bool        `json:"kit_itens_sell_apart"`
	SKU                   string      `json:"sku" gorm:"uniqueIndex;not null"`
	Price                 float64     `json:"price" gorm:"not null"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
	Images                []ItemImage `json:"images,omitempty"`
}

func (Item) TableName() string {
	return "items"
}
