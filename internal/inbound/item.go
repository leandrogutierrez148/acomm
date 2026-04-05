package inbound

import (
	"time"

	"github.com/leandrogutierrez148/acomm/internal/models"
)

type CreateItemRequest struct {
	ProductID             uint       `json:"product_id"`
	SKU                   string     `json:"sku"`
	Price                 float64    `json:"price"`
	IsActive              bool       `json:"is_active"`
	ActivateIfPossible    bool       `json:"activate_if_possible"`
	Name                  string     `json:"name"`
	RefID                 string     `json:"ref_id"`
	PackagedHeight        float64    `json:"packaged_height"`
	PackagedLength        float64    `json:"packaged_length"`
	PackagedWidth         float64    `json:"packaged_width"`
	PackagedWeightKg      float64    `json:"packaged_weight_kg"`
	Height                *float64   `json:"height"`
	Length                *float64   `json:"length"`
	Width                 *float64   `json:"width"`
	WeightKg              *float64   `json:"weight_kg"`
	CubicWeight           float64    `json:"cubic_weight"`
	IsKit                 bool       `json:"is_kit"`
	RewardValue           *float64   `json:"reward_value"`
	EstimatedDateArrival  *time.Time `json:"estimated_date_arrival"`
	ManufacturerCode      string     `json:"manufacturer_code"`
	CommercialConditionID uint       `json:"commercial_condition_id"`
	MeasurementUnit       string     `json:"measurement_unit"`
	UnitMultiplier        float64    `json:"unit_multiplier"`
	ModalType             *string    `json:"modal_type"`
	KitItensSellApart     bool       `json:"kit_itens_sell_apart"`
}

func (r *CreateItemRequest) ToDomain() *models.Item {
	return &models.Item{
		ProductID:             r.ProductID,
		SKU:                   r.SKU,
		Price:                 r.Price,
		IsActive:              r.IsActive,
		ActivateIfPossible:    r.ActivateIfPossible,
		Name:                  r.Name,
		RefID:                 r.RefID,
		PackagedHeight:        r.PackagedHeight,
		PackagedLength:        r.PackagedLength,
		PackagedWidth:         r.PackagedWidth,
		PackagedWeightKg:      r.PackagedWeightKg,
		Height:                r.Height,
		Length:                r.Length,
		Width:                 r.Width,
		WeightKg:              r.WeightKg,
		CubicWeight:           r.CubicWeight,
		IsKit:                 r.IsKit,
		RewardValue:           r.RewardValue,
		EstimatedDateArrival:  r.EstimatedDateArrival,
		ManufacturerCode:      r.ManufacturerCode,
		CommercialConditionID: r.CommercialConditionID,
		MeasurementUnit:       r.MeasurementUnit,
		UnitMultiplier:        r.UnitMultiplier,
		ModalType:             r.ModalType,
		KitItensSellApart:     r.KitItensSellApart,
	}
}

type UpdateItemRequest struct {
	ProductID             uint       `json:"product_id"`
	SKU                   string     `json:"sku"`
	Price                 float64    `json:"price"`
	IsActive              bool       `json:"is_active"`
	ActivateIfPossible    bool       `json:"activate_if_possible"`
	Name                  string     `json:"name"`
	RefID                 string     `json:"ref_id"`
	PackagedHeight        float64    `json:"packaged_height"`
	PackagedLength        float64    `json:"packaged_length"`
	PackagedWidth         float64    `json:"packaged_width"`
	PackagedWeightKg      float64    `json:"packaged_weight_kg"`
	Height                *float64   `json:"height"`
	Length                *float64   `json:"length"`
	Width                 *float64   `json:"width"`
	WeightKg              *float64   `json:"weight_kg"`
	CubicWeight           float64    `json:"cubic_weight"`
	IsKit                 bool       `json:"is_kit"`
	RewardValue           *float64   `json:"reward_value"`
	EstimatedDateArrival  *time.Time `json:"estimated_date_arrival"`
	ManufacturerCode      string     `json:"manufacturer_code"`
	CommercialConditionID uint       `json:"commercial_condition_id"`
	MeasurementUnit       string     `json:"measurement_unit"`
	UnitMultiplier        float64    `json:"unit_multiplier"`
	ModalType             *string    `json:"modal_type"`
	KitItensSellApart     bool       `json:"kit_itens_sell_apart"`
}

func (r *UpdateItemRequest) ToDomain() *models.Item {
	return &models.Item{
		ProductID:             r.ProductID,
		SKU:                   r.SKU,
		Price:                 r.Price,
		IsActive:              r.IsActive,
		ActivateIfPossible:    r.ActivateIfPossible,
		Name:                  r.Name,
		RefID:                 r.RefID,
		PackagedHeight:        r.PackagedHeight,
		PackagedLength:        r.PackagedLength,
		PackagedWidth:         r.PackagedWidth,
		PackagedWeightKg:      r.PackagedWeightKg,
		Height:                r.Height,
		Length:                r.Length,
		Width:                 r.Width,
		WeightKg:              r.WeightKg,
		CubicWeight:           r.CubicWeight,
		IsKit:                 r.IsKit,
		RewardValue:           r.RewardValue,
		EstimatedDateArrival:  r.EstimatedDateArrival,
		ManufacturerCode:      r.ManufacturerCode,
		CommercialConditionID: r.CommercialConditionID,
		MeasurementUnit:       r.MeasurementUnit,
		UnitMultiplier:        r.UnitMultiplier,
		ModalType:             r.ModalType,
		KitItensSellApart:     r.KitItensSellApart,
	}
}
