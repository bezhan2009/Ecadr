package service

import (
	"Ecadr/internal/app/models"
	redisService "Ecadr/internal/app/service/redis"
	"Ecadr/internal/app/service/validators"
	"Ecadr/internal/repository"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"Ecadr/pkg/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

func GetAllUsers(search string) (users []models.User, err error) {
	users, err = repository.GetAllUsers(search)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByID(id uint) (user models.User, charter models.Charter, err error) {
	user, err = redisService.GetCachedUser(id)
	if err != nil {
		user, err = repository.GetUserByID(strconv.Itoa(int(id)))
		if err != nil {
			return user, charter, err
		}
	}

	charter, err = redisService.GetCachedCharter(id)
	resp := ""
	if err != nil {
		resp, _ = utils.SendTextToGemini(fmt.Sprintf(`Скинь вот такой обьект json чартеров 
{
	"subject": "Предмет из школы",
	"score": "Статистика этого предмета основопологаясь на данных пользователя(от 1 до 100)"
}
вот данные пользователя
%v
`, user))

		err = json.Unmarshal([]byte(resp), &charter)
		if err != nil {
			logger.Error.Println(err)
		}
	}

	redisService.SetUserCache(user)
	redisService.SetCharterCache(charter, user.ID)

	return user, charter, nil
}

func CreateUser(user models.User, company models.Company) (uint, error) {
	usernameExists, emailExists, err := repository.UserExists(user.Username, user.Email)
	if err != nil {
		return 0, fmt.Errorf("failed to check existing user: %w", err)
	}

	if user.Password == "" || user.Email == "" || user.Username == "" {
		return 0, errs.ErrInvalidData
	}

	if usernameExists {
		logger.Error.Printf("[service.CreateUser] user with username %s already exists", user.Username)

		return 0, errs.ErrUsernameUniquenessFailed
	}

	if emailExists {
		logger.Error.Printf("[service.CreateUser] user with email %s already exists", user.Email)
		return 0, errs.ErrEmailUniquenessFailed
	}

	if user.RoleID == 3 {
		if err = validators.ValidateCompany(company); err != nil {
			return 0, err
		}
	}

	user.Password = utils.GenerateHash(user.Password)

	var userDB models.User

	if userDB, err = repository.CreateUser(user, company); err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	redisService.SetUserCache(user)

	return userDB.ID, nil
}
