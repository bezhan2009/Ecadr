package repository

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/db"
	"Ecadr/pkg/logger"
)

func GetAllCourses(search string) (courses []models.Course, err error) {
	query := db.GetDBConn()

	if search != "" {
		// Пример: ищем по title и description
		likeSearch := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", likeSearch, likeSearch)
	}

	if err = query.Find(&courses).Error; err != nil {
		logger.Error.Printf("[repository.GetAllCourses] error while getting courses: %v", err)

		return []models.Course{}, TranslateGormError(err)
	}

	return courses, nil
}

func GetCourseById(courseId int) (course models.Course, err error) {
	if err = db.GetDBConn().Model(&models.Course{}).Where("id = ?", courseId).Find(&course).Error; err != nil {
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
	if err = db.GetDBConn().Model(&models.Course{}).Save(&course).Error; err != nil {
		logger.Error.Printf("[repository.UpdateCourse] error while updating course: %v", err)

		return TranslateGormError(err)
	}

	return nil
}

func DeleteCourse(userID, courseID int) (err error) {
	if err = db.GetDBConn().Model(&models.Course{}).Where("user_id = ?", userID).Delete(&models.Course{}, courseID).Error; err != nil {
		logger.Error.Printf("[repository.DeleteCourse] error while deleting course: %v", err)

		return TranslateGormError(err)
	}

	return nil
}
