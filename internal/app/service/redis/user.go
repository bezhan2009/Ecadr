package redisService

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/db"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func GetCachedUser(userID uint) (models.User, error) {
	cacheKey := fmt.Sprintf("user:%d", userID)
	val, err := db.RedisClient.Get(context.Background(), cacheKey).Result()

	if err == redis.Nil {
		// В кэше нет — достаем из БД
		return models.User{}, errs.ErrUserNotFound
	} else if err != nil {
		return models.User{}, err
	}

	var user models.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		logger.Error.Printf("[redisService.SetUserCache] Error unmarshalling user: %v", err)

		return models.User{}, err
	}

	return user, nil
}

func SetUserCache(user models.User) error {
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		logger.Error.Printf("[redisService.SetUserCache] Error marshalling user: %v", err)

		return err
	}

	return db.RedisClient.Set(context.Background(), cacheKey, data, time.Minute*10).Err()
}
