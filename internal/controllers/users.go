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

func GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error.Printf("[controllers.GetUserByID] invalid id: %s\n", c.Param("id"))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})

		return
	}

	user, err := service.GetUserByID(uint(id))
	if err != nil {
		HandleError(c, err)
		logger.Error.Printf("[controllers.GetUserByID] error: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, user)

}

func GetMyData(c *gin.Context) {
	userId := c.GetUint(middlewares.UserIDCtx)

	user, err := service.GetUserByID(userId)
	if err != nil {
		logger.Error.Printf("[controllers.GetMyData] error: %v\n", err)

		HandleError(c, err)

		return
	}

	c.JSON(http.StatusOK, user)
}
