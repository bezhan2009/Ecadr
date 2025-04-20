package models

import (
	"gorm.io/gorm"
	"time"
)

type SchoolGrade struct {
	gorm.Model
	UserID     int       `json:"user_id" gorm:"not null"`
	User       User      `json:"-" gorm:"foreignKey:UserID"`
	Subject    string    `json:"subject" gorm:"not null"`
	Grade      int8      `json:"grade" gorm:"type:smallint;not null"`
	Class      int8      `json:"class" gorm:"type:smallint;not null"`
	Teacher    string    `json:"teacher" gorm:"not null"`
	LessonDate time.Time `json:"lesson_date" gorm:"not null"`
}
