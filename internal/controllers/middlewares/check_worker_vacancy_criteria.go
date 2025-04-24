package middlewares

import (
	"Ecadr/internal/app/service"
	"Ecadr/pkg/errs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CheckWorkerVacancyCriteria(c *gin.Context) {
	userID := c.GetUint(UserIDCtx)
	criteriaStrID := c.Param("id")
	criteriaID, err := strconv.Atoi(criteriaStrID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error": errs.ErrInvalidID.Error(),
			})
		return
	}

	criteria, err := service.GetVacancyCriteriaByID(uint(criteriaID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"error": errs.ErrRecordNotFound.Error(),
			})
		return
	}

	vacancy, err := service.GetVacancyByID(int(criteria.VacancyID))
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
