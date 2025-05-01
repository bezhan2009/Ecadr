package aiSerivce

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"Ecadr/pkg/utils"
	"encoding/json"
	"fmt"
)

func GetAnalyseForCourseUser(coursesWorker models.Course,
	usersAnalyse []models.User) (analysedUsers []models.User, err error) {
	jsons, err := utils.SerializeData(usersAnalyse,
		nil,
		nil,
		nil,
		nil,
		[]models.Course{coursesWorker},
		nil,
	)
	if err != nil {
		return nil, err
	}

	var text = fmt.Sprintf(
		"Ниже приведены данные о пользователях:\n\n"+
			"1. Users: %v\n"+
			"Также прилагается курс: %v\n\n"+
			"Твоя задача — на основе предоставленных данных о пользователях "+
			"выбрать наиболее подходящего пользователя для того курса."+
			"Учитывай как уровень образования, так и достижения пользователя.\n\n"+
			"Ответ верни строго в виде **JSON-массива пользователя** от 0 до 10 элементов (включительно).\n"+
			"- Если нет ни одного подходящего пользователя — верни `[]` (пустой массив).\n"+
			"- Без лишнего текста, только JSON. Ни комментариев, ни пояснений.\n\n",
		string(jsons[utils.Users]), string(jsons[utils.Courses]),
	)

	GeminiText, err := utils.SendTextToGemini(text)
	if err != nil {
		return nil, err
	}

	var GeminiTextParse = utils.AddBrackets(GeminiText[8 : len(GeminiText)-5])

	if GeminiTextParse == "[]" || len(GeminiTextParse) < 10 {
		return analysedUsers, errs.ErrNoCourseFound

	}

	if err := json.Unmarshal([]byte(GeminiTextParse), &analysedUsers); err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForCourseUser] Error parsing gemini text to Users: %v", err)
		return nil, err
	}

	return analysedUsers, err
}
