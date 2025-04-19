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

func GetAllVacancies(c *gin.Context) {
	searchText := c.Query("search")

	vacancies, err := service.GetAllVacancies(searchText)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": vacancies})
}

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

	c.JSON(http.StatusOK, gin.H{"data": vacancy})
}

func CreateVacancy(c *gin.Context) {
	var vacancy models.Vacancy
	if err := c.ShouldBindJSON(&vacancy); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

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
