package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItem struct {
	CreatedAt time.Time `json:"created_at"`
	CartID    uuid.UUID `json:"-" gorm:"primaryKey;index"`
	Cart      Cart      `json:"cart" gorm:"foreignKey:CartID;references:ID"`
	ProductID uint      `json:"-" gorm:"primaryKey;index"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID;references:ID"`
	Quantity  uint      `json:"quantity"`
}

type Cart struct {
	ID        uuid.UUID      `json:"id" gorm:"index;primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
