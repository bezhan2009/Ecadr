package repository

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/db"
	"Ecadr/pkg/logger"
	"strings"
)

func GetCompaniesWithVacanciesAndCourses() ([]models.CompanyInfo, error) {
	var companies []models.Company
	if err := db.GetDBConn().Find(&companies).Error; err != nil {
		return nil, err
	}

	var result []models.CompanyInfo

	for _, company := range companies {
		// Получаем вакансии через твою функцию
		vacancies, err := GetAllWorkerVacancies(company.WorkerID, "")
		if err != nil {
			return nil, err
		}

		// Фильтруем вакансии по компании
		var companyVacancies []models.Vacancy
		var salarySum int
		var salaryCount int

		for _, vacancy := range vacancies {
			if vacancy.CompanyID == company.ID {
				companyVacancies = append(companyVacancies, vacancy)

				if strings.ToLower(vacancy.Type) == strings.ToLower("Работа") ||
					strings.ToLower(vacancy.Type) == strings.ToLower("Work") {
					if vacancy.Salary.Min > 0 && vacancy.Salary.Max > 0 {
						salarySum += (vacancy.Salary.Min + vacancy.Salary.Max) / 2
						salaryCount++
					}
				}
			}
		}

		// Получаем курсы по WorkerID компании
		var courses []models.Course
		if err := db.GetDBConn().Where("worker_id = ?", company.WorkerID).Find(&courses).Error; err != nil {
			return nil, err
		}

		// Считаем среднюю зарплату
		var averageSalary float64
		if salaryCount > 0 {
			averageSalary = float64(salarySum) / float64(salaryCount)
		}

		info := models.CompanyInfo{
			Company:       company,
			Vacancies:     companyVacancies,
			Courses:       courses,
			AverageSalary: averageSalary,
		}
		result = append(result, info)
	}

	return result, nil
}

func GetCompanyByID(companyID uint) (company models.Company, err error) {
	if err = db.GetDBConn().Model(&models.Company{}).Where("id = ?", companyID).First(&company).Error; err != nil {
		logger.Error.Printf("[repository.GetCompanyByID] error while getting company by ID: %v\n", err)

		return models.Company{}, TranslateGormError(err)
	}

	return company, nil
}

func CreateCompany(company models.Company) (err error) {
	if err = db.GetDBConn().Model(&models.Company{}).Create(&company).Error; err != nil {
		logger.Error.Printf("[repository.CreateCompany] error while creating company by ID: %v\n", err)

		return TranslateGormError(err)
	}

	return nil
}

func UpdateCompany(company models.Company) (err error) {
	if err = db.GetDBConn().
		Model(&models.Company{}).
		Where("id = ?", company.ID).
		Updates(company).Error; err != nil {

		logger.Error.Printf("[repository.UpdateCompany] error while updating company by ID: %v\n", err)
		return TranslateGormError(err)
	}

	return nil
}

func DeleteCompany(companyID uint) (err error) {
	if err = db.GetDBConn().Model(&models.Company{}).Delete(&models.Company{}, companyID).Error; err != nil {
		logger.Error.Printf("[repository.DeleteCompany] error while deleting company by ID: %v\n", err)

		return TranslateGormError(err)
	}

	return nil
}
