package aiSerivce

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"encoding/json"
	"fmt"
)

func GetUsersStatistic(usersAnalyse []models.User) (usersStatistic []models.UsersStatistic, err error) {
	jsons, err := serializeData(
		usersAnalyse,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	var text = fmt.Sprintf(
		"Ниже приведены данные о пользователях:\n\n"+
			"1. Users: %v\n"+
			"Твоя задача — на основе предоставленных данных о пользователях "+
			"вывести статистику пользователей"+
			"Учитывай как уровень образования, так и достижения пользователя.\n\n"+
			"Ответ верни строго в виде **JSON-массива**\n"+
			"- Если нет ни одной статистики, либо ты просто не смог анализировать всё, то верни пустой массив `[]`\n"+
			"- По данным каждого пользователя определи, в каком предемете он развит(или в направлении)\n"+
			"- Верни вот такой JSON: { \"subject\": \"Это предмет\", \"quantity\": \"Это количество пользователей, которые тоже в этом предмете\""+
			"- Без лишнего текста, только JSON. Ни комментариев, ни пояснений.\n\n",
		string(jsons[users]),
	)

	GeminiText, err := sendTextToGemini(text)
	if err != nil {
		return nil, err
	}

	var GeminiTextParse = addBrackets(GeminiText[8 : len(GeminiText)-5])

	if GeminiTextParse == "[]" || len(GeminiTextParse) < 10 {
		return usersStatistic, errs.ErrNoUsersStatisticFound
	}

	if err := json.Unmarshal([]byte(GeminiTextParse), &usersStatistic); err != nil {
		logger.Error.Printf("[aiService.GetUsersStatistic] Error parsing gemini text to users: %v", err)
		return nil, err
	}

	return usersStatistic, err
}
