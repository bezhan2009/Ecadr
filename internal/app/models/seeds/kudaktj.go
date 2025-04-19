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

func SeedKinderGarten(db *gorm.DB) error {
	url := "https://bf06c50f7c3aa09d.mokky.dev/kudaktj"

	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка при отправке запроса: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Ошибка при чтении ответа: %v", err)
	}

	var kinderGardens []models.KindergartenNoteRequest
	if err := json.Unmarshal(body, &kinderGardens); err != nil {
		log.Fatalf("Ошибка при декодировании JSON: %v", err)
	}

	for _, kinderGarden := range kinderGardens {
		var existingGarten models.KindergartenNote
		if err := db.First(&existingGarten, "note = ?", kinderGarden.Note).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&models.KindergartenNote{
					UserID:   kinderGarden.UserID,
					Note:     kinderGarden.Note,
					Author:   kinderGarden.Author,
					NoteDate: kinderGarden.NoteDate,
				})
			} else {
				logger.Error.Printf("[seeds.SeedKinderGarten] Error seeding kinderGardens: %v", err)

				return err
			}
		}
	}

	return nil
}
