package controllers

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	"Ecadr/pkg/errs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetVacancyCriteria godoc
// @Summary Получить критерии вакансии по ID вакансии
// @Description Возвращает список критериев для указанной вакансии
// @Tags Criteria
// @Param id path int true "ID Вакансии"
// @Success 200 {array} models.CriteriaReq
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /vacancy/criteria/{id} [get]
func GetVacancyCriteria(c *gin.Context) {
	vacancyStrID := c.Param("id")
	vacancyID, err := strconv.Atoi(vacancyStrID)
	if err != nil {
		HandleError(c, err)
		return
	}

	criteria, err := service.GetVacancyCriteria(uint(vacancyID))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, criteria)
}

// GetVacancyCriteriaByID godoc
// @Summary Получить критерий по его ID
// @Description Возвращает критерий по его уникальному ID
// @Tags Criteria
// @Param id path int true "ID Критерия"
// @Success 200 {object} models.CriteriaReq
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /vacancy/criteria/spec/{id} [get]
func GetVacancyCriteriaByID(c *gin.Context) {
	criteriaStrID := c.Param("id")
	criteriaID, err := strconv.Atoi(criteriaStrID)
	if err != nil {
		HandleError(c, err)
		return
	}

	criteria, err := service.GetVacancyCriteria(uint(criteriaID))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, criteria)
}

// CreateVacancyCriteria godoc
// @Summary Создать новый критерий
// @Description Создает новый критерий для вакансии
// @Tags Criteria
// @Accept json
// @Produce json
// @Param criteria body models.CriteriaReq true "Данные критерия"
// @Success 200 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /vacancy/criteria [post]
func CreateVacancyCriteria(c *gin.Context) {
	var criteria models.Criteria
	if err := c.ShouldBindJSON(&criteria); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	err := service.CreateCriteria(criteria)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "criteria created successfully"})
}

// UpdateVacancyCriteria godoc
// @Summary Обновить критерий
// @Description Обновляет существующий критерий по его ID
// @Tags Criteria
// @Accept json
// @Produce json
// @Param id path int true "ID Критерия"
// @Success 200 {object} models.CriteriaResp
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /vacancy/criteria/{id} [patch]
func UpdateVacancyCriteria(c *gin.Context) {
	criteriaStrID := c.Param("id")
	criteriaID, err := strconv.Atoi(criteriaStrID)
	if err != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	var criteria models.Criteria
	if err := c.ShouldBindJSON(&criteria); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	criteria.ID = uint(criteriaID)

	err = service.UpdateCriteria(criteria)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated criteria successfully"})
}

// DeleteVacancyCriteria godoc
// @Summary Удалить критерий
// @Description Удаляет критерий по его ID
// @Tags Criteria
// @Param id path int true "ID Критерия"
// @Success 200 {object} models.CriteriaResp
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /vacancy/criteria/{id} [delete]
func DeleteVacancyCriteria(c *gin.Context) {
	criteriaStrID := c.Param("id")
	criteriaID, err := strconv.Atoi(criteriaStrID)
	if err != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	err = service.DeleteCriteria(uint(criteriaID))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delete criteria successfully"})
}
