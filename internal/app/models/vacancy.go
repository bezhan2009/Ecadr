package models

import (
	"gorm.io/gorm"
)

type Vacancy struct {
	gorm.Model

	Title       string `json:"title"`
	Description string `json:"description" gorm:"type:text"`
	WorkerID    int    `json:"worker_id" gorm:"not null"`
	Worker      User   `json:"-" gorm:"foreignkey:WorkerID"`
	Contact     string `json:"contact"` // Куда отправлять ответы
}
