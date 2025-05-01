package aiSerivce

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"Ecadr/pkg/utils"
	"encoding/json"
	"fmt"
)

func GetAnalyseForUserVacancies(vacanciesWorker []models.Vacancy,
	kinderNote []models.KindergartenNote,
	schoolGrade []models.SchoolGrade,
	achievementsUser []models.Achievement) (analysedVacancies []models.Vacancy, err error) {

	jsons, err := utils.SerializeData(
		nil,
		nil,
		kinderNote,
		schoolGrade,
		achievementsUser,
		nil,
		vacanciesWorker,
	)
	if err != nil {
		return nil, err
	}

	var text = fmt.Sprintf(
		"Ниже приведены данные о пользователе:\n\n"+
			"1. Успеваемость в детском саду: %v\n"+
			"2. Успеваемость в школе: %v\n"+
			"3. Достижения: %v\n\n"+
			"Также прилагается список доступных вакансий: %v\n\n"+
			"Твоя задача — на основе предоставленных данных о пользователе "+
			"выбрать наиболее подходящие вакансии из предложенного списка. "+
			"Учитывай как уровень образования, так и достижения пользователя.\n\n"+
			"Ответ верни строго в виде **JSON-массива** от 0 до 10 элементов (включительно).\n"+
			"- Если нет ни одной подходящей вакансии — верни `[]` (пустой массив).\n"+
			"- Без лишнего текста, только JSON. Ни комментариев, ни пояснений.\n\n",
		string(jsons[utils.KinderNotes]), string(jsons[utils.SchoolGrades]), string(jsons[utils.Achievements]), string(jsons[utils.Vacancies]),
	)

	GeminiText, err := utils.SendTextToGemini(text)
	if err != nil {
		return nil, err
	}

	var GeminiTextParse = utils.AddBrackets(GeminiText[8 : len(GeminiText)-5])

	if GeminiTextParse == "[]" || len(GeminiTextParse) < 10 {
		return analysedVacancies, errs.ErrNoVacancyFound
	}

	if err := json.Unmarshal([]byte(GeminiTextParse), &analysedVacancies); err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForUserVacancies] Error parsing gemini text to vacanciesWorker: %v", err)
		return nil, err
	}

	return analysedVacancies, err
}
