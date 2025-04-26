package service

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service/validators"
	"Ecadr/internal/repository"
)

func GetCompanyByID(companyID uint) (company models.Company, err error) {
	company, err = repository.GetCompanyByID(companyID)
	if err != nil {
		return models.Company{}, err
	}

	return company, nil
}

func GetCompaniesProducts() (dataCompanies []models.CompanyInfo, err error) {
	dataCompanies, err = repository.GetCompaniesWithVacanciesAndCourses()
	if err != nil {
		return nil, err
	}

	return dataCompanies, nil
}

func GetCompaniesProductsStatistics() (statistics []models.CompanyProductsStatistic, err error) {
	dataCompanies, err := GetCompaniesProducts()
	if err != nil {
		return nil, err
	}

	for _, dataCompany := range dataCompanies {
		statistics = append(statistics, models.CompanyProductsStatistic{
			QuantityVacancies: len(dataCompany.Vacancies),
			QuantityCourses:   len(dataCompany.Courses),
			Company:           dataCompany.Company,
			AverageSalary:     dataCompany.AverageSalary,
		})
	}

	return statistics, nil
}

func CreateCompany(company models.Company) (err error) {
	if err = validators.ValidateCompany(company); err != nil {
		err = repository.DeleteUserByID(company.WorkerID)
		if err != nil {
			return err
		}

		return err
	}

	err = repository.CreateCompany(company)
	if err != nil {
		err = repository.DeleteUserByID(company.WorkerID)
		if err != nil {
			return err
		}

		return err
	}

	return nil
}

func UpdateCompany(company models.Company) (err error) {
	if err = validators.ValidateCompany(company); err != nil {
		return err
	}

	err = repository.UpdateCompany(company)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCompany(companyID uint) (err error) {
	err = repository.DeleteCompany(companyID)
	if err != nil {
		return err
	}

	return nil
}
