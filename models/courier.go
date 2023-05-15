package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Couriers struct {
	Couriers []Courier `json:"couriers"`
}

type Courier struct {
	gorm.Model

	TypeID *string `json:"courier_type,omitempty"`
	Type   *Type   `json:"courier_type_full,omitempty"`

	Regions      pq.Int64Array  `json:"regions" gorm:"type:integer[]"`
	WorkingHours pq.StringArray `json:"working_hours" gorm:"type:text[]"`
}

type ResponseCouriers struct {
	Couriers []ResponseCourier `json:"couriers"`
	Limit    int               `json:"limit,omitempty"`
	Offset   int               `json:"offset,omitempty"`
}

type ResponseCourier struct {
	CourierId    uint     `json:"courier_id"`
	CourierType  string   `json:"courier_type"`
	Regions      []int    `json:"regions"`
	WorkingHours []string `json:"working_hours,omitempty"`
}

type ResponseMetaInfo struct {
	ResponseCourier
	Rating   int `json:"rating"`
	Earnings int `json:"earnings"`
}
