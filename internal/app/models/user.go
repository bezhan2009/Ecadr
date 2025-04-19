package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username  string `json:"username" gorm:"unique;not null"`
	Email     string `json:"email" gorm:"unique;not null"`
	Password  string `json:"-" gorm:"size:255;not null"`
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`

	RoleID int  `json:"role_id" gorm:"not null"`
	Role   Role `json:"-" gorm:"foreignKey:RoleID;not null"`

	KindergartenNotes []KindergartenNote `gorm:"foreignKey:UserID"`
	SchoolGrades      []SchoolGrade      `gorm:"foreignKey:UserID"`
	Achievements      []Achievement      `gorm:"foreignKey:UserID"`
	TestSubmissions   []TestSubmission   `gorm:"foreignKey:UserID"`
	//BirthDate time.Time
	//PassportNumber string
	//Phone          string
}

type Role struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}
