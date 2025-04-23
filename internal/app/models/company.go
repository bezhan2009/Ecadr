package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model

	Title       string         `json:"title" gorm:"unique;not null"`
	Description string         `json:"description" gorm:"not null"`
	Criteria    pq.StringArray `json:"criteria" gorm:"type:text[]"`
}
