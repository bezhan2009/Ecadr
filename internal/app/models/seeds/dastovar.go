package seeds

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/logger"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
)

func SeedDastovar(db *gorm.DB) error {
	url := "https://bf06c50f7c3aa09d.mokky.dev/dastovartj"

	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка при отправке запроса: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Ошибка при чтении ответа: %v", err)
	}

	var achievements []models.Achievement
	if err := json.Unmarshal(body, &achievements); err != nil {
		log.Fatalf("Ошибка при декодировании JSON: %v", err)
	}

	IdAch := 1

	for _, achievement := range achievements {
		var existingAchievement models.Achievement
		if err := db.First(&existingAchievement, "id = ?", IdAch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&achievement)
			} else {
				logger.Error.Printf("[seeds.SeedDastovar] Error seeding achievements: %v", err)

				return err
			}
		}

		IdAch++
	}

	return nil
}
