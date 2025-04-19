package service

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/repository"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/utils"
	"errors"
)

func SignIn(userDataCheck, password string) (user models.User, accessToken string, refreshToken string, err error) {
	if userDataCheck == "" {
		return user, "", "", errs.ErrInvalidData
	}

	user, err = repository.GetUserByEmailAndPassword(userDataCheck, password)
	if err != nil {
		if !errors.Is(err, errs.ErrRecordNotFound) {
			return user, "", "", err
		}

		user, err = repository.GetUserByUsernameAndPassword(userDataCheck, password)
		if err != nil {
			if !errors.Is(err, errs.ErrRecordNotFound) {
				return user, "", "", err
			}

			return user, "", "", errs.ErrInvalidCredentials
		}
	}

	accessToken, refreshToken, err = utils.GenerateToken(user.ID, uint(user.RoleID), user.Username)
	if err != nil {
		return user, "", "", err
	}

	return user, accessToken, refreshToken, nil
}
