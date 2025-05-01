package controllers

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	"Ecadr/internal/app/service/validators"
	"Ecadr/internal/controllers/middlewares"
	"Ecadr/pkg/errs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetAllCompanyVacancies godoc
// @Summary Получить все воркера вакансии
// @Description Возвращает список всех вакансий воркера с возможностью поиска
// @Tags Vacancies
// @Produce json
// @Param id path int true "workerID вакансии"
// @Param search query string false "Поисковый запрос"
// @Success 200 {object} models.VacancyResponse "Вакансия"
// @Failure 500 {object} errs.ErrorResp
// @Router /vacancy/company/{company_id} [get]
func GetAllCompanyVacancies(c *gin.Context) {
	search := c.Query("search")
	companyStrID := c.Param("company_id")
	companyID, err := strconv.Atoi(companyStrID)
	if err != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	company, err := service.GetCompanyByID(uint(companyID))
	if err != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	vacancies, err := service.GetAllWorkerVacancies(company.WorkerID, search)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, vacancies)
}

// GetVacancyByID godoc
// @Summary Получить все вакансии
// @Description Возвращает список всех вакансий с возможностью поиска
// @Tags Vacancies
// @Produce json
// @Param id path int true "ID вакансии"
// @Param search query string false "Поисковый запрос"
// @Success 200 {object} models.VacancyResponse "Вакансия"
// @Failure 500 {object} errs.ErrorResp
// @Router /vacancy/{id} [get]
func GetVacancyByID(c *gin.Context) {
	vacancyIDStr := c.Param("id")
	vacancyID, err := strconv.Atoi(vacancyIDStr)
	if err != nil {
		HandleError(c, errs.ErrInvalidVacancyID)
		return
	}

	vacancy, err := service.GetVacancyByID(vacancyID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, vacancy)
}

// CreateVacancy godoc
// @Summary Создать вакансию
// @Description Создаёт новую вакансию
// @Tags Vacancies
// @Accept json
// @Produce json
// @Param vacancy body models.VacancyResponse true "Информация о вакансии"
// @Success 200 {object} models.VacancyReq "ID созданной вакансии"
// @Failure 400 {object} errs.ErrorResp
// @Failure 500 {object} errs.ErrorResp
// @Router /vacancy [post]
// @Security ApiKeyAuth
func CreateVacancy(c *gin.Context) {
	userID := c.GetUint(middlewares.UserIDCtx)

	var vacancy models.Vacancy
	if err := c.ShouldBindJSON(&vacancy); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	vacancy.WorkerID = int(userID)

	err := validators.ValidateVacancy(vacancy)
	if err != nil {
		HandleError(c, err)
		return
	}

	vacancyId, err := service.CreateVacancy(vacancy)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vacancy created successfully",
		"id":      vacancyId,
	})
}

// UpdateVacancy godoc
// @Summary Обновить вакансию
// @Description Обновляет вакансию по ID (только владелец может обновить)
// @Tags Vacancies
// @Accept json
// @Produce json
// @Param id path int true "ID вакансии"
// @Param vacancy body models.VacancyResponse true "Обновлённые данные"
// @Success 200 {object} models.VacancyReq "Успешное сообщение и ID"
// @Failure 400 {object} errs.ErrorResp
// @Failure 403 {object} errs.ErrorResp
// @Failure 500 {object} errs.ErrorResp
// @Router /vacancy/{id} [patch]
// @Security ApiKeyAuth
func UpdateVacancy(c *gin.Context) {
	vacancyIDStr := c.Param("id")
	vacancyID, err := strconv.Atoi(vacancyIDStr)
	if err != nil {
		HandleError(c, errs.ErrInvalidVacancyID)
		return
	}

	var vacancy models.Vacancy
	if err := c.ShouldBindJSON(&vacancy); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	err = validators.ValidateVacancy(vacancy)
	if err != nil {
		HandleError(c, err)
		return
	}

	userID := c.GetUint(middlewares.UserIDCtx)

	vacancy.WorkerID = int(userID)

	err = service.UpdateVacancy(int(userID), vacancyID, vacancy)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vacancy updated successfully",
		"id":      vacancyID,
	})
}

// DeleteVacancyByID
// @Summary Удалить вакансию
// @Description Удаляет вакансию по ID (только владелец может удалить)
// @Tags Vacancies
// @Produce json
// @Param id path int true "ID вакансии"
// @Success 200 {object} models.VacancyReq "Успешное сообщение и ID"
// @Failure 400 {object} errs.ErrorResp
// @Failure 403 {object} errs.ErrorResp
// @Failure 500 {object} errs.ErrorResp
// @Router /vacancy/{id} [delete]
// @Security ApiKeyAuth
func DeleteVacancyByID(c *gin.Context) {
	vacancyIDStr := c.Param("id")
	vacancyID, err := strconv.Atoi(vacancyIDStr)
	if err != nil {
		HandleError(c, errs.ErrInvalidVacancyID)
		return
	}

	userID := c.GetUint(middlewares.UserIDCtx)

	err = service.DeleteVacancy(int(userID), vacancyID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vacancy deleted successfully",
		"id":      vacancyID,
	})
}
