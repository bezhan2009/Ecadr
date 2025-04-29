package aiSerivce

import (
	"Ecadr/internal/app/models"
	"encoding/json"
	"fmt"
)

func SendMessageToGeminiAI(messages []models.Message, kinderNote []models.KindergartenNote,
	schoolGrade []models.SchoolGrade,
	achievementsUser []models.Achievement,
	msg string) (respAI *string, err error) {

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

	msgJson, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}

	text := fmt.Sprintf(
		`Ниже приведены данные о пользователе:

1. Успеваемость в детском саду: %s
2. Успеваемость в школе: %s
3. Достижения: %s

4. История сообщений: %s

Теперь твоя задача — ответить на вопрос от самого пользователя.
Обрати внимание: данные выше просто описывают пользователя, чтобы ты лучше его понимал.

⚠️ Очень важно: ответь пользователю **на том языке, на котором он задал вопрос**.

Вот сам вопрос:
%s`,
		string(jsons[kinderNotes]),
		string(jsons[schoolGrades]),
		string(jsons[achievements]),
		string(msgJson),
		msg,
	)

	fmt.Println(text)

	GeminiText, err := sendTextToGemini(text)
	if err != nil {
		return nil, err
	}

	return &GeminiText, err
}
