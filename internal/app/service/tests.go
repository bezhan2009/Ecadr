package service

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/repository"
	"sort"
)

func SortSubmissionsByScore(testID uint) ([]models.TestSubmission, error) {
	subs, err := repository.GetAllSubmissions(testID)
	if err != nil {
		return nil, err
	}

	for i := range subs {
		total := len(subs[i].Answers)
		correct := 0
		for _, ans := range subs[i].Answers {
			for _, sel := range ans.SelectedChoiceIDs {
				for _, ch := range ans.Question.Choices {
					if int64(ch.ID) == sel && ch.IsCorrect {
						correct++
					}
				}
			}
		}
		if total > 0 {
			subs[i].CorrectPercentage = float64(correct) / float64(total) * 100
		}
	}
	sort.SliceStable(subs, func(i, j int) bool {
		return subs[i].CorrectPercentage > subs[j].CorrectPercentage
	})
	return subs, nil
}

func GetTestsByTypeAndID(targetType string, targetID uint) (tests []models.Test, err error) {
	tests, err = repository.GetTasksByTypeAndID(targetType, targetID)
	if err != nil {
		return nil, err
	}

	return tests, nil
}

func GetTestByID(testID uint) (test models.Test, err error) {
	test, err = repository.GetTestByID(testID)
	if err != nil {
		return test, err
	}

	return test, nil
}

func CreateTest(test models.Test) (testID uint, err error) {
	testID, err = repository.CreateTest(test)
	if err != nil {
		return testID, err
	}

	return testID, nil
}

func UpdateTest(test models.Test) (testID uint, err error) {
	testID, err = repository.UpdateTest(test)
	if err != nil {
		return testID, err
	}

	return testID, nil
}

func DeleteTest(testID uint) error {
	err := repository.DeleteTest(testID)
	if err != nil {
		return err
	}

	return nil
}

func CreateTestSubmission(submission models.TestSubmission) (err error) {
	err = repository.CreateSubmission(submission)
	if err != nil {
		return err
	}

	return nil
}
