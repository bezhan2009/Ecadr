package models

import (
	"gorm.io/gorm"
	"time"
)

type KindergartenNote struct {
	gorm.Model

	UserID   int       `json:"user_id" gorm:"not null"`
	User     User      `json:"-" gorm:"foreignKey:UserID"`
	Note     string    `json:"note" gorm:"type:text"`
	Author   string    `json:"author" gorm:"type:text"`
	NoteDate time.Time `json:"note_date"`
}

type KindergartenNoteRequest struct {
	UserID   int       `json:"user_id" gorm:"not null"`
	User     User      `json:"-" gorm:"foreignKey:UserID"`
	Note     string    `json:"note" gorm:"type:text"`
	Author   string    `json:"author" gorm:"type:text"`
	NoteDate time.Time `json:"note_date"`
}
