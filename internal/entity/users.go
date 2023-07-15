package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}
