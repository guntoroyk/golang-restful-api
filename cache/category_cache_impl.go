package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/guntoroyk/golang-restful-api/model/domain"
	"log"
	"time"
)

type CategoryCacheImpl struct {
	redisClient *redis.Client
}

func NewCategoryCache(redisClient *redis.Client) CategoryCache {
	return &CategoryCacheImpl{redisClient: redisClient}
}

func (c CategoryCacheImpl) SetCategory(ctx context.Context, category domain.Category) error {
	key := fmt.Sprintf("category:%d", category.Id)
	expiration := 5 * time.Minute

	json, err := json.Marshal(category)
	if err != nil {
		log.Println("problem marshalling cache, err: ", err.Error())
		return err
	}

	err = c.redisClient.Set(ctx, key, json, expiration).Err()
	if err != nil {
		log.Println("problem saving cache, err: ", err.Error())
		return err
	}

	return nil
}

func (c CategoryCacheImpl) SetCategoryBatch(ctx context.Context, categories []domain.Category) error {
	json, err := json.Marshal(categories)
	expiration := 5 * time.Minute

	if err != nil {
		log.Println("problem marshalling cache, err: ", err.Error())
		return err
	}

	err = c.redisClient.Set(ctx, "categories", json, expiration).Err()
	if err != nil {
		log.Println("problem saving cache, err: ", err.Error())
		return err
	}

	return nil
}

func (c CategoryCacheImpl) Delete(ctx context.Context, category domain.Category) error {
	key := fmt.Sprintf("category:%d", category.Id)

	err := c.redisClient.Del(ctx, key).Err()
	if err != nil {
		log.Println("problem deleting cache, err: ", err.Error())
		return err
	}
	return nil
}

func (c CategoryCacheImpl) GetCategory(ctx context.Context, categoryId int) (domain.Category, error) {
	var resp domain.Category
	key := fmt.Sprintf("category:%d", categoryId)

	val, err := c.redisClient.Get(ctx, key).Result()

	if err != nil {
		return resp, err
	}

	if err := json.Unmarshal([]byte(val), &resp); err != nil {
		log.Println("problem unmarshalling cache, err: ", err.Error())
		return resp, err
	}

	return resp, nil
}

func (c CategoryCacheImpl) GetCategoryBatch(ctx context.Context) ([]domain.Category, error) {
	var resp []domain.Category

	val, err := c.redisClient.Get(ctx, "categories").Result()

	fmt.Printf("error: %s", err)
	fmt.Printf("error text: %s", err.Error())

	if err != nil {
		return resp, err
	}

	if err := json.Unmarshal([]byte(val), &resp); err != nil {
		log.Println("problem unmarshalling cache, err: ", err.Error())
		return resp, err
	}

	return resp, nil

}
