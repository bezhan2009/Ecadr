package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model

	Description string `json:"description" gorm:"type:text; not null"`
	Rate        int    `json:"rate" gorm:"type:int; not null"`
}
