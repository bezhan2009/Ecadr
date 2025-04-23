package models

import "gorm.io/gorm"

type Criteria struct {
	gorm.Model

	Title     string  `json:"title"`
	VacancyID uint    `json:"vacancy_id"`
	Vacancy   Vacancy `json:"vacancy" gorm:"foreignKey:VacancyID"`
}
