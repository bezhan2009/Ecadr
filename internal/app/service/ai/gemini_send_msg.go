package aiSerivce

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func SendMessageToGeminiAI(kinderNote []models.KindergartenNote,
	schoolGrades []models.SchoolGrade,
	achievements []models.Achievement,
	msg string) (respAI *string, err error) {

	jsonKinderNote, err := json.Marshal(kinderNote)
	if err != nil {
		logger.Error.Printf("[aiService.SendMessageToGeminiAI] Error marshalling kindernote\n\tKinderNote:%v\n\tError: %v", kinderNote, err)
		return nil, err
	}

	jsonSchoolGrades, err := json.Marshal(schoolGrades)
	if err != nil {
		logger.Error.Printf("[aiService.SendMessageToGeminiAI] Error marshalling SchoolGrades\n\tSchoolGrades:%v\n\tError: %v", schoolGrades, err)
		return nil, err
	}

	jsonAchievements, err := json.Marshal(achievements)
	if err != nil {
		logger.Error.Printf("[aiService.SendMessageToGeminiAI] Error marshalling achievements\n\tachievements:%v\n\tError: %v", achievements, err)
		return nil, err
	}

	text := fmt.Sprintf(
		`Ниже приведены данные о пользователе:

1. Успеваемость в детском саду: %s
2. Успеваемость в школе: %s
3. Достижения: %s

Теперь твоя задача — ответить на вопрос от самого пользователя.
Обрати внимание: данные выше просто описывают пользователя, чтобы ты лучше его понимал.

⚠️ Очень важно: ответь пользователю **на том языке, на котором он задал вопрос**.

Вот сам вопрос:
%s`,
		string(jsonKinderNote),
		string(jsonSchoolGrades),
		string(jsonAchievements),
		msg,
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
		logger.Error.Printf("[aiService.SendMessageToGeminiAI] Error marshalling json body: %v", err)
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
		logger.Error.Printf("[aiService.SendMessageToGeminiAI] Error getting analyse: %v", err)
		return nil, err
	}
	defer analyse.Body.Close()

	body, err := ioutil.ReadAll(analyse.Body)
	if err != nil {
		logger.Error.Printf("[aiService.SendMessageToGeminiAI] Error reading body analyse: %v", err)
		return nil, err
	}

	var GeminiResp models.Gemini
	if err := json.Unmarshal(body, &GeminiResp); err != nil {
		logger.Error.Printf("[aiService.SendMessageToGeminiAI] Error parsing body: %v", err)
		return nil, err
	}

	var GeminiText = GeminiResp.Candidates[0].Content.Parts[0].Text

	return &GeminiText, err
}
