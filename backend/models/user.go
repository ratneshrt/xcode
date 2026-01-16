package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name"`
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
}
