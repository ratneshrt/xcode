package models

import "gorm.io/gorm"

type Submission struct {
	gorm.Model

	UserID    uint `gorm:"index;not null"`
	ProblemID uint `gorm:"index;not null"`

	Language string `gorm:"not null"`
	Code     string `gorm:"type:text;not null"`

	Status string `gorm:"index;not null"`
	Result string `gorm:"type:text"`
}
