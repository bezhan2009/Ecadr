package ai

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	aiSerivce "Ecadr/internal/app/service/ai"
	"Ecadr/internal/controllers"
	"Ecadr/internal/security"
	"Ecadr/pkg/db"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func getUserCourseFromRedis(key string) ([]models.User, error) {
	UserStr, err := db.GetCache(key)
	if UserStr != "" {
		var usersJson []models.User
		err := json.Unmarshal([]byte(UserStr), &usersJson)
		if err != nil {
			db.DeleteCache(key)
		} else {
			return usersJson, nil
		}
	}

	return nil, err
}

// GetAnalyseForCourseUser godoc
// @Summary Анализ для курсов чтобы искать пользователей
// @Description Возвращает список рекомендованных пользователей для курса
// @Tags AI
// @Accept  json
// @Produce  json
// @Param search query string false "Поисковый запрос по пользователям"
// @Success 201 {object} []models.UserRequest "Рекомендованные пользователи на основе анализа AI"
// @Failure 400 {object} errs.ErrorResp "Неверный запрос или пользователь не найден"
// @Failure 500 {object} errs.ErrorResp "Внутренняя ошибка сервера"
// @Router /course/users/{id} [get]
// @Security ApiKeyAuth
func GetAnalyseForCourseUser(c *gin.Context) {
	search := c.Query("search")
	courseIDStr := c.Param("id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		controllers.HandleError(c, errs.ErrInvalidID)
		return
	}

	if search != "" {
		keyCacheRedisSearch := fmt.Sprintf("searched_users_%s", search)

		searchVacancy, err := getUserVacancyFromRedis(keyCacheRedisSearch)
		if err == nil && len(searchVacancy) > 0 {
			c.JSON(200, searchVacancy)
			return
		}

		users, err := service.GetAllUsers(search)
		if err != nil {
			controllers.HandleError(c, err)
			return
		}

		usersJson, err := json.Marshal(users)
		if err != nil {
			logger.Error.Printf("[ai.GetAnalyseForVacanciesUser] Error marshalling users json: %v", err)
		} else {
			db.SetCache(
				keyCacheRedisSearch,
				usersJson,
				time.Duration(security.AppSettings.RedisParams.TTLMinutes)*time.Minute,
			)
		}

		c.JSON(200, users)
		return
	}

	keyCacheRedis := fmt.Sprintf("analyzed_users_course_%d", courseID)

	analysedUsersCache, err := getUserCourseFromRedis(keyCacheRedis)
	if err == nil && len(analysedUsersCache) > 0 {
		c.JSON(200, analysedUsersCache)
		return
	}

	users, err := service.GetAllUsers(search)
	if err != nil {
		controllers.HandleError(c, err)
		return
	}

	course, err := service.GetCourseById(courseID)
	if err != nil {
		controllers.HandleError(c, err)
		return
	}

	analysedUsers, err := aiSerivce.GetAnalyseForCourseUser(
		course,
		users,
	)
	if err != nil {
		if errors.Is(err, errs.ErrNoCourseFound) {
			usersJson, err := json.Marshal(users)
			if err != nil {
				logger.Error.Printf("[ai.GetAnalyseForVacanciesUser] Error marshalling analysed users json: %v", err)
			} else {
				db.SetCache(
					keyCacheRedis,
					usersJson,
					time.Duration(security.AppSettings.RedisParams.TTLMinutes)*time.Minute,
				)
			}

			c.JSON(200, users)
			return
		}

		controllers.HandleError(c, err)
		return
	}

	userJson, err := json.Marshal(analysedUsers)
	if err != nil {
		logger.Error.Printf("[ai.GetAnalyseForVacanciesUser] Error marshalling analysed vacancies json: %v", err)
	} else {
		db.SetCache(
			keyCacheRedis,
			userJson,
			time.Duration(security.AppSettings.RedisParams.TTLMinutes)*time.Minute,
		)
	}

	c.JSON(201, analysedUsers)
}
