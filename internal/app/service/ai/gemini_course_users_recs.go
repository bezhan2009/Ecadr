package aiSerivce

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func GetAnalyseForCourseUser(course models.Course,
	users []models.User) (analysedUsers []models.User, err error) {
	jsonUsers, err := json.Marshal(users)
	if err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForCourseUser] Error marshalling kindernote\n\tUsers:%v\n\tError: %v", users, err)
		return nil, err
	}

	jsonCourse, err := json.Marshal(course)
	if err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForVacanciesUser] Error marshalling Vacancy\n\tVacancy:%v\n\tError: %v", course, err)
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
		string(jsonUsers), string(jsonCourse),
	)

	geminiReq := models.GeminiCandidateReq{
		Content: []models.GeminiContents{
			{
				Parts: []models.GeminiParts{
					{
						Text: text,
					},
				},
			},
		},
	}

	// сериализуем в JSON
	jsonBody, err := json.Marshal(geminiReq)
	if err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForVacanciesUser] Error marshalling json body: %v", err)
		return nil, err
	}

	var geminiURL = fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s",
		os.Getenv("GEMINI_API_KEY"),
	)

	analyse, err := http.Post(
		geminiURL,
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForCourseUser] Error getting analyse: %v", err)
		return nil, err
	}
	defer analyse.Body.Close()

	body, err := ioutil.ReadAll(analyse.Body)
	if err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForCourseUser] Error reading body analyse: %v", err)
		return nil, err
	}

	var GeminiResp models.Gemini
	if err := json.Unmarshal(body, &GeminiResp); err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForCourseUser] Error parsing body: %v", err)
		return nil, err
	}

	var GeminiText = GeminiResp.Candidates[0].Content.Parts[0].Text

	var GeminiTextParse = addBrackets(GeminiText[8 : len(GeminiText)-5])

	if GeminiTextParse == "[]" || len(GeminiTextParse) < 10 {
		return analysedUsers, errs.ErrNoCourseFound

	}

	if err := json.Unmarshal([]byte(GeminiTextParse), &analysedUsers); err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForCourseUser] Error parsing gemini text to users: %v", err)
		return nil, err
	}

	return analysedUsers, err
}
