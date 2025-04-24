package middlewares

import (
	"Ecadr/internal/app/service"
	"Ecadr/pkg/errs"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CheckWorkerCourse(c *gin.Context) {
	userID := c.GetUint(UserIDCtx)
	courseStrID := c.Param("id")
	if courseStrID == "" {
		c.Next()
		return
	}

	courseID, err := strconv.Atoi(courseStrID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error": errs.ErrInvalidID.Error(),
			})
		return
	}

	course, err := service.GetCourseById(courseID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"error": errs.ErrRecordNotFound.Error(),
			})
		return
	}

	if uint(course.WorkerID) != userID {
		fmt.Printf("Course ID not match: %d != %d\n", course.WorkerID, userID)
		c.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"error": errs.ErrRecordNotFound.Error(),
			})
		return
	}

	c.Next()
}
