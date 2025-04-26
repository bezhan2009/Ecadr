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
	"strings"
)

const (
	users            = "users"
	kinderNotes      = "kinder_note"
	schoolGrades     = "school_grades"
	achievements     = "achievements"
	courses          = "courses"
	vacancies        = "vacancies"
	companyInfoStats = "company_info_stats"
)

func cleanGeminiText(text string) string {
	const prefix = "```json\n"
	const suffix = "```"

	// Убираем префикс
	if len(text) >= len(prefix) && text[:len(prefix)] == prefix {
		text = text[len(prefix):]
	}

	// Убираем суффикс
	if len(text) >= len(suffix) && text[len(text)-len(suffix):] == suffix {
		text = text[:len(text)-len(suffix)]
	}

	return text
}

func addBrackets(text string) string {
	trimmedText := strings.TrimSpace(text)

	if !strings.HasPrefix(trimmedText, "[") {
		trimmedText = "[" + trimmedText
	}
	if !strings.HasSuffix(trimmedText, "]") {
		trimmedText = trimmedText + "]"
	}

	return trimmedText
}

func serializeData(usersAnalyse []models.User,
	companyInfo []models.CompanyInfo,
	kinderNote []models.KindergartenNote,
	schoolGrade []models.SchoolGrade,
	achievementsUser []models.Achievement,
	coursesWorker []models.Course,
	vacanciesWorker []models.Vacancy) (jsons map[string][]byte, err error) {
	jsons = make(map[string][]byte)

	var (
		jsonKinderNote   []byte
		jsonSchoolGrades []byte
		jsonAchievements []byte
		jsonCourses      []byte
		jsonVacancies    []byte
		jsonUsers        []byte
		jsonCompanyInfo  []byte
	)

	if usersAnalyse == nil {
		jsonKinderNote, err = json.Marshal(kinderNote)
		if err != nil {
			logger.Error.Printf("[aiService.serializeData] Error marshalling kindernote\n\tKinderNote:%v\n\tError: %v", kinderNote, err)
			return nil, err
		}

		jsonSchoolGrades, err = json.Marshal(schoolGrade)
		if err != nil {
			logger.Error.Printf("[aiService.serializeData] Error marshalling SchoolGrades\n\tSchoolGrades:%v\n\tError: %v", schoolGrades, err)
			return nil, err
		}

		jsonAchievements, err = json.Marshal(achievementsUser)
		if err != nil {
			logger.Error.Printf("[aiService.serializeData] Error marshalling achievements\n\tachievements:%v\n\tError: %v", achievements, err)
			return nil, err
		}
	}

	jsonUsers, err = json.Marshal(usersAnalyse)
	if err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForCourseUser] Error marshalling kindernote\n\tUsers:%v\n\tError: %v", users, err)
		return nil, err
	}

	jsonVacancies, err = json.Marshal(vacanciesWorker)
	if err != nil {
		logger.Error.Printf("[aiService.serializeData] Error marshalling vacancies\n\tcourses:%v\n\tError: %v", courses, err)
		return nil, err
	}

	jsonCourses, err = json.Marshal(coursesWorker)
	if err != nil {
		logger.Error.Printf("[aiService.serializeData] Error marshalling courses\n\tcourses:%v\n\tError: %v", courses, err)
		return nil, err
	}

	jsonCompanyInfo, err = json.Marshal(companyInfo)
	if err != nil {
		logger.Error.Printf("[aiService.serializeData] Error while marshalling company info stats\n\tcompany info:%v\n\tError: %v", companyInfo, err)
		return nil, err
	}

	jsons[users] = jsonUsers
	jsons[companyInfoStats] = jsonCompanyInfo
	jsons[kinderNotes] = jsonKinderNote
	jsons[schoolGrades] = jsonSchoolGrades
	jsons[achievements] = jsonAchievements
	jsons[courses] = jsonCourses
	jsons[vacancies] = jsonVacancies

	return jsons, nil
}

func sendTextToGemini(text string) (response string, err error) {
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
		logger.Error.Printf("[aiService.GetAnalyseForUserCourse] Error marshalling json body: %v", err)
		return "", err
	}

	var geminiURL = fmt.Sprintf(
		"%s?key=%s",
		os.Getenv("GEMINI_AI_API"),
		os.Getenv("GEMINI_API_KEY"),
	)

	analyse, err := http.Post(
		geminiURL,
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForUserCourse] Error getting analyse: %v", err)
		return "", err
	}
	defer analyse.Body.Close()

	body, err := ioutil.ReadAll(analyse.Body)
	if err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForUserCourse] Error reading body analyse: %v", err)
		return "", err
	}

	var GeminiResp models.Gemini
	if err := json.Unmarshal(body, &GeminiResp); err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForCourseUser] Error parsing body: %v", err)
		return "", err
	}

	var GeminiText = GeminiResp.Candidates[0].Content.Parts[0].Text

	return GeminiText, nil
}
