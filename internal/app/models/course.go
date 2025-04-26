package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type Course struct {
	gorm.Model

	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Subject     string         `json:"subject" gorm:"not null"`
	Tags        pq.StringArray `json:"tags" gorm:"type:text[]"`
	WorkerID    int            `json:"worker_id" gorm:"not null"`
	Worker      User           `json:"-" gorm:"foreignkey:WorkerID"`
	CompanyID   int            `json:"company_id" gorm:"not null"`
	Company     Company        `json:"-" gorm:"foreignkey:CompanyID"`
	StartDate   time.Time      `json:"start_date" gorm:"type:date;"`
	EndDate     time.Time      `json:"end_date" gorm:"type:date;"`
}
