package api

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type UserService interface {
	GetByUUID(uuid string) UserProfile
	Create(profile UserProfile)
}

type LeaderboardService interface {
	GetGlobalLeaderboard(query LeaderboardQuery) []LeaderboardRow
	GetCountryLeaderboard(query LeaderboardQuery) []LeaderboardRow
	SubmitScore(submission ScoreSubmission)
}

type ActuatorService interface {
	GenerateRecords(n int)
	EvictCache()
}

type Handler interface {
	Register(echo *echo.Echo)
}

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
}
