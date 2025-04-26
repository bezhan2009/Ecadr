package middlewares

import (
	"Ecadr/internal/app/service"
	"Ecadr/pkg/errs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CheckWorkerCompany(c *gin.Context) {
	userID := c.GetUint(UserIDCtx)
	companyStrID := c.Param("id")
	if companyStrID == "" {
		c.Next()
		return
	}

	companyID, err := strconv.Atoi(companyStrID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error": errs.ErrInvalidID.Error(),
			})
		return
	}

	company, err := service.GetCompanyByID(uint(companyID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"error": errs.ErrRecordNotFound.Error(),
			})
		return
	}

	if company.WorkerID != userID {
		c.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"error": errs.ErrRecordNotFound.Error(),
			})
		return
	}

	c.Next()
}
