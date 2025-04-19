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

func SeedGrades(db *gorm.DB) error {
	url := "https://bf06c50f7c3aa09d.mokky.dev/edonish"

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

	var grades []models.SchoolGrade
	if err := json.Unmarshal(body, &grades); err != nil {
		log.Fatalf("Ошибка при декодировании JSON: %v", err)
	}

	IdGrades := 1

	for _, grade := range grades {
		var existingGrade models.SchoolGrade
		if err := db.First(&existingGrade, "id = ?", IdGrades).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&grade)
			} else {
				logger.Error.Printf("[seeds.SeedGrades] Error seeding grades: %v", err)

				return err
			}
		}

		IdGrades++
	}

	return nil
}
