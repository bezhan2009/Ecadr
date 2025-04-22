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

func GetRecommendsForUser(kinderNote []models.KindergartenNote,
	schoolGrades []models.SchoolGrade,
	achievements []models.Achievement) (recommendsUser []models.AiUserRecommends, err error) {

	jsonKinderNote, err := json.Marshal(kinderNote)
	if err != nil {
		logger.Error.Printf("[aiService.GetRecommendsForUser] Error marshalling kindernote\n\tKinderNote:%v\n\tError: %v", kinderNote, err)
		return nil, err
	}

	jsonSchoolGrades, err := json.Marshal(schoolGrades)
	if err != nil {
		logger.Error.Printf("[aiService.GetRecommendsForUser] Error marshalling SchoolGrades\n\tSchoolGrades:%v\n\tError: %v", schoolGrades, err)
		return nil, err
	}

	jsonAchievements, err := json.Marshal(achievements)
	if err != nil {
		logger.Error.Printf("[aiService.GetRecommendsForUser] Error marshalling achievements\n\tachievements:%v\n\tError: %v", achievements, err)
		return nil, err
	}

	var text = fmt.Sprintf(
		"Ниже приведены данные о пользователе:\n\n"+
			"1. Успеваемость в детском саду: %v\n"+
			"2. Успеваемость в школе: %v\n"+
			"3. Достижения: %v\n\n"+
			"Твоя задача — написать ему несколько рекомендаций, что он может улучшить в себе(по предметам)"+
			"Ответ верни строго в виде **JSON-массива** от 1 до 3 элементов (включительно).\n"+
			"- Если нет ни одной подходящей рекомендации от тебя к нему — верни `[]` (пустой массив).\n"+
			"- Без лишнего текста, только JSON. Ни комментариев, ни пояснений.\n\n"+
			"ты должен возвращать такой JSON: { Subject: \"Предмет который надо подтянуть\", Recommend: \"Сама рекомендаця\" } и желательно массив из этих объектов даже если объект один",
		string(jsonKinderNote), string(jsonSchoolGrades), string(jsonAchievements),
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
		logger.Error.Printf("[aiService.GetRecommendsForUser] Error marshalling json body: %v", err)
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
		logger.Error.Printf("[aiService.GetRecommendsForUser] Error getting analyse: %v", err)
		return nil, err
	}
	defer analyse.Body.Close()

	body, err := ioutil.ReadAll(analyse.Body)
	if err != nil {
		logger.Error.Printf("[aiService.GetRecommendsForUser] Error reading body analyse: %v", err)
		return nil, err
	}

	var GeminiResp models.Gemini

	if err := json.Unmarshal(body, &GeminiResp); err != nil {
		logger.Error.Printf("[aiService.GetRecommendsForUser] Error parsing body: %v", err)
		return nil, err
	}

	var GeminiText = GeminiResp.Candidates[0].Content.Parts[0].Text

	var GeminiTextParse = addBrackets(GeminiText[8 : len(GeminiText)-5])

	if GeminiTextParse == "[]" || len(GeminiTextParse) < 5 {
		return recommendsUser, errs.ErrNoAIRecommends
	}

	if err = json.Unmarshal([]byte(GeminiTextParse), &recommendsUser); err != nil {
		logger.Error.Printf("[aiService.GetRecommendsForUser] Error parsing body: %v", err)

		return recommendsUser, err
	}

	return recommendsUser, nil
}
