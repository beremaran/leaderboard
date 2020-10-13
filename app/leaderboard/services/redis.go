package services

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type SingleRedisService struct {
	context context.Context
	client  *redis.Client
}

func NewSingleRedisService(client *redis.Client) *SingleRedisService {
	return &SingleRedisService{client: client, context: context.Background()}
}

func (o *SingleRedisService) Set(key string, value string) {
	o.client.Set(o.context, key, value, 8*time.Hour)
}

func (o *SingleRedisService) Get(key string) (string, error) {
	result, err := o.client.Get(o.context, key).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (o *SingleRedisService) Add(sortedSetName string, z ...*redis.Z) {
	o.client.ZAdd(o.context, sortedSetName, z...)
}

func (o *SingleRedisService) FlushAll() {
	o.client.FlushAll(o.context)
}

func (o *SingleRedisService) GetSortedSetSize(sortedSetName string) (int64, error) {
	result, err := o.client.ZCard(o.context, sortedSetName).Result()
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (o *SingleRedisService) GetRank(sortedSetName string, key string) (int64, error) {
	result, err := o.client.ZRevRank(o.context, sortedSetName, key).Result()
	if err != nil {
		o.Add(sortedSetName, &redis.Z{
			Score:  0,
			Member: key,
		})

		rank, err := o.GetRank(sortedSetName, key)
		if err != nil {
			return 0.0, err
		}

		return rank, nil
	}

	return result + 1, nil
}

func (o *SingleRedisService) GetScore(sortedSetName string, key string) (float64, error) {
	result, err := o.client.ZScore(o.context, sortedSetName, key).Result()
	if err != nil {
		o.Add(sortedSetName, &redis.Z{
			Score:  0,
			Member: key,
		})

		score, err := o.GetScore(sortedSetName, key)
		if err != nil {
			return 0.0, err
		}

		return score, nil
	}

	return result, nil
}

func (o *SingleRedisService) GetPage(sortedSetName string, startIndex int64, endIndex int64) ([]redis.Z, error) {
	result, err := o.client.ZRevRangeWithScores(o.context, sortedSetName, startIndex, endIndex).Result()
	if err != nil {
		return []redis.Z{}, err
	}

	return result, nil
}
