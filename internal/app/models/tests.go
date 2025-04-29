package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type Test struct {
	ID          int    `gorm:"primaryKey"`
	Title       string `json:"title"`
	Description string `json:"description" gorm:"type:text"`
	TargetType  int    `json:"target_type"` // "vacancy" или "course"
	TargetID    int    `json:"target_id"`   // ID вакансии или курса
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	Questions []Question `gorm:"foreignKey:TestID"`
}

type TestSubmission struct {
	ID          int       `gorm:"primaryKey"`
	TestID      int       `json:"test_id" gorm:"not null"`
	Test        Test      `gorm:"foreignKey:TestID"` // ← добавлено поле для связи
	UserID      int       `json:"user_id" gorm:"not null"`
	SubmittedAt time.Time `json:"submitted_at"`
	Answers     []Answer  `gorm:"foreignKey:SubmissionID"`
	// не забудьте добавить сюда при необходимости:
	CorrectPercentage float64 `json:"correct_percentage" gorm:"-"`
}

type Answer struct {
	ID                int           `gorm:"primaryKey"`
	SubmissionID      int           `json:"submission_id" gorm:"not null"`
	QuestionID        int           `json:"question_id" gorm:"not null"`
	Question          Question      `gorm:"foreignKey:QuestionID"` // ← связь
	TextAnswer        string        `json:"text_answer"`
	SelectedChoiceIDs pq.Int64Array `json:"selected_choice_ids" gorm:"type:int[]"`
}

type Question struct {
	ID        int      `gorm:"primaryKey"`
	TestID    int      `json:"test_id" gorm:"not null"`
	Test      Test     `json:"-" gorm:"foreignKey:TestID"`
	Text      string   `json:"text" gorm:"type:text"`
	Type      string   `json:"type"`
	Choices   []Choice `gorm:"foreignKey:QuestionID"` // ← связь
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Choice struct {
	ID         int      `gorm:"primaryKey"`
	QuestionID int      `json:"question_id" gorm:"not null"`
	Question   Question `json:"-" gorm:"foreignKey:QuestionID"`
	Text       string   `json:"text"`
	IsCorrect  bool     `json:"is_correct"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
