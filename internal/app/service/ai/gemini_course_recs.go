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

func GetAnalyseForUserCourse(courses []models.Course,
	kinderNote []models.KindergartenNote,
	schoolGrades []models.SchoolGrade,
	achievements []models.Achievement) (analysedCourses []models.Course, err error) {

	jsonKinderNote, err := json.Marshal(kinderNote)
	if err != nil {
		logger.Error.Printf("Error marshalling kindernote\n\tKinderNote:%v\n\tError: %v", kinderNote, err)
		return nil, err
	}

	jsonSchoolGrades, err := json.Marshal(schoolGrades)
	if err != nil {
		logger.Error.Printf("Error marshalling SchoolGrades\n\tSchoolGrades:%v\n\tError: %v", schoolGrades, err)
		return nil, err
	}

	jsonAchievements, err := json.Marshal(achievements)
	if err != nil {
		logger.Error.Printf("Error marshalling achievements\n\tachievements:%v\n\tError: %v", achievements, err)
		return nil, err
	}

	jsonVacancies, err := json.Marshal(courses)
	if err != nil {
		logger.Error.Printf("Error marshalling courses\n\tcourses:%v\n\tError: %v", courses, err)
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
		string(jsonKinderNote), string(jsonSchoolGrades), string(jsonAchievements), string(jsonVacancies),
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
		logger.Error.Printf("[aiService.GetAnalyseForUserVacancies] Error marshalling json body: %v", err)
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
	fmt.Println(analyse)
	if err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForUserVacancies] Error getting analyse: %v", err)
		return nil, err
	}
	defer analyse.Body.Close()

	body, err := ioutil.ReadAll(analyse.Body)
	if err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForUserVacancies] Error reading body analyse: %v", err)
		return nil, err
	}

	var GeminiResp models.Gemini
	if err := json.Unmarshal(body, &GeminiResp); err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForUserVacancies] Error parsing body: %v", err)
		return nil, err
	}

	var GeminiText = GeminiResp.Candidates[0].Content.Parts[0].Text

	var GeminiTextParse = addBrackets(GeminiText[8 : len(GeminiText)-5])
	fmt.Println("GeminiTextParse: ", GeminiTextParse)

	if GeminiTextParse == "[]" {
		return analysedCourses, errs.ErrNoCourseFound
	}

	if err := json.Unmarshal([]byte(GeminiTextParse), &analysedCourses); err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForUserVacancies] Error parsing gemini text to courses: %v", err)
		return nil, err
	}

	return analysedCourses, err
}
