package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model

	Text   string `json:"text" gorm:"type:text; not null"`
	UserID uint   `json:"user_id"`

	User User `json:"user" gorm:"foreignkey:UserID"`
}
