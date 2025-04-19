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

func SeedCourse(db *gorm.DB) error {
	url := "https://bf06c50f7c3aa09d.mokky.dev/courses"

	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка при отправке запроса: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Ошибка при чтении ответа: %v", err)
	}

	var courses []models.Course
	if err := json.Unmarshal(body, &courses); err != nil {
		log.Fatalf("Ошибка при декодировании JSON: %v", err)
	}

	IdCourse := 1

	for _, vacancy := range courses {
		var existingCourse models.Course
		if err := db.First(&existingCourse, "id = ?", IdCourse).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&vacancy)
			} else {
				logger.Error.Printf("[seeds.SeedCourse] Error seeding courses: %v", err)

				return err
			}
		}

		IdCourse++
	}

	return nil
}
