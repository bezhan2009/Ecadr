package middlewares

import (
	"Ecadr/internal/app/service"
	"Ecadr/pkg/errs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CheckWorkerVacancy(c *gin.Context) {
	userID := c.GetUint(UserIDCtx)
	vacancyStrID := c.Param("id")
	if vacancyStrID == "" {
		c.Next()
		return
	}

	vacancyID, err := strconv.Atoi(vacancyStrID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error": errs.ErrInvalidID.Error(),
			})
		return
	}

	vacancy, err := service.GetVacancyByID(vacancyID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"error": errs.ErrRecordNotFound.Error(),
			})
		return
	}

	if uint(vacancy.WorkerID) != userID {
		c.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"error": errs.ErrRecordNotFound.Error(),
			})
		return
	}

	c.Next()
}
