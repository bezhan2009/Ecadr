package aiSerivce

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"encoding/json"
	"fmt"
)

func GetRecommendsForUser(kinderNote []models.KindergartenNote,
	schoolGrade []models.SchoolGrade,
	achievementsUser []models.Achievement) (recommendsUser []models.AiUserRecommends, err error) {

	jsons, err := serializeData(
		nil,
		nil,
		kinderNote,
		schoolGrade,
		achievementsUser,
		nil,
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
			"Твоя задача — написать ему несколько рекомендаций, что он может улучшить в себе(по предметам)"+
			"Ответ верни строго в виде **JSON-массива** от 1 до 3 элементов (включительно).\n"+
			"- Если нет ни одной подходящей рекомендации от тебя к нему — верни `[]` (пустой массив).\n"+
			"- Без лишнего текста, только JSON. Ни комментариев, ни пояснений.\n\n"+
			"ты должен возвращать такой JSON: { Subject: \"Предмет который надо подтянуть\", Recommend: \"Сама рекомендаця\" } и желательно массив из этих объектов даже если объект один",
		string(jsons[kinderNotes]), string(jsons[schoolGrades]), string(jsons[achievements]),
	)

	GeminiText, err := sendTextToGemini(text)
	if err != nil {
		return nil, err
	}

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
