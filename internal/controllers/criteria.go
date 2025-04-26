package controllers

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	"Ecadr/pkg/errs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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
