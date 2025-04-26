package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Vacancy struct {
	gorm.Model
	Title       string         `json:"title"`
	Description string         `json:"description" gorm:"type:text"`
	Subject     string         `json:"subject" gorm:"type:text"`
	Tags        pq.StringArray `json:"tags" gorm:"type:text[]"`
	WorkerID    int            `json:"worker_id" gorm:"not null"`
	CompanyID   uint           `json:"company_id"`
	Contact     string         `json:"contact"`                                       // Куда отправлять ответы
	Salary      Salary         `json:"salary" gorm:"embedded;embeddedPrefix:salary_"` // Embedding salary
	Location    string         `json:"location"`
	Experience  string         `json:"experience"`
	Type        string         `json:"type"`

	Worker   User       `json:"-" gorm:"foreignkey:WorkerID"`
	Company  Company    `json:"company" gorm:"foreignKey:CompanyID"`
	Criteria []Criteria `json:"criteria" gorm:"foreignkey:VacancyID"`
}

type Salary struct {
	Currency string `json:"currency"`
	Min      int    `json:"min"`
	Max      int    `json:"max"`
	Period   string `json:"period"`
}
