package service

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/repository"
	"Ecadr/pkg/errs"
)

func GetAllCourses(search string) (courses []models.Course, err error) {
	courses, err = repository.GetAllCourses(search)
	if err != nil {
		return []models.Course{}, err
	}

	return courses, nil
}

func GetAllWorkerCourses(workerID uint, search string) (courses []models.Course, err error) {
	courses, err = repository.GetAllWorkerCourses(workerID, search)
	if err != nil {
		return []models.Course{}, err
	}

	return courses, nil
}

func GetCourseById(id int) (course models.Course, err error) {
	course, err = repository.GetCourseById(id)
	if err != nil {
		return models.Course{}, err
	}

	return course, nil
}

func CreateCourse(course models.Course) (id int, err error) {
	id, err = repository.CreateCourse(course)
	if err != nil {
		return id, err
	}

	return id, nil
}

func UpdateCourse(userID int, course models.Course) (id int, err error) {
	if course.WorkerID != userID {
		return id, errs.ErrRecordNotFound
	}

	err = repository.UpdateCourse(course)
	if err != nil {
		return id, err
	}

	return id, nil
}

func DeleteCourse(userID, courseID int) (err error) {
	err = repository.DeleteCourse(userID, courseID)
	if err != nil {
		return err
	}

	return nil
}
