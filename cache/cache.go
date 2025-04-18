package cache

import (
	"context"
	"fmt"
	"os"
	bt "restapi/basic_types"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Cache *redis.Client
	ctx   context.Context
}

func NewRedisCache() (*RedisCache, error) {
	cache := &RedisCache{}
	cache.ctx = context.Background()

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if err := client.Ping(cache.ctx).Err(); err != nil {
		return nil, err
	}

	cache.Cache = client
	return cache, nil
}

func (rc *RedisCache) Set(task *bt.Task) error {
	id := strconv.Itoa(task.ID)

	exists, err := rc.Cache.HExists(rc.ctx, id, "name").Result()
	if err != nil {
		return fmt.Errorf("failed to check if task %d exists in cache: %v", task.ID, err)
	}
	if exists {
		return fmt.Errorf("task %d already exists in cache", task.ID)
	}

	err = rc.Cache.HSet(rc.ctx, id, map[string]interface{}{
		"name":        task.Name,
		"description": task.Description,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to insert task %d into cache: %v", task.ID, err)
	}

	err = rc.Cache.Expire(rc.ctx, id, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set expired time of task %d in cache: %v", task.ID, err)
	}

	return nil
}

func (rc *RedisCache) Get(taskID int) (*bt.Task, error) {
	id := strconv.Itoa(taskID)

	data, err := rc.Cache.HGetAll(rc.ctx, id).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get task %d from cache: %v", taskID, err)
	}
	if len(data) == 0 {
		return nil, ErrTaskNotFound
	}

	task := &bt.Task{
		ID:          taskID,
		Name:        data["name"],
		Description: data["description"],
	}

	return task, nil
}

func (rc *RedisCache) Delete(taskID int) error {
	id := strconv.Itoa(taskID)

	removed, err := rc.Cache.Del(rc.ctx, id).Result()
	if err != nil {
		return fmt.Errorf("failed to delete task %d from cache: %v", taskID, err)
	}

	if removed == 0 {
		return ErrTaskNotFound
	}

	return nil
}
