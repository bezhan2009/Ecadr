package service

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/repository"
	"Ecadr/pkg/errs"
)

func CreateMessage(msg models.Message) (err error) {
	if msg.Text == "" {
		return errs.ErrEmptyMessageText
	}

	err = repository.CreateMessage(&msg)
	if err != nil {
		return err
	}

	return nil
}
