package ai

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	aiSerivce "Ecadr/internal/app/service/ai"
	"Ecadr/internal/controllers/middlewares"
	"Ecadr/internal/security"
	"Ecadr/pkg/db"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func getAIRecommendsFromRedis(key string) ([]models.AiUserRecommends, error) {
	aiRecsStr, err := db.GetCache(key)
	if aiRecsStr != "" {
		var aiUserRecs []models.AiUserRecommends
		err := json.Unmarshal([]byte(aiRecsStr), &aiUserRecs)
		if err != nil {
			db.DeleteCache(key)
		} else {
			return aiUserRecs, nil
		}
	}

	return nil, err
}

// GetAIRecommendsForUser godoc
// @Summary Рекомендаций от ИИ пользователю
// @Description Возвращает список рекомендаций для улучшения на основе данных пользователя (оценки, достижения и т.д.)
// @Tags AI
// @Accept  json
// @Produce  json
// @Success 201 {object} []models.AiUserRecommends "Рекомендаций на основе анализа AI"
// @Failure 400 {object} errs.ErrorResp "Неверный запрос или пользователь не найден"
// @Failure 500 {object} errs.ErrorResp "Внутренняя ошибка сервера"
// @Router /ai/recommends [get]
// @Security ApiKeyAuth
func GetAIRecommendsForUser(c *gin.Context) (interface{}, error) {
	userID := c.GetUint(middlewares.UserIDCtx)

	user, _, err := service.GetUserByID(userID)
	if err != nil {
		//controllers.HandleError(c, err)
		return nil, err
	}

	keyCacheAIRecs := fmt.Sprintf("ai_recs:%d", userID)

	aiRecommendsFromRedis, err := getAIRecommendsFromRedis(keyCacheAIRecs)
	if err == nil && len(aiRecommendsFromRedis) > 0 {
		//c.JSON(http.StatusOK, aiRecommendsFromRedis)
		return aiRecommendsFromRedis, nil
	}

	aiRecommends, err := aiSerivce.GetRecommendsForUser(
		user.KindergartenNotes,
		user.SchoolGrades,
		user.Achievements,
	)
	if err != nil {
		//controllers.HandleError(c, err)
		return nil, err
	}

	aiRecommendsJson, err := json.Marshal(aiRecommends)
	if err == nil {
		db.SetCache(
			keyCacheAIRecs,
			aiRecommendsJson,
			time.Duration(security.AppSettings.RedisParams.TTLMinutes)*time.Minute,
		)
	}

	//c.JSON(http.StatusCreated, aiRecommends)
	return aiRecommends, nil
}
