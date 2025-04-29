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

func SeedCriteria(db *gorm.DB) error {
	url := "https://bf06c50f7c3aa09d.mokky.dev/vacancy-criteria"

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

	var criteria []models.Criteria
	if err := json.Unmarshal(body, &criteria); err != nil {
		log.Fatalf("Ошибка при декодировании JSON: %v", err)
	}

	IdCriteria := 1

	for _, crit := range criteria {
		var existingGrade models.Criteria
		if err := db.First(&existingGrade, "id = ?", IdCriteria).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&crit)
			} else {
				logger.Error.Printf("[seeds.SeedGrades] Error seeding criteria: %v", err)

				return err
			}
		}

		IdCriteria++
	}

	return nil
}
