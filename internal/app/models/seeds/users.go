package seeds

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/logger"
	"Ecadr/pkg/utils"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
)

func SeedUsers(db *gorm.DB) error {
	url := "https://bf06c50f7c3aa09d.mokky.dev/users"

	// Отправьте GET-запрос
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка при отправке запроса: %v", err)
	}
	defer response.Body.Close()

	// Прочитайте тело ответа
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Ошибка при чтении ответа: %v", err)
	}

	var users []models.UserRequest
	if err := json.Unmarshal(body, &users); err != nil {
		log.Fatalf("Ошибка при декодировании JSON: %v", err)
	}

	for _, user := range users {
		user.Password = "bezhan2009"
		var existingUser models.User
		if err := db.First(&existingUser, "username = ?", user.Username).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				user.Password = utils.GenerateHash(user.Password)

				db.Create(&models.User{
					Username:  user.Username,
					Email:     user.Email,
					Password:  user.Password,
					FirstName: user.FirstName,
					LastName:  user.LastName,
					RoleID:    int(user.RoleID),
				})
			} else {
				logger.Error.Printf("[seeds.SeedUsers] Error seeding users: %v", err)

				return err
			}
		}
	}

	return nil
}
