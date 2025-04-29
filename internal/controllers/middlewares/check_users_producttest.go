package middlewares

import (
	"Ecadr/internal/app/service"
	"Ecadr/internal/controllers"
	"Ecadr/pkg/errs"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CheckUserTest(c *gin.Context) {
	userID := c.GetUint(UserIDCtx)

	testIDStr := c.Param("id")
	testID, err := strconv.Atoi(testIDStr)
	if err != nil {
		controllers.HandleError(c, errs.ErrValidationFailed)
		return
	}

	test, err := service.GetTestByID(uint(testID))
	if err != nil {
		controllers.HandleError(c, err)
		return
	}

	if test.TargetType == 1 {
		vacancy, err := service.GetVacancyByID(test.TargetID)
		if err != nil {
			controllers.HandleError(c, err)
			return
		}

		if uint(vacancy.WorkerID) == userID {
			controllers.HandleError(c, errs.ErrPermissionDenied)
			return
		}
	} else if test.TargetType == 2 {
		course, err := service.GetCourseById(test.TargetID)
		if err != nil {
			controllers.HandleError(c, err)
			return
		}

		if uint(course.WorkerID) == userID {
			controllers.HandleError(c, errs.ErrPermissionDenied)
			return
		}
	}

	c.Next()
}
