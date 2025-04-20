package ai

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	aiSerivce "Ecadr/internal/app/service/ai"
	"Ecadr/internal/controllers"
	"Ecadr/internal/controllers/middlewares"
	"Ecadr/pkg/errs"
	"errors"
	"github.com/gin-gonic/gin"
)

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
// @Security BearerAuth
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
		courses, err = service.GetAllCourses(search)
		if err != nil {
			controllers.HandleError(c, err)
			return
		}

		c.JSON(200, courses)
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
			c.JSON(200, courses)
			return
		}

		controllers.HandleError(c, err)
		return
	}

	c.JSON(201, analysedCourse)
}
