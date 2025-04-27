package ai

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	aiService "Ecadr/internal/app/service/ai"
	"Ecadr/internal/controllers"
	"Ecadr/internal/security"
	"Ecadr/pkg/db"
	"Ecadr/pkg/logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func getAnalysedDataUsersFromRedis(key string) ([]models.UsersStatistic, error) {
	analysedDataStr, err := db.GetCache(key)
	if analysedDataStr != "" {
		var analysedDataJson []models.UsersStatistic
		err = json.Unmarshal([]byte(analysedDataStr), &analysedDataJson)
		if err != nil {
			db.DeleteCache(key)
		} else {
			return analysedDataJson, nil
		}
	}

	return nil, err
}

// GetUsersStatistic godoc
// @Summary Статистика пользователей
// @Description Возвращает статистику пользователей
// @Tags AI
// @Accept  json
// @Produce  json
// @Success 201 {object} []models.UsersStatistic "Сама статистика от сервиса и ИИ"
// @Failure 400 {object} errs.ErrorResp "Неверный запрос"
// @Failure 500 {object} errs.ErrorResp "Внутренняя ошибка сервера"
// @Router /analyse/users [get]
// @Security ApiKeyAuth
func GetUsersStatistic(c *gin.Context) {
	usersInfo, err := service.GetAllUsers("")
	if err != nil {
		controllers.HandleError(c, err)
		return
	}

	keyCacheRedis := "analyzed_users"

	analysedUsersCache, err := getAnalysedDataUsersFromRedis(keyCacheRedis)
	if err == nil && len(analysedUsersCache) > 0 {
		c.JSON(http.StatusOK, analysedUsersCache)
		return
	}

	analysedData, err := aiService.GetUsersStatistic(
		usersInfo,
	)
	if err != nil {
		controllers.HandleError(c, err)
		return
	}

	analysedDataJson, err := json.Marshal(analysedData)
	if err != nil {
		logger.Error.Printf("[ai.GetUsersStatistic] Error marshalling analysed data json: %v", err)
	} else {
		db.SetCache(
			keyCacheRedis,
			analysedDataJson,
			time.Duration(security.AppSettings.RedisParams.TTLMinutes)*time.Minute,
		)
	}

	c.JSON(http.StatusCreated, analysedData)
}
