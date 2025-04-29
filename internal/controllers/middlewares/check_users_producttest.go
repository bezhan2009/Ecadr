package middlewares

import (
	"Ecadr/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CheckUserTest(c *gin.Context) {
	userID := c.GetUint(UserIDCtx)

	testIDStr := c.Param("id")
	testID, err := strconv.Atoi(testIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Permission denied"})
		return
	}

	test, err := service.GetTestByID(uint(testID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Permission denied"})
		return
	}

	if test.TargetType == 1 {
		vacancy, err := service.GetVacancyByID(test.TargetID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Permission denied"})
			return
		}

		if uint(vacancy.WorkerID) == userID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Permission denied"})
			return
		}
	} else if test.TargetType == 2 {
		course, err := service.GetCourseById(test.TargetID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Permission denied"})
			return
		}

		if uint(course.WorkerID) == userID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Permission denied"})
			return
		}
	}

	c.Next()
}
