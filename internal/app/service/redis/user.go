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
		logger.Error.Printf("[redisService.GetCachedUser] Error unmarshalling user: %v", err)

		return models.User{}, err
	}

	return user, nil
}

func GetCachedCharter(userID uint) (models.Charter, error) {
	cacheKey := fmt.Sprintf("userch:%d", userID)
	val, err := db.RedisClient.Get(context.Background(), cacheKey).Result()

	if err == redis.Nil {
		// В кэше нет — достаем из БД
		return models.Charter{}, errs.ErrUserNotFound
	} else if err != nil {
		return models.Charter{}, err
	}

	var charter models.Charter
	if err := json.Unmarshal([]byte(val), &charter); err != nil {
		logger.Error.Printf("[redisService.GetCachedCharter] Error unmarshalling user: %v", err)

		return models.Charter{}, err
	}

	return charter, nil
}

func SetCharterCache(charter models.Charter, userID uint) error {
	cacheKey := fmt.Sprintf("userch:%d", userID)
	data, err := json.Marshal(charter)
	if err != nil {
		logger.Error.Printf("[redisService.SetCharterCache] Error marshalling user: %v", err)

		return err
	}

	return db.RedisClient.Set(context.Background(), cacheKey, data, time.Minute*10).Err()
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
