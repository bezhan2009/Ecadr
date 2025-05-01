package utils

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
	"strings"
)

const (
	Users            = "Users"
	KinderNotes      = "kinder_note"
	SchoolGrades     = "school_grades"
	Achievements     = "Achievements"
	Courses          = "Courses"
	Vacancies        = "Vacancies"
	CompanyInfoStats = "company_info_stats"
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

func AddBrackets(text string) string {
	trimmedText := strings.TrimSpace(text)

	if !strings.HasPrefix(trimmedText, "[") {
		trimmedText = "[" + trimmedText
	}
	if !strings.HasSuffix(trimmedText, "]") {
		trimmedText = trimmedText + "]"
	}

	return trimmedText
}

func SerializeData(usersAnalyse []models.User,
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
			logger.Error.Printf("[aiService.SerializeData] Error marshalling kindernote\n\tKinderNote:%v\n\tError: %v", kinderNote, err)
			return nil, err
		}

		jsonSchoolGrades, err = json.Marshal(schoolGrade)
		if err != nil {
			logger.Error.Printf("[aiService.SerializeData] Error marshalling SchoolGrades\n\tSchoolGrades:%v\n\tError: %v", SchoolGrades, err)
			return nil, err
		}

		jsonAchievements, err = json.Marshal(achievementsUser)
		if err != nil {
			logger.Error.Printf("[aiService.SerializeData] Error marshalling Achievements\n\tAchievements:%v\n\tError: %v", Achievements, err)
			return nil, err
		}
	}

	jsonUsers, err = json.Marshal(usersAnalyse)
	if err != nil {
		logger.Error.Printf("[aiService.GetAnalyseForCourseUser] Error marshalling kindernote\n\tUsers:%v\n\tError: %v", Users, err)
		return nil, err
	}

	jsonVacancies, err = json.Marshal(vacanciesWorker)
	if err != nil {
		logger.Error.Printf("[aiService.SerializeData] Error marshalling Vacancies\n\tCourses:%v\n\tError: %v", Courses, err)
		return nil, err
	}

	jsonCourses, err = json.Marshal(coursesWorker)
	if err != nil {
		logger.Error.Printf("[aiService.SerializeData] Error marshalling Courses\n\tCourses:%v\n\tError: %v", Courses, err)
		return nil, err
	}

	jsonCompanyInfo, err = json.Marshal(companyInfo)
	if err != nil {
		logger.Error.Printf("[aiService.SerializeData] Error while marshalling company info stats\n\tcompany info:%v\n\tError: %v", companyInfo, err)
		return nil, err
	}

	jsons[Users] = jsonUsers
	jsons[CompanyInfoStats] = jsonCompanyInfo
	jsons[KinderNotes] = jsonKinderNote
	jsons[SchoolGrades] = jsonSchoolGrades
	jsons[Achievements] = jsonAchievements
	jsons[Courses] = jsonCourses
	jsons[Vacancies] = jsonVacancies

	return jsons, nil
}

func SendTextToGemini(text string) (response string, err error) {
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

	fmt.Println(analyse)

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

	fmt.Println(GeminiResp)

	if len(GeminiResp.Candidates) == 0 {
		logger.Error.Printf("[aiService.GetAnalyseForCourseUser] No candidates returned from Gemini response: %s", string(body))
		return "", errs.ErrGeminiIsNotWorking
	}
	if len(GeminiResp.Candidates[0].Content.Parts) == 0 {
		logger.Error.Printf("[aiService.GetAnalyseForCourseUser] No parts in first candidate's content: %s", string(body))
		return "", errs.ErrGeminiIsNotWorking
	}

	var GeminiText = GeminiResp.Candidates[0].Content.Parts[0].Text

	return GeminiText, nil
}
