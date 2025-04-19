package controllers

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	"Ecadr/internal/app/service/validators"
	"Ecadr/internal/controllers/middlewares"
	"Ecadr/pkg/errs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetAllCourses(c *gin.Context) {
	search := c.Query("search")

	courses, err := service.GetAllCourses(search)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, courses)
}

func GetCourseByID(c *gin.Context) {
	idStr := c.Param("id")
	id, convErr := strconv.Atoi(idStr)
	if convErr != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	course, err := service.GetCourseById(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, course)
}

func CreateCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	if err := validators.ValidateCourse(course); err != nil {
		HandleError(c, err)
		return
	}

	courseID, err := service.CreateCourse(course)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Course Created successfully",
		"course_id": courseID,
	})
}

func UpdateCourse(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		HandleError(c, errs.ErrInvalidVacancyID)
		return
	}

	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	if err := validators.ValidateCourse(course); err != nil {
		HandleError(c, err)
		return
	}

	courseID, err = service.UpdateCourse(courseID, course)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Course Updated successfully",
		"course_id": courseID,
	})
}

func DeleteCourse(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		HandleError(c, errs.ErrInvalidVacancyID)
		return
	}

	userID := c.GetInt(middlewares.UserIDCtx)

	err = service.DeleteCourse(courseID, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Course Deleted successfully",
		"course_id": courseID,
	})
}
