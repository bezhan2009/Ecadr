package models

type Recommend struct {
	ID        uint `json:"id" gorm:"primary_key;auto_increment"`
	UserID    uint `json:"userID" gorm:"not null"`
	VacancyID uint `json:"vacancyID"`
	CourseID  uint `json:"courseID"`

	User    User    `json:"-" gorm:"foreignkey:UserID"`
	Vacancy Vacancy `json:"-" gorm:"foreignkey:VacancyID"`
	Course  Course  `json:"-" gorm:"foreignkey:CourseID"`
}
