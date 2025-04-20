package ai

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	aiSerivce "Ecadr/internal/app/service/ai"
	"Ecadr/internal/controllers"
	"Ecadr/internal/controllers/middlewares"
	"Ecadr/internal/security"
	"Ecadr/pkg/db"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func getCourseFromRedis(key string) ([]models.Course, error) {
	CourseStr, err := db.GetCache(key)
	if CourseStr != "" {
		var courseJson []models.Course
		err := json.Unmarshal([]byte(CourseStr), &courseJson)
		if err != nil {
			db.DeleteCache(key)
		} else {
			return courseJson, nil
		}
	}

	return nil, err
}

// GetAnalyseForUserCourse godoc
// @Summary Анализ курсов для пользователя
// @Description Возвращает список рекомендованных курсов на основе данных пользователя (оценки, достижения и т.д.)
// @Tags AI
// @Accept  json
// @Produce  json
// @Param search query string false "Поисковый запрос по курсам"
// @Success 200 {array} models.CourseResponse "Успешный ответ с курсами (если нет подходящих)"
// @Success 201 {object} []models.CourseResponse "Рекомендованные курсы на основе анализа AI"
// @Failure 400 {object} errs.ErrorResp "Неверный запрос или пользователь не найден"
// @Failure 500 {object} errs.ErrorResp "Внутренняя ошибка сервера"
// @Router /course [get]
// @Security ApiKeyAuth
func GetAnalyseForUserCourse(c *gin.Context) {
	search := c.Query("search")
	userID := c.GetUint(middlewares.UserIDCtx)

	if userID == 0 {
		controllers.HandleError(c, errs.ErrUserNotFound)
		return
	}

	userData, err := service.GetUserByID(userID)
	if err != nil {
		controllers.HandleError(c, err)
		return
	}

	var courses []models.Course
	if search != "" {
		keyCacheRedisSearch := fmt.Sprintf("searched_course_%s", search)

		searchCourse, err := getCourseFromRedis(keyCacheRedisSearch)
		if err == nil && len(searchCourse) > 0 {
			c.JSON(200, searchCourse)
			return
		}

		courses, err = service.GetAllCourses(search)
		if err != nil {
			controllers.HandleError(c, err)
			return
		}

		courseJson, err := json.Marshal(courses)
		if err != nil {
			logger.Error.Printf("[ai.GetAnalyseForUserCourse] Error marshalling courses json: %v", err)
		} else {
			db.SetCache(
				keyCacheRedisSearch,
				courseJson,
				time.Duration(security.AppSettings.RedisParams.TTLMinutes)*time.Minute,
			)
		}
		c.JSON(200, courses)
		return
	}

	keyCacheRedis := fmt.Sprintf("analyzed_course_%d", userID)

	analysedCoursesCache, err := getCourseFromRedis(keyCacheRedis)
	if err == nil && len(analysedCoursesCache) > 0 {
		c.JSON(200, analysedCoursesCache)
		return
	}

	courses, err = service.GetAllCourses(search)
	if err != nil {
		controllers.HandleError(c, err)
		return
	}

	analysedCourse, err := aiSerivce.GetAnalyseForUserCourse(
		courses,
		userData.KindergartenNotes,
		userData.SchoolGrades,
		userData.Achievements,
	)
	if err != nil {
		if errors.Is(err, errs.ErrNoCourseFound) {
			courseJson, err := json.Marshal(courses)
			if err != nil {
				logger.Error.Printf("[ai.GetAnalyseForUserCourse] Error marshalling analysed courses json: %v", err)
			} else {
				db.SetCache(
					keyCacheRedis,
					courseJson,
					time.Duration(security.AppSettings.RedisParams.TTLMinutes)*time.Minute,
				)
			}

			c.JSON(200, courses)
			return
		}

		controllers.HandleError(c, err)
		return
	}

	courseJson, err := json.Marshal(analysedCourse)
	if err != nil {
		logger.Error.Printf("[ai.GetAnalyseForUserCourse] Error marshalling analysed courses json: %v", err)
	} else {
		db.SetCache(
			keyCacheRedis,
			courseJson,
			time.Duration(security.AppSettings.RedisParams.TTLMinutes)*time.Minute,
		)
	}

	c.JSON(201, analysedCourse)
}
