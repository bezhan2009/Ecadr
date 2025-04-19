package service

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/repository"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"Ecadr/pkg/utils"
	"fmt"
	"strconv"
)

func GetAllUsers() (users []models.User, err error) {
	users, err = repository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByID(id uint) (user models.User, err error) {
	user, err = repository.GetUserByID(strconv.Itoa(int(id)))
	if err != nil {
		return user, err
	}

	return user, nil
}

func CreateUser(user models.User) (uint, error) {
	usernameExists, emailExists, err := repository.UserExists(user.Username, user.Email)
	if err != nil {
		return 0, fmt.Errorf("failed to check existing user: %w", err)
	}

	if user.Password == "" || user.Email == "" || user.Username == "" {
		return 0, errs.ErrInvalidData
	}

	if usernameExists {
		logger.Error.Printf("user with username %s already exists", user.Username)
		return 0, errs.ErrUsernameUniquenessFailed
	}

	if emailExists {
		logger.Error.Printf("user with email %s already exists", user.Email)
		return 0, errs.ErrEmailUniquenessFailed
	}

	user.Password = utils.GenerateHash(user.Password)

	var userDB models.User

	if userDB, err = repository.CreateUser(user); err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userDB.ID, nil
}
