package models

type Type struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Ratio       int
	RatingRatio int
	MaxWeight   float32
	MaxOrders   int
	MaxRegions  int
}
