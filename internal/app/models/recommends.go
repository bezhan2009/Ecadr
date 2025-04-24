package models

type Recommend struct {
	ID         uint   `json:"id" gorm:"primary_key;auto_increment"`
	UserID     uint   `json:"user_id" gorm:"not null"`
	TargetID   uint   `json:"target_id" gorm:"not null"`
	TargetType string `json:"target_type" gorm:"not null"`

	User User `json:"-" gorm:"foreignkey:UserID"`
}
