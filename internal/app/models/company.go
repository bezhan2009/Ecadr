package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model

	Title       string         `json:"title" gorm:"unique;not null"`
	Description string         `json:"description" gorm:"not null"`
	Subjects    pq.StringArray `json:"subjects" gorm:"type:text[]"`
	Criteria    pq.StringArray `json:"criteria" gorm:"type:text[]"`
	WorkerID    uint           `json:"worker_id" gorm:"not null"`
	Worker      User           `json:"-" gorm:"foreignkey:WorkerID"`
}
