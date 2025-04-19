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

func SeedVacansy(db *gorm.DB) error {
	url := "https://bf06c50f7c3aa09d.mokky.dev/vacancy"

	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка при отправке запроса: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Ошибка при чтении ответа: %v", err)
	}

	var vacancies []models.Vacancy
	if err := json.Unmarshal(body, &vacancies); err != nil {
		log.Fatalf("Ошибка при декодировании JSON: %v", err)
	}

	IdVac := 1

	for _, vacancy := range vacancies {
		var existingVacancy models.Vacancy
		if err := db.First(&existingVacancy, "id = ?", IdVac).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&vacancy)
			} else {
				logger.Error.Printf("[seeds.SeedVacansy] Error seeding vacancies: %v", err)

				return err
			}
		}

		IdVac++
	}

	return nil
}
