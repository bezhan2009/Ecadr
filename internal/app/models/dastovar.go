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
	Place           int       `json:"place" gorm:"not null"`
	Level           string    `json:"level" gorm:"not null"`
	Location        string    `json:"location" gorm:"not null"`
	Type            string    `json:"type" gorm:"not null"`
	Organization    string    `json:"organization" gorm:"not null"`
	AchievementDate time.Time `json:"achievement_date" gorm:"type:date;not null"`
}
