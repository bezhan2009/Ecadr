package controllers

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	"Ecadr/pkg/errs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetCompanyByID(c *gin.Context) {
	companyIDStr := c.Param("id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	company, err := service.GetCompanyByID(uint(companyID))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, company)
}

func UpdateCompany(c *gin.Context) {
	companyIDStr := c.Param("id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	company.ID = uint(companyID)
	if err := service.UpdateCompany(company); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update company successfully"})
}

func DeleteCompany(c *gin.Context) {
	companyIDStr := c.Param("id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	if err = service.DeleteCompany(uint(companyID)); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delete company successfully"})
}
