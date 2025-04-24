package controllers

import (
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Обработка ошибок, которые приводят к статусу 400 (Bad Request)
func handleBadRequestErrors(err error) bool {
	return errors.Is(err, errs.ErrUsernameUniquenessFailed) ||
		errors.Is(err, errs.ErrUsernameIsRequired) ||
		errors.Is(err, errs.ErrPasswordIsRequired) ||
		errors.Is(err, errs.ErrRoleIsRequired) ||
		errors.Is(err, errs.ErrWrongRoleID) ||
		errors.Is(err, errs.ErrInvalidVacancyID) ||
		errors.Is(err, errs.ErrNoAIRecommends) ||
		errors.Is(err, errs.ErrInvalidDate) ||
		errors.Is(err, errs.ErrEndDateBeforeStartDate) ||
		errors.Is(err, errs.ErrInvalidCourseID) ||
		errors.Is(err, errs.ErrInvalidType) ||
		errors.Is(err, errs.ErrAlreadyExists) ||
		errors.Is(err, errs.ErrFirstNameIsRequired) ||
		errors.Is(err, errs.ErrLastNameIsRequired) ||
		errors.Is(err, errs.ErrEmailIsRequired) ||
		errors.Is(err, errs.ErrIncorrectUsernameOrPassword) ||
		errors.Is(err, errs.ErrInvalidCredentials) ||
		errors.Is(err, errs.ErrInvalidStoreReviewID) ||
		errors.Is(err, errs.ErrPathParametrized) ||
		errors.Is(err, errs.ErrInvalidProductID) ||
		errors.Is(err, errs.ErrInvalidAddressID) ||
		errors.Is(err, errs.ErrInvalidOrderID) ||
		errors.Is(err, errs.ErrInvalidProductReviewID) ||
		errors.Is(err, errs.ErrInvalidAccountID) ||
		errors.Is(err, errs.ErrInvalidFeaturedProductID) ||
		errors.Is(err, errs.ErrInvalidAddressName) ||
		errors.Is(err, errs.ErrInvalidAccountNumber) ||
		errors.Is(err, errs.ErrInvalidRecommendIDs) ||
		errors.Is(err, errs.ErrInvalidMinPrice) ||
		errors.Is(err, errs.ErrInvalidMaxPrice) ||
		errors.Is(err, errs.ErrInvalidPrice) ||
		errors.Is(err, errs.ErrInvalidID) ||
		errors.Is(err, errs.ErrInvalidPaymentID) ||
		errors.Is(err, errs.ErrInvalidRating) ||
		errors.Is(err, errs.ErrInvalidContent) ||
		errors.Is(err, errs.ErrInvalidComment) ||
		errors.Is(err, errs.ErrInvalidField) ||
		errors.Is(err, errs.ErrInvalidCategory) ||
		errors.Is(err, errs.ErrEmailIsEmpty) ||
		errors.Is(err, errs.ErrPasswordIsEmpty) ||
		errors.Is(err, errs.ErrUsernameIsEmpty) ||
		errors.Is(err, errs.ErrInvalidStore) ||
		errors.Is(err, errs.ErrInvalidStoreID) ||
		errors.Is(err, errs.ErrValidationFailed) ||
		errors.Is(err, errs.ErrNotEnoughProductInStock) ||
		errors.Is(err, errs.ErrDeleteFailed) ||
		errors.Is(err, errs.ErrInvalidTitle) ||
		errors.Is(err, errs.ErrOrderAlreadyPayed) ||
		errors.Is(err, errs.ErrInvalidDescription) ||
		errors.Is(err, errs.ErrInvalidAmount) ||
		errors.Is(err, errs.ErrInvalidQuantity) ||
		errors.Is(err, errs.ErrInsufficientFunds)
}

// Обработка ошибок, которые приводят к статусу 404 (Not Found)
func handleNotFoundErrors(err error) bool {
	return errors.Is(err, errs.ErrRecordNotFound)
}

// Обработка ошибок, которые приводят к статусу 401 (Unauthorized)
func handleUnauthorizedErrors(err error) bool {
	return errors.Is(err, errs.ErrInvalidToken) ||
		errors.Is(err, errs.ErrUnauthorized) ||
		errors.Is(err, errs.ErrUserNotFound) ||
		errors.Is(err, errs.ErrRefreshTokenExpired)
}

// HandleError Основная функция обработки ошибок
func HandleError(c *gin.Context, err error) {
	if handleBadRequestErrors(err) {
		c.JSON(http.StatusBadRequest, newErrorResponse(err.Error()))
	} else if errors.Is(err, errs.ErrPermissionDenied) {
		c.JSON(http.StatusForbidden, newErrorResponse(err.Error()))
	} else if handleNotFoundErrors(err) {
		c.JSON(http.StatusNotFound, newErrorResponse(err.Error()))
	} else if handleUnauthorizedErrors(err) {
		c.JSON(http.StatusUnauthorized, newErrorResponse(err.Error()))
	} else {
		logger.Error.Printf("Err: %s", err)
		c.JSON(http.StatusInternalServerError, newErrorResponse(errs.ErrSomethingWentWrong.Error()))
	}
}
