package models

import (
	"gorm.io/gorm"
	"time"
)

type Achievement struct {
	gorm.Model

	UserID          uint      `json:"user_id" gorm:"not null"`
	User            User      `json:"-" gorm:"foreignKey:UserID"`
	Title           string    `json:"title"`
	Description     string    `json:"description" gorm:"type:text"`
	Organization    string    `json:"organization"`
	AchievementDate time.Time `json:"achievement_date" gorm:"type:date;not null"`
}
