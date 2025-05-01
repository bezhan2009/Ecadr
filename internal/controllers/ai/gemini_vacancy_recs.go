package ai

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	aiSerivce "Ecadr/internal/app/service/ai"
	"Ecadr/internal/controllers/middlewares"
	"Ecadr/internal/security"
	"Ecadr/pkg/db"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func getVacancyFromRedis(key string) ([]models.Vacancy, error) {
	VacanciesStr, err := db.GetCache(key)
	if VacanciesStr != "" {
		var vacanciesJson []models.Vacancy
		err := json.Unmarshal([]byte(VacanciesStr), &vacanciesJson)
		if err != nil {
			db.DeleteCache(key)
		} else {
			return vacanciesJson, nil
		}
	}

	return nil, err
}

// GetAnalyseForUserVacancies godoc
// @Summary Анализ вакансий для пользователя
// @Description Возвращает список рекомендованных вакансий на основе данных пользователя (оценки, достижения и т.д.)
// @Tags AI
// @Accept  json
// @Produce  json
// @Param search query string false "Поисковый запрос по вакансиям"
// @Success 200 {array} models.VacancyResponse "Успешный ответ с вакансиям (если нет подходящих)"
// @Success 201 {object} []models.VacancyResponse "Рекомендованные вакансий на основе анализа AI"
// @Failure 400 {object} errs.ErrorResp "Неверный запрос или пользователь не найден"
// @Failure 500 {object} errs.ErrorResp "Внутренняя ошибка сервера"
// @Router /vacancy [get]
// @Security ApiKeyAuth
func GetAnalyseForUserVacancies(c *gin.Context) (interface{}, error) {
	search := c.Query("search")
	userID := c.GetUint(middlewares.UserIDCtx)

	if userID == 0 {
		//controllers.HandleError(c, errs.ErrUserNotFound)
		return nil, errs.ErrUserNotFound
	}

	userData, _, err := service.GetUserByID(userID)
	if err != nil {
		//controllers.HandleError(c, err)
		return nil, err
	}

	var vacancies []models.Vacancy
	if search != "" {
		keyCacheRedisSearch := fmt.Sprintf("searched_vacancy_%s", search)

		searchVacancy, err := getVacancyFromRedis(keyCacheRedisSearch)
		if err == nil && len(searchVacancy) > 0 {
			//c.JSON(http.StatusOK, searchVacancy)
			return searchVacancy, nil
		}

		vacancies, err = service.GetAllVacancies(search)
		if err != nil {
			//controllers.HandleError(c, err)
			return nil, err
		}

		vacancyJson, err := json.Marshal(vacancies)
		if err != nil {
			logger.Error.Printf("[ai.GetAnalyseForUserVacancies] Error marshalling vacancies json: %v", err)
		} else {
			db.SetCache(
				keyCacheRedisSearch,
				vacancyJson,
				time.Duration(security.AppSettings.RedisParams.TTLMinutes)*time.Minute,
			)
		}

		//c.JSON(http.StatusOK, vacancies)
		return vacancies, nil
	}

	keyCacheRedis := fmt.Sprintf("analyzed_vacancies_%d", userID)

	analysedVacanciesCache, err := getVacancyFromRedis(keyCacheRedis)
	if err == nil && len(analysedVacanciesCache) > 0 {
		//c.JSON(200, analysedVacanciesCache)
		return analysedVacanciesCache, err
	}

	vacancies, err = service.GetAllVacancies(search)
	if err != nil {
		//controllers.HandleError(c, err)
		return nil, err
	}

	analysedVacancies, err := aiSerivce.GetAnalyseForUserVacancies(
		vacancies,
		userData.KindergartenNotes,
		userData.SchoolGrades,
		userData.Achievements,
	)
	if err != nil {
		if errors.Is(err, errs.ErrNoVacancyFound) {
			vacancyJson, err := json.Marshal(vacancies)
			if err != nil {
				logger.Error.Printf("[ai.GetAnalyseForUserVacancies] Error marshalling analysed vacancies json: %v", err)
			} else {
				db.SetCache(
					keyCacheRedis,
					vacancyJson,
					time.Duration(security.AppSettings.RedisParams.TTLMinutes)*time.Minute,
				)
			}

			c.JSON(http.StatusOK, vacancies)
			return vacancies, err
		}

		//controllers.HandleError(c, err)
		return nil, err
	}

	vacancyJson, err := json.Marshal(analysedVacancies)
	if err != nil {
		logger.Error.Printf("[ai.GetAnalyseForUserVacancies] Error marshalling analysed vacancies json: %v", err)
	} else {
		db.SetCache(
			keyCacheRedis,
			vacancyJson,
			time.Duration(security.AppSettings.RedisParams.TTLMinutes)*time.Minute,
		)
	}

	//c.JSON(http.StatusCreated, analysedVacancies)
	return analysedVacancies, nil
}
