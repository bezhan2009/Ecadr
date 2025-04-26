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

func SeedCompany(db *gorm.DB) error {
	url := "https://bf06c50f7c3aa09d.mokky.dev/company"

	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка при отправке запроса: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Ошибка при чтении ответа: %v", err)
	}

	var companies []models.Company
	if err := json.Unmarshal(body, &companies); err != nil {
		log.Fatalf("Ошибка при декодировании JSON: %v", err)
	}

	for _, company := range companies {
		var existingCompany models.Company
		if err := db.First(&existingCompany, "title = ?", company.Title).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&models.Company{
					WorkerID:    company.WorkerID,
					Title:       company.Title,
					Description: company.Description,
				})
			} else {
				logger.Error.Printf("[seeds.SeedCompany] Error seeding companies: %v", err)

				return err
			}
		}
	}

	return nil
}
