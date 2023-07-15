package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id" gorm:"index;primaryKey"`
	UserID    uint      `json:"-" gorm:"index"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

type OrderItem struct {
	OrderID   uuid.UUID `json:"-" gorm:"primaryKey;index"`
	Order     Order     `json:"order" gorm:"foreignKey:OrderID;references:ID"`
	ProductID uint      `json:"-" gorm:"primaryKey;index"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID;references:ID"`
	Quantity  uint      `json:"quantity"`
	UnitPrice uint      `json:"unit_price"`
}
