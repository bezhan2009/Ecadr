package ai

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	aiSerivce "Ecadr/internal/app/service/ai"
	"Ecadr/internal/controllers"
	"Ecadr/internal/controllers/middlewares"
	"Ecadr/pkg/errs"
	"errors"
	"github.com/gin-gonic/gin"
)

func GetAnalyseForUserCourse(c *gin.Context) {
	search := c.Query("search")
	userID := c.GetUint(middlewares.UserIDCtx)

	if userID == 0 {
		controllers.HandleError(c, errs.ErrUserNotFound)
		return
	}

	userData, err := service.GetUserByID(userID)
	if err != nil {
		controllers.HandleError(c, err)
		return
	}

	var courses []models.Course
	if search != "" {
		courses, err = service.GetAllCourses(search)
		if err != nil {
			controllers.HandleError(c, err)
			return
		}

		c.JSON(200, courses)
		return
	}

	courses, err = service.GetAllCourses(search)
	if err != nil {
		controllers.HandleError(c, err)
		return
	}

	analysedCourse, err := aiSerivce.GetAnalyseForUserCourse(
		courses,
		userData.KindergartenNotes,
		userData.SchoolGrades,
		userData.Achievements,
	)
	if err != nil {
		if errors.Is(err, errs.ErrCourseNotFound) {
			c.JSON(200, courses)
			return
		}

		controllers.HandleError(c, err)
		return
	}

	c.JSON(201, analysedCourse)
}
