package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type Orders struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	gorm.Model
	Name string `json:"name" gorm:"uniqueIndex"`
	Cost int    `json:"cost"`
	//RegionID *uint   `json:"region_id,omitempty"`
	//Region   *Region `json:"regions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Weight float32 `json:"weight"`

	DeliveryHours pq.StringArray `json:"delivery_hours" gorm:"type:text[]"`
	Regions       int            `json:"regions"`

	CompletedTime time.Time `json:"completed_time,omitempty"`

	CourierID *uint    `json:"courier_id,omitempty"`
	Courier   *Courier `json:"courier,omitempty"`
}

type ResponseOrder struct {
	Cost          int       `json:"cost"`
	DeliveryHours []string  `json:"delivery_hours"`
	OrderID       uint      `json:"order_id"`
	Regions       int       `json:"regions"`
	Weight        float32   `json:"weight"`
	CompletedTime time.Time `json:"completed_time,omitempty"`
}

type CompleteOrders struct {
	Orders []CompleteOrder `json:"complete_info"`
}

type CompleteOrder struct {
	ID            uint      `json:"order_id,omitempty"`
	CompletedTime time.Time `json:"completed_time,omitempty"`
	CourierID     uint      `json:"courier_id,omitempty"`
}

type AssignCouriers struct {
	Date     time.Time       `json:"date"`
	Couriers []AssignCourier `json:"orders"`
}

type AssignCourier struct {
	CourierId uint           `json:"courier_id"`
	Orders    []AssignOrders `json:"orders"`
}

type AssignOrders struct {
	GroupOrderId uint            `json:"group_order_id"`
	Orders       []ResponseOrder `json:"orders"`
}
