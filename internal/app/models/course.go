package models

import (
	"gorm.io/gorm"
	"time"
)

type Course struct {
	gorm.Model

	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	WorkerID    int       `json:"worker_id" gorm:"not null"`
	Worker      User      `json:"-" gorm:"foreignkey:WorkerID"`
	StartDate   time.Time `json:"start_date" gorm:"type:date;"`
	EndDate     time.Time `json:"end_date" gorm:"type:date;"`
}
