package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Title       string `json:"title"`
	Price       uint   `json:"price"`
	Inventory   uint   `json:"inventory"`
	Image       string `json:"image"`
	Description string `json:"description"`
}
