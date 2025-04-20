package models

type Experience struct {
	ExperienceOption string `json:"experience_option" gorm:"not null;unique"`
}
