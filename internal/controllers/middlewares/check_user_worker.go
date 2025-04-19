package middlewares

import "github.com/gin-gonic/gin"

func CheckUserWorker(c *gin.Context) {
	if c.GetUint(UserRoleIDCtx) != 3 {
		c.AbortWithStatusJSON(403, gin.H{
			"message": "Permission denied",
		})
		return
	}

	c.Next()
}
