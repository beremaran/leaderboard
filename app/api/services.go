package api

import (
	"github.com/go-redis/redis/v8"
)

type RedisService interface {
	Set(key string, value string)
	Get(key string) (string, error)
	Add(sortedSetName string, z ...*redis.Z)
	FlushAll()
	GetSortedSetSize(sortedSetName string) (int64, error)
	GetRank(sortedSetName string, key string) (int64, error)
	GetScore(sortedSetName string, key string) (float64, error)
	GetPage(sortedSetName string, startIndex int64, endIndex int64) ([]redis.Z, error)
	GetProfile(id string) (*UserProfile, error)
	SetProfile(profile *UserProfile) (err error)
	HSet(key string, values ...interface{}) *redis.IntCmd
	HGetAll(key string) *redis.StringStringMapCmd
	Exists(key string) (bool, error)
	GetOrDefault(key string, defaultVal string) string
}

type LeaderboardService interface {
	GetPage(boardName string, page int64, pageSize int64) ([]*LeaderboardRow, error)
}
