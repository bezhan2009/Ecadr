package controllers

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	"Ecadr/internal/controllers/middlewares"
	"Ecadr/pkg/errs"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetUserRecommendsCourse godoc
// @Summary Получить рекомендованные курсы для пользователя
// @Description Возвращает список курсов, рекомендованных пользователю
// @Tags Recommendations
// @Produce json
// @Success 200 {array} models.CourseResponse "Список рекомендованных курсов"
// @Failure 400 {object} errs.ErrorResp
// @Failure 500 {object} errs.ErrorResp
// @Router /course/recommends [get]
// @Security ApiKeyAuth
func GetUserRecommendsCourse(c *gin.Context) {
	userID := c.GetUint(middlewares.UserIDCtx)

	recommends, err := service.GetUserRecommendsCourse(userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, recommends)
}

// GetUserRecommendsVacancy godoc
// @Summary Получить рекомендованные вакансии для пользователя
// @Description Возвращает список вакансий, рекомендованных пользователю
// @Tags Recommendations
// @Produce json
// @Success 200 {array} models.VacancyResponse "Список рекомендованных вакансий"
// @Failure 400 {object} errs.ErrorResp
// @Failure 500 {object} errs.ErrorResp
// @Router /recommends/vacancy [get]
// @Security ApiKeyAuth
func GetUserRecommendsVacancy(c *gin.Context) {
	userID := c.GetUint(middlewares.UserIDCtx)

	recommends, err := service.GetUserRecommendsVacancy(userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, recommends)
}

// GetUserRecommendByID godoc
// @Summary Получить конкретную рекомендацию по ID
// @Description Возвращает одну рекомендацию пользователя по её ID
// @Tags Recommendations
// @Produce json
// @Param id path int true "ID рекомендации"
// @Success 200 {object} models.Recommend "Рекомендация пользователя"
// @Failure 400 {object} errs.ErrorResp
// @Failure 404 {object} errs.ErrorResp
// @Failure 500 {object} errs.ErrorResp
// @Router /recommends/{id} [get]
// @Security ApiKeyAuth
func GetUserRecommendByID(c *gin.Context) {
	userID := c.GetUint(middlewares.UserIDCtx)
	recommendStrID := c.Param("id")
	recommendID, err := strconv.Atoi(recommendStrID)
	if err != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	recommend, err := service.GetUserRecommendByID(userID, uint(recommendID))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, recommend)
}

// CreateRecommendCourse godoc
// @Summary Создать рекомендацию курса
// @Description Создаёт новую рекомендацию курса для пользователя
// @Tags Recommendations
// @Accept json
// @Produce json
// @Param recommend body models.Recommend true "Данные рекомендации (только поле CourseID должно быть заполнено)"
// @Success 201 {object} models.RecommendReq "Успешное сообщение"
// @Failure 400 {object} errs.ErrorResp
// @Failure 500 {object} errs.ErrorResp
// @Router /recommends/course [post]
// @Security ApiKeyAuth
func CreateRecommendCourse(c *gin.Context) {
	userID := c.GetUint(middlewares.UserIDCtx)

	var recommend models.Recommend
	err := c.BindJSON(&recommend)
	if err != nil {
		HandleError(c, err)
		return
	}

	recommend.UserID = userID
	recommend.VacancyID = 0

	err = service.CreateRecommend(recommend)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, gin.H{"message": "recommend course created successfully"})
}

// CreateRecommendVacancy godoc
// @Summary Создать рекомендацию вакансии
// @Description Создаёт новую рекомендацию вакансии для пользователя
// @Tags Recommendations
// @Accept json
// @Produce json
// @Param recommend body models.Recommend true "Данные рекомендации (только поле VacancyID должно быть заполнено)"
// @Success 201 {object} models.RecommendReq "Успешное сообщение"
// @Failure 400 {object} errs.ErrorResp
// @Failure 500 {object} errs.ErrorResp
// @Router /recommends/vacancy [post]
// @Security ApiKeyAuth
func CreateRecommendVacancy(c *gin.Context) {
	userID := c.GetUint(middlewares.UserIDCtx)

	var recommend models.Recommend
	err := c.BindJSON(&recommend)
	if err != nil {
		HandleError(c, err)
		return
	}

	recommend.UserID = userID
	recommend.CourseID = 0

	err = service.CreateRecommend(recommend)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, gin.H{"message": "recommend vacancy created successfully"})
}
