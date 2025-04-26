package aiSerivce

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"encoding/json"
	"fmt"
)

func getNeededCompanyByName(companyStatisticUtil models.CompanyStatisticJSONUtil,
	companyStatisticsAnalyse []models.CompanyProductsStatistic) models.CompanyProductsStatistic {

	for _, analyse := range companyStatisticsAnalyse {
		if analyse.Company.Title == companyStatisticUtil.Title {
			return analyse
		}
	}

	// Если не нашли совпадения, вернем пустую структуру
	return models.CompanyProductsStatistic{}
}

func GetCompaniesStatistic(companiesInfo []models.CompanyInfo) ([]models.CompanyStatistic, error) {
	jsons, err := serializeData(
		nil,
		companiesInfo,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	if len(jsons) == 0 {
		return nil, errs.ErrNoUsersStatisticFound
	}

	var text = fmt.Sprintf(
		`Ниже приведены данные о компаниях:

1. Company: %v
Твоя задача — на основе предоставленных данных о компаниях вывести статистику компаний.
Учитывай все данные и во всех данных есть CreatedAt, анализируя их, так же выведи Частоту публикации компаний (тип данных float32).

Ответ верни строго в виде **JSON-массива**:
- Если нет ни одной статистики, либо ты просто не смог анализировать всё, то верни пустой массив [].
- По данным каждой компании определи, в каких предметах она развита (или в каком направлении).
- Верни вот такой JSON для абсолютно каждой компании по очередям:
{ "title": "название", "subject": ["Предмет1", "Предмет2"], "frequency_of_publications": "частота" }.

Без лишнего текста, только JSON. Ни комментариев, ни пояснений.
`,
		string(jsons[companyInfoStats]), // Поясни, пожалуйста, что такое companyInfoStats!
	)

	GeminiText, err := sendTextToGemini(text)
	if err != nil {
		return nil, err
	}

	if len(GeminiText) < 10 {
		return nil, errs.ErrNoUsersStatisticFound
	}

	var GeminiTextParse = addBrackets(GeminiText[8 : len(GeminiText)-5])

	if GeminiTextParse == "[]" || len(GeminiTextParse) < 10 {
		return nil, errs.ErrNoUsersStatisticFound
	}

	var companiesStatisticsUtils []models.CompanyStatisticJSONUtil
	if err := json.Unmarshal([]byte(GeminiTextParse), &companiesStatisticsUtils); err != nil {
		logger.Error.Printf("[aiService.GetCompaniesStatistic] Error parsing gemini text to companies: %v", err)
		return nil, err
	}

	companyStatisticsAnalyse, err := service.GetCompaniesProductsStatistics()
	if err != nil {
		return nil, err
	}

	var companiesStatistic []models.CompanyStatistic

	for _, companyStatisticUtil := range companiesStatisticsUtils {
		statistic := models.CompanyStatistic{
			Subject:                 companyStatisticUtil.Subject,
			FrequencyOfPublications: companyStatisticUtil.FrequencyOfPublications,
			CompanyProductsStatistic: getNeededCompanyByName(
				companyStatisticUtil, companyStatisticsAnalyse,
			),
		}
		companiesStatistic = append(companiesStatistic, statistic)
	}

	return companiesStatistic, nil
}
