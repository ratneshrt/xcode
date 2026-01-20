package models

import (
	"time"

	"gorm.io/gorm"
)

type Problem struct {
	gorm.Model

	Title       string
	Slug        string `gorm:"uniqueIndex"` // it is nothing but url-friendly identifier
	Description string

	Difficulty  string
	Constraints string

	InputFormat  string
	OutputFormat string

	TimeLimitMs   int
	MemoryLimitMs string

	Status      string
	PublishedAt *time.Time
}

type ProblemExample struct {
	gorm.Model
	ProblemID uint

	Input       string
	Output      string
	Explanation string
}

type ProblemTestCase struct {
	gorm.Model
	ProblemID uint

	Input          string
	ExpectedOutput string
	IsHidden       bool
}
