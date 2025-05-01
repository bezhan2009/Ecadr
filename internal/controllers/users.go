package controllers

import (
	"Ecadr/internal/app/service"
	"Ecadr/internal/controllers/middlewares"
	"Ecadr/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetAllUsers(c *gin.Context) {
	search := c.Query("search")

	users, err := service.GetAllUsers(search)
	if err != nil {
		logger.Error.Printf("[controllers.GetAllUsers] error: %v\n", err)

		HandleError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// GetUserByID godoc
// @Summary Получить пользователя по ID
// @Description Возвращает пользователя по его ID
// @Tags users
// @Security ApiKeyAuth
// @Param id path int true "ID пользователя"
// @Produce json
// @Success 200 {object} models.UserRequest
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error.Printf("[controllers.GetUserByID] invalid id: %s\n", c.Param("id"))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})

		return
	}

	user, _, err := service.GetUserByID(uint(id))
	if err != nil {
		HandleError(c, err)
		logger.Error.Printf("[controllers.GetUserByID] error: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetMyData godoc
// @Summary Получить мои данные
// @Description Возвращает данные текущего авторизованного пользователя и его чартер
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/me [get]
func GetMyData(c *gin.Context) {
	userId := c.GetUint(middlewares.UserIDCtx)

	user, charter, err := service.GetUserByID(userId)
	if err != nil {
		logger.Error.Printf("[controllers.GetMyData] error: %v\n", err)

		HandleError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "charter": charter})
}
