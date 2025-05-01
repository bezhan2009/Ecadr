package ai

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	aiService "Ecadr/internal/app/service/ai"
	"Ecadr/internal/security"
	"Ecadr/pkg/db"
	"Ecadr/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func getAnalysedDataCompaniesFromRedis(key string) ([]models.CompanyStatistic, error) {
	analysedDataStr, err := db.GetCache(key)
	if analysedDataStr != "" {
		var analysedDataJson []models.CompanyStatistic
		err = json.Unmarshal([]byte(analysedDataStr), &analysedDataJson)
		if err != nil {
			db.DeleteCache(key)
		} else {
			return analysedDataJson, nil
		}
	}

	return nil, err
}

// GetCompaniesStatistic godoc
// @Summary Статистика компаний
// @Description Возвращает статистику компаний учитывая все их продукты(вакансий, курсы и тд)
// @Tags AI
// @Accept  json
// @Produce  json
// @Failure 400 {object} errs.ErrorResp "Неверный запрос"
// @Failure 500 {object} errs.ErrorResp "Внутренняя ошибка сервера"
// @Router /analyse/companies [get]
// @Security ApiKeyAuth
func GetCompaniesStatistic(c *gin.Context) (interface{}, error) {
	fmt.Println("GetCompaniesStatistic")
	companiesInfo, err := service.GetCompaniesProducts()
	if err != nil {
		//controllers.HandleError(c, err)
		return nil, err
	}

	keyCacheRedis := "analyzed_companies"

	analysedCompaniesCache, err := getAnalysedDataCompaniesFromRedis(keyCacheRedis)
	if err == nil && len(analysedCompaniesCache) > 0 {
		//c.JSON(http.StatusOK, analysedCompaniesCache)
		return analysedCompaniesCache, nil
	}

	analysedData, err := aiService.GetCompaniesStatistic(
		companiesInfo,
	)
	if err != nil {
		//controllers.HandleError(c, err)
		return nil, err
	}

	analysedDataJson, err := json.Marshal(analysedData)
	if err != nil {
		logger.Error.Printf("[ai.GetCompaniesStatistic] Error marshalling analysed vacancies json: %v", err)
	} else {
		db.SetCache(
			keyCacheRedis,
			analysedDataJson,
			time.Duration(security.AppSettings.RedisParams.TTLMinutes)*time.Minute,
		)
	}

	//c.JSON(http.StatusCreated, analysedData)
	return analysedData, nil
}
