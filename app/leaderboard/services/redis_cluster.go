package services

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type ClusterRedisService struct {
	context context.Context
	client  *redis.ClusterClient
}

func NewClusterRedisService(client *redis.ClusterClient) *ClusterRedisService {
	return &ClusterRedisService{client: client, context: context.Background()}
}

func (o *ClusterRedisService) Set(key string, value string) {
	o.client.Set(o.context, key, value, 8*time.Hour)
}

func (o *ClusterRedisService) Get(key string) (string, error) {
	result, err := o.client.Get(o.context, key).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (o *ClusterRedisService) Add(sortedSetName string, z ...*redis.Z) {
	o.client.ZAdd(o.context, sortedSetName, z...)
}

func (o *ClusterRedisService) FlushAll() {
	o.client.FlushAll(o.context)
}

func (o *ClusterRedisService) GetSortedSetSize(sortedSetName string) (int64, error) {
	result, err := o.client.ZCard(o.context, sortedSetName).Result()
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (o *ClusterRedisService) GetRank(sortedSetName string, key string) (int64, error) {
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

func (o *ClusterRedisService) GetScore(sortedSetName string, key string) (float64, error) {
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

func (o *ClusterRedisService) GetPage(sortedSetName string, startIndex int64, endIndex int64) ([]redis.Z, error) {
	result, err := o.client.ZRevRangeWithScores(o.context, sortedSetName, startIndex, endIndex).Result()
	if err != nil {
		return []redis.Z{}, err
	}

	return result, nil
}
