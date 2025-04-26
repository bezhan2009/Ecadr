package repository

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/db"
	"Ecadr/pkg/logger"
)

func GetAllCourses(search string) (courses []models.Course, err error) {
	query := db.GetDBConn().Where("deleted_at IS NULL") // исключаем удалённые

	if search != "" {
		likeSearch := "%" + search + "%"

		query = query.Where(`
			title ILIKE ?
			OR description ILIKE ?
			OR subject ILIKE ?
			OR CAST(worker_id AS TEXT) ILIKE ?
			OR CAST(start_date AS TEXT) ILIKE ?
			OR CAST(end_date AS TEXT) ILIKE ?
			OR EXISTS (
				SELECT 1 FROM unnest(tags) AS tag WHERE tag ILIKE ?
			)
		`, likeSearch, likeSearch, likeSearch, likeSearch, likeSearch, likeSearch, likeSearch)
	}

	if err = query.Find(&courses).Error; err != nil {
		logger.Error.Printf("[repository.GetAllCourses] error while getting courses: %v", err)
		return []models.Course{}, TranslateGormError(err)
	}

	return courses, nil
}

func GetAllWorkerCourses(workerID uint, search string) (courses []models.Course, err error) {
	query := db.GetDBConn().Where("deleted_at IS NULL") // исключаем удалённые

	if search != "" {
		likeSearch := "%" + search + "%"

		query = query.Where(`
			title ILIKE ?
			OR description ILIKE ?
			OR subject ILIKE ?
			OR CAST(worker_id AS TEXT) ILIKE ?
			OR CAST(start_date AS TEXT) ILIKE ?
			OR CAST(end_date AS TEXT) ILIKE ?
			OR EXISTS (
				SELECT 1 FROM unnest(tags) AS tag WHERE tag ILIKE ?
			)
		`, likeSearch, likeSearch, likeSearch, likeSearch, likeSearch, likeSearch, likeSearch)
	}

	if err = query.Where("worker_id = ? AND deleted_at IS NULL", workerID).Find(&courses).Error; err != nil {
		logger.Error.Printf("[repository.GetAllWorkerCourses] Error while getting worker courses: %v", err)
		return []models.Course{}, TranslateGormError(err)
	}

	return courses, nil
}

func GetCourseById(courseId int) (course models.Course, err error) {
	if err = db.GetDBConn().Model(&models.Course{}).Where("id = ?", courseId).First(&course).Error; err != nil {
		logger.Error.Printf("[repository.GetCourseById] error while getting course: %v", err)

		return models.Course{}, TranslateGormError(err)
	}

	return course, nil
}

func CreateCourse(course models.Course) (id int, err error) {
	if err = db.GetDBConn().Model(&models.Course{}).Create(&course).Error; err != nil {
		logger.Error.Printf("[repository.CreateCourse] error while creating course: %v", err)

		return 0, TranslateGormError(err)
	}

	return int(course.ID), nil
}

func UpdateCourse(course models.Course) (err error) {
	if err = db.GetDBConn().Model(&models.Course{}).Where("id = ?", course.ID).Updates(&course).Error; err != nil {
		logger.Error.Printf("[repository.UpdateCourse] error while updating course: %v", err)

		return TranslateGormError(err)
	}

	return nil
}

func DeleteCourse(userID, courseID int) (err error) {
	if err = db.GetDBConn().Model(&models.Course{}).Where("worker_id = ?", userID).Delete(&models.Course{}, courseID).Error; err != nil {
		logger.Error.Printf("[repository.DeleteCourse] error while deleting course: %v", err)

		return TranslateGormError(err)
	}

	return nil
}
