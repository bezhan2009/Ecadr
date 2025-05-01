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

func SeedTestSubmissions(db *gorm.DB) error {
	url := "https://bf06c50f7c3aa09d.mokky.dev/test_submissions"

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

	var tests []models.TestSubmission
	if err := json.Unmarshal(body, &tests); err != nil {
		log.Fatalf("Ошибка при декодировании JSON: %v", err)
	}

	IdTest := 1

	for _, test := range tests {
		var existingTest models.TestSubmission
		if err := db.First(&existingTest, "id = ?", IdTest).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&test)
			} else {
				logger.Error.Printf("[seeds.SeedTestSubmissions] Error seeding tests: %v", err)

				return err
			}
		}

		IdTest++
	}

	return nil
}
