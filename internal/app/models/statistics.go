package models

type UsersStatistic struct {
	Subject  string `json:"subject"`
	Quantity int    `json:"quantity"`
}

type CompanyStatisticJSONUtil struct {
	Title                   string   `json:"title"`
	Subject                 []string `json:"subject"`
	FrequencyOfPublications float32  `json:"frequency_of_publications"`
}

type CompanyStatistic struct {
	Subject                  []string                 `json:"subject"`
	FrequencyOfPublications  float32                  `json:"frequency_of_publications"`
	CompanyProductsStatistic CompanyProductsStatistic `json:"company_products"`
}

type CompanyProductsStatistic struct {
	QuantityVacancies int     `json:"quantity_vacancies"`
	QuantityCourses   int     `json:"quantity_courses"`
	Company           Company `json:"company"`
	AverageSalary     float64 `json:"average_salary"`
}

type CompanyInfo struct {
	Company       Company   `json:"company"`
	Vacancies     []Vacancy `json:"vacancies"`
	Courses       []Course  `json:"courses"`
	AverageSalary float64   `json:"average_salary"`
}
