package repository

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/db"
)

func GetAllUserMessages(userID uint) (msgs []models.Message, err error) {
	if err = db.GetDBConn().Model(&models.Message{}).Find(&msgs).Error; err != nil {
		return nil, err
	}

	return msgs, nil
}

func CreateMessage(msg *models.Message) (err error) {
	if err = db.GetDBConn().Model(&models.Message{}).Create(msg).Error; err != nil {
		return err
	}

	return nil
}
