package repository

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/db"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"errors"
	"gorm.io/gorm"
)

func GetAllUsers(search string) (users []models.User, err error) {
	query := db.GetDBConn().
		Preload("KindergartenNotes").
		Preload("SchoolGrades").
		Preload("Achievements").
		Preload("TestSubmissions.Answers.Question").
		Joins("LEFT JOIN roles ON roles.id = users.role_id")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(`
			users.username ILIKE ? OR 
			users.email ILIKE ? OR 
			users.first_name ILIKE ? OR 
			users.last_name ILIKE ? OR 
			roles.name ILIKE ?`,
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	err = query.Find(&users).Error
	if err != nil {
		logger.Error.Printf("[repository.GetAllUsers] error getting all users: %s\n", err.Error())
		return nil, TranslateGormError(err)
	}

	return users, nil
}

func GetUserByID(id string) (user models.User, err error) {
	err = db.GetDBConn().
		Preload("KindergartenNotes").
		Preload("SchoolGrades").
		Preload("Achievements").
		Preload("TestSubmissions.Answers.Question").
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		logger.Error.Printf("[repository.GetUserByID] error getting user by id: %v\n", err)
		return user, TranslateGormError(err)
	}
	return user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := db.GetDBConn().Where("username = ?", username).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		logger.Error.Printf("[repository.GetUserByUsername] error getting user by username: %v\n", err)
		return nil, err
	}
	return &user, nil
}

func UserExists(username, email string) (bool, bool, error) {
	users, err := GetAllUsers("")
	if err != nil {
		return false, false, err
	}

	var usernameExists, emailExists bool
	for _, user := range users {
		if user.Username == username {
			usernameExists = true
		}
		if user.Email == email {
			emailExists = true
		}
	}
	return usernameExists, emailExists, nil
}

func CreateUser(user models.User) (userDB models.User, err error) {
	//logger.Debug.Println(user.ID)
	if err = db.GetDBConn().Create(&user).Error; err != nil {
		logger.Error.Printf("[repository.CreateUser] error creating user: %v\n", err)
		return userDB, TranslateGormError(err)
	}

	//logger.Debug.Println(user.ID)
	userDB = user
	return userDB, nil
}

func GetUserByUsernameAndPassword(username string, password string) (user models.User, err error) {
	err = db.GetDBConn().Where("username = ? AND password = ?", username, password).First(&user).Error
	if err != nil {
		logger.Error.Printf("[repository.GetUserByUsernameAndPassword] error getting user by username and password: %v\n", err)
		return user, TranslateGormError(err)
	}

	return user, nil
}

func GetUserByEmailAndPassword(email string, password string) (user models.User, err error) {
	err = db.GetDBConn().Where("email = ? AND password = ?", email, password).First(&user).Error
	if err != nil {
		logger.Error.Printf("[repository.GetUserByEmailAndPassword] error getting user by email and password: %v\n", err)
		return user, TranslateGormError(err)
	}

	return user, nil
}

func GetUserByEmailPasswordAndUsername(username, email, password string) (user models.User, err error) {
	err = db.GetDBConn().Where("email = ? AND password = ? AND username = ?", email, password, username).First(&user).Error
	if err != nil {
		logger.Error.Printf("[repository.GetUserByEmailPasswordAndUsername] error getting user by username, email and password: %v\n", err)
		return user, TranslateGormError(err)
	}

	return user, nil
}
