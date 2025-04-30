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

// GetTestsByTypeAndID godoc
// @Summary      Get tests by target type and ID
// @Description  Get list of tests for vacancy/course
// @Tags         tests
// @Accept       json
// @Produce      json
// @Param        request  body      models.TestSearchRequest  true  "Target info"
// @Success      200      {array}   models.TestResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /tests/search [post]
func GetTestsByTypeAndID(c *gin.Context) {
	var testReq models.TestSearchRequest
	err := c.BindJSON(&testReq)
	if err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	tests, err := service.GetTestsByTypeAndID(uint(testReq.TargetType), uint(testReq.TargetID))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tests)
}

// GetTestByID godoc
// @Summary      Get test by ID
// @Description  Returns a test by its ID
// @Tags         tests
// @Param        id   path      int  true  "Test ID"
// @Success      200  {object}  models.TestResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /tests/{id} [get]
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

// CreateTest godoc
// @Summary      Create new test
// @Description  Creates a test
// @Tags         tests
// @Accept       json
// @Produce      json
// @Param        test  body      models.TestRequest  true  "Test object"
// @Success      200   {object}  models.DefaultResponse
// @Failure      400   {object}  models.ErrorResponse
// @Router       /tests [post]
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

// UpdateTest godoc
// @Summary      Update test
// @Description  Updates an existing test
// @Tags         tests
// @Accept       json
// @Produce      json
// @Param        id    path      int             true  "Test ID"
// @Param        test  body      models.TestRequest true  "Test object"
// @Success      200   {object}  models.TestsRouteCRUDReq
// @Failure      400   {object}  models.ErrorResponse
// @Router       /tests/{id} [put]
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

// DeleteTest godoc
// @Summary      Delete test
// @Description  Deletes a test by ID
// @Tags         tests
// @Param        id   path      int  true  "Test ID"
// @Success      200  {object}  models.TestsRouteCRUDReq
// @Failure      400  {object}  models.ErrorResponse
// @Router       /tests/{id} [delete]
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

// GetSortedSubmissionsHandler godoc
// @Summary      Submit test answers
// @Description  Submit test answers for a test
// @Tags         tests
// @Accept       json
// @Produce      json
// @Success      200         {object}  models.TestResponse
// @Failure      400         {object}  models.ErrorResponse
// @Router       /tests/submission/{id} [get]
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

// CreateTestSubmission godoc
// @Summary      Submit test answers
// @Description  Submit test answers for a test
// @Tags         tests
// @Accept       json
// @Produce      json
// @Param        submission  body      models.TestSubmissionRequest  true  "Submission"
// @Success      200         {object}  models.TestSubmissionResponse
// @Failure      400         {object}  models.ErrorResponse
// @Router       /tests/submission [post]
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
