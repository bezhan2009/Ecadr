package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckWorkerDepartment(c *gin.Context) {
	roleID := c.GetUint(UserRoleIDCtx)
	if roleID != 5 {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{
				"message": "Permission denied",
			},
		)
		return
	}

	c.Next()
}
