package services

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type cacheService struct {
	cacheServer *redis.Client
}

type CacheService interface {
	GetKey(key string) (string, error)
	SetKey(key string, value string) error
	DeleteKey(key string) error
}

func CreateCacheClient(addr string, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return client
}

func NewCacheService(cacheClient *redis.Client) CacheService {
	return &cacheService{
		cacheServer: cacheClient,
	}
}

func (s *cacheService) GetKey(key string) (string, error) {
	ctx := context.Background()
	val, err := s.cacheServer.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist")
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (s *cacheService) SetKey(key string, value string) error {
	ctx := context.Background()
	err := s.cacheServer.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *cacheService) DeleteKey(key string) error {
	ctx := context.Background()
	err := s.cacheServer.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
