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

// GetAllWorkerCourses godoc
// @Summary Получить курс по workerID
// @Description Возвращает курс по указанному workerID
// @Tags Courses
// @Produce json
// @Param worker_id path int true "workerID курса"
// @Success 200 {object} models.CourseResponse
// @Failure 400 {object} errs.ErrorResp
// @Failure 404 {object} errs.ErrorResp
// @Router /course/company/{company_id} [get]
func GetAllWorkerCourses(c *gin.Context) {
	search := c.Query("search")
	companyIDStr := c.Param("company_id")
	companyID, convErr := strconv.Atoi(companyIDStr)
	if convErr != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	company, err := service.GetCompanyByID(uint(companyID))
	if err != nil {
		HandleError(c, errs.ErrInvalidID)
		return
	}

	courses, err := service.GetAllWorkerCourses(company.WorkerID, search)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, courses)
}

// GetCourseByID godoc
// @Summary Получить курс по ID
// @Description Возвращает курс по указанному ID
// @Tags Courses
// @Produce json
// @Param id path int true "ID курса"
// @Success 200 {object} models.CourseResponse
// @Failure 400 {object} errs.ErrorResp
// @Failure 404 {object} errs.ErrorResp
// @Router /course/{id} [get]
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

// CreateCourse godoc
// @Summary Создать курс
// @Description Создаёт новый курс
// @Tags Courses
// @Accept json
// @Produce json
// @Param course body models.CourseResponse true "Информация о курсе"
// @Success 200 {object} models.CourseReq "ID созданного курса"
// @Failure 400 {object} errs.ErrorResp
// @Failure 500 {object} errs.ErrorResp
// @Router /course [post]
// @Security ApiKeyAuth
func CreateCourse(c *gin.Context) {
	userID := c.GetUint(middlewares.UserIDCtx)

	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	course.WorkerID = int(userID)

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

// UpdateCourse godoc
// @Summary Обновить курс
// @Description Обновляет данные существующего курса по ID
// @Tags Courses
// @Accept json
// @Produce json
// @Param id path int true "ID курса"
// @Param course body models.CourseResponse true "Обновлённые данные курса"
// @Success 200 {object} models.CourseReq "Успешное сообщение и ID"
// @Failure 400 {object} errs.ErrorResp
// @Failure 500 {object} errs.ErrorResp
// @Router /course/{id} [patch]
// @Security ApiKeyAuth
func UpdateCourse(c *gin.Context) {
	userID := c.GetUint(middlewares.UserIDCtx)
	courseIDStr := c.Param("id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		HandleError(c, errs.ErrInvalidCourseID)
		return
	}

	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	course.ID = uint(courseID)
	course.WorkerID = int(userID)

	if err := validators.ValidateCourse(course); err != nil {
		HandleError(c, err)
		return
	}

	courseID, err = service.UpdateCourse(course)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Course Updated successfully",
		"course_id": courseID,
	})
}

// DeleteCourse godoc
// @Summary Удалить курс
// @Description Удаляет курс по ID (только автор может удалить)
// @Tags Courses
// @Produce json
// @Param id path int true "ID курса"
// @Success 200 {object} models.CourseReq "Успешное сообщение и ID"
// @Failure 400 {object} errs.ErrorResp
// @Failure 403 {object} errs.ErrorResp
// @Failure 500 {object} errs.ErrorResp
// @Router /course/{id} [delete]
// @Security ApiKeyAuth
func DeleteCourse(c *gin.Context) {
	userID := c.GetUint(middlewares.UserIDCtx)
	courseIDStr := c.Param("id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		HandleError(c, errs.ErrInvalidVacancyID)
		return
	}

	err = service.DeleteCourse(int(userID), courseID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Course Deleted successfully",
		"course_id": courseID,
	})
}
