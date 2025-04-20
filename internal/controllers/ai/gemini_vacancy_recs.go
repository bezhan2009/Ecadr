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
func GetAnalyseForUserVacancies(c *gin.Context) {
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

	var vacancies []models.Vacancy
	if search != "" {
		vacancies, err = service.GetAllVacancies(search)
		if err != nil {
			controllers.HandleError(c, err)
			return
		}

		c.JSON(200, vacancies)
		return
	}

	vacancies, err = service.GetAllVacancies(search)
	if err != nil {
		controllers.HandleError(c, err)
		return
	}

	analysedVacancies, err := aiSerivce.GetAnalyseForUserVacancies(
		vacancies,
		userData.KindergartenNotes,
		userData.SchoolGrades,
		userData.Achievements,
	)
	if err != nil {
		if errors.Is(err, errs.ErrNoVacancyFound) {
			c.JSON(200, vacancies)
			return
		}

		controllers.HandleError(c, err)
		return
	}

	c.JSON(201, analysedVacancies)
}
