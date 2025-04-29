package controllers

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	"Ecadr/internal/controllers/middlewares"
	"Ecadr/pkg/errs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetTestsByTypeAndID(c *gin.Context) {
	var testReq models.TestSearchRequest
	err := c.BindJSON(&testReq)
	if err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	tests, err := service.GetTestsByTypeAndID(testReq.TargetType, testReq.TargetID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tests)
}

func GetTestByID(c *gin.Context) {
	testIdStr := c.Param("id")
	testId, err := strconv.Atoi(testIdStr)
	if err != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	test, err := service.GetTestByID(uint(testId))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, test)
}

func CreateTest(c *gin.Context) {
	var test models.Test
	err := c.BindJSON(&test)
	if err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	testId, err := service.CreateTest(test)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"testId":  testId,
		"message": "test created successfully",
	})
}

func UpdateTest(c *gin.Context) {
	testIdStr := c.Param("id")
	testId, err := strconv.Atoi(testIdStr)
	if err != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	var testReq models.Test
	err = c.BindJSON(&testReq)
	if err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	testReq.ID = testId
	_, err = service.UpdateTest(testReq)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"testId":  testId,
		"message": "test updated successfully",
	})
}

func DeleteTest(c *gin.Context) {
	testIdStr := c.Param("id")
	testId, err := strconv.Atoi(testIdStr)
	if err != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	err = service.DeleteTest(uint(testId))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"testId":  testId,
		"message": "test deleted successfully",
	})
}

func GetSortedSubmissionsHandler(c *gin.Context) {
	testIdStr := c.Param("id")
	testId, err := strconv.Atoi(testIdStr)
	if err != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	sorted, err := service.SortSubmissionsByScore(uint(testId))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, sorted)
}

func CreateTestSubmission(c *gin.Context) {
	userID := c.GetUint(middlewares.UserIDCtx)

	var testReq models.TestSubmission
	err := c.BindJSON(&testReq)
	if err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	testReq.UserID = int(userID)

	err = service.CreateTestSubmission(testReq)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "test submission created successfully",
	})
}
