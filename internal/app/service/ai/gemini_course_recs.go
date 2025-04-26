package aiSerivce

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"encoding/json"
	"fmt"
)

func GetAnalyseForUserCourse(coursesWorker []models.Course,
	kinderNote []models.KindergartenNote,
	schoolGrade []models.SchoolGrade,
	achievementsUser []models.Achievement) (analysedCourses []models.Course, err error) {

	jsons, err := serializeData(nil,
		nil,
		kinderNote,
		schoolGrade,
		achievementsUser,
		coursesWorker,
		nil,
	)
	if err != nil {
		return nil, err
	}

	var text = fmt.Sprintf(
		"Ниже приведены данные о пользователе:\n\n"+
			"1. Успеваемость в детском саду: %v\n"+
			"2. Успеваемость в школе: %v\n"+
			"3. Достижения: %v\n\n"+
			"Также прилагается список доступных курсов: %v\n\n"+
			"Твоя задача — на основе предоставленных данных о пользователе "+
			"выбрать наиболее подходящие курсы из предложенного списка. "+
			"Учитывай как уровень образования, так и достижения пользователя.\n\n"+
			"Ответ верни строго в виде **JSON-массива** от 0 до 10 элементов (включительно).\n"+
			"- Если нет ни одной подходящих курсов — верни `[]` (пустой массив).\n"+
			"- Без лишнего текста, только JSON. Ни комментариев, ни пояснений.\n\n",
		string(jsons[kinderNotes]), string(jsons[schoolGrades]), string(jsons[achievements]), string(jsons[courses]),
	)

	GeminiText, err := sendTextToGemini(text)
	if err != nil {
		return nil, err
	}

	var GeminiTextParse = addBrackets(GeminiText[8 : len(GeminiText)-5])

	if GeminiTextParse == "[]" || len(GeminiTextParse) < 10 {
		return analysedCourses, errs.ErrNoCourseFound
	}

	if err := json.Unmarshal([]byte(GeminiTextParse), &analysedCourses); err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForUserCourse] Error parsing gemini text to courses: %v", err)
		return nil, err
	}

	return analysedCourses, err
}
