package services

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"leaderboard/app/api"
	"sync"
	"time"
)

type RedisService struct {
	context              context.Context
	client               redis.UniversalClient
	leaderboardKeyPrefix string
	leaderboardKeys      map[string]string
	leaderboardKeysMux   sync.Mutex
}

func NewRedisService(client redis.UniversalClient, leaderboardKeyPrefix string) *RedisService {
	return &RedisService{client: client, context: context.Background(), leaderboardKeyPrefix: leaderboardKeyPrefix}
}

func (o *RedisService) HSet(key string, values ...interface{}) *redis.IntCmd {
	return o.client.HSet(o.context, key, values...)
}

func (o *RedisService) HGetAll(key string) *redis.StringStringMapCmd {
	return o.client.HGetAll(o.context, key)
}

func (o *RedisService) SetProfile(profile *api.UserProfile) (err error) {
	_, err = o.client.HSet(
		o.context, profile.UserId,
		"display_name", profile.DisplayName,
		"country", profile.Country,
		"points", profile.Points,
	).Result()

	return
}

func (o *RedisService) GetProfile(id string) (*api.UserProfile, error) {
	resultMap, err := o.client.HGetAll(o.context, id).Result()
	if err != nil {
		return nil, err
	}

	profile := new(api.UserProfile)
	profile.UserId = id
	profile.DisplayName = resultMap["display_name"]
	profile.Country = resultMap["country"]
	profile.Points, _ = o.GetScore("GLOBAL", id)

	return profile, nil
}

func (o *RedisService) Set(key string, value string) {
	o.client.Set(o.context, key, value, 8*time.Hour)
}

func (o *RedisService) Get(key string) (string, error) {
	result, err := o.client.Get(o.context, key).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (o *RedisService) getBoardKey(name string) string {
	o.leaderboardKeysMux.Lock()
	defer o.leaderboardKeysMux.Unlock()

	if o.leaderboardKeys == nil {
		o.leaderboardKeys = map[string]string{}
	}

	var boardKey string
	if val, ok := o.leaderboardKeys[name]; ok {
		return val
	}

	boardKey = fmt.Sprintf("%s%s", o.leaderboardKeyPrefix, name)
	o.leaderboardKeys[name] = boardKey

	return boardKey
}

func (o *RedisService) Add(sortedSetName string, z ...*redis.Z) {
	o.client.ZAdd(o.context, o.getBoardKey(sortedSetName), z...)
}

func (o *RedisService) FlushAll() {
	o.client.FlushAll(o.context)
}

func (o *RedisService) GetSortedSetSize(sortedSetName string) (int64, error) {
	result, err := o.client.ZCard(o.context, o.getBoardKey(sortedSetName)).Result()
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (o *RedisService) GetRank(sortedSetName string, key string) (int64, error) {
	result, err := o.client.ZRevRank(o.context, o.getBoardKey(sortedSetName), key).Result()
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

func (o *RedisService) GetScore(sortedSetName string, key string) (float64, error) {
	result, err := o.client.ZScore(o.context, o.getBoardKey(sortedSetName), key).Result()
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

func (o *RedisService) GetPage(sortedSetName string, startIndex int64, endIndex int64) ([]redis.Z, error) {
	result, err := o.client.ZRevRangeWithScores(o.context, o.getBoardKey(sortedSetName), startIndex, endIndex).Result()
	if err != nil {
		return []redis.Z{}, err
	}

	return result, nil
}
