package services

import (
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"leaderboard/app/api"
)

type UserService struct {
	redisService         api.RedisService
	leaderboardKeyPrefix string
}

func NewUserService(redisService api.RedisService, leaderboardKeyPrefix string) *UserService {
	return &UserService{redisService: redisService, leaderboardKeyPrefix: leaderboardKeyPrefix}
}

func (us *UserService) Create(profile *api.UserProfile) (string, error) {
	if len(profile.UserId) == 0 {
		profile.UserId = uuid.New().String()
	}

	err := us.redisService.SetProfile(profile)
	if err != nil {
		return "", err
	}

	us.redisService.Add("GLOBAL", &redis.Z{
		Score:  profile.Points,
		Member: profile.UserId,
	})

	us.redisService.Add(profile.Country, &redis.Z{
		Score:  profile.Points,
		Member: profile.UserId,
	})

	return profile.UserId, nil
}

func (us *UserService) GetByID(guid string) (*api.UserProfile, error) {
	return us.redisService.GetProfile(guid)
}

func (us *UserService) GetByIDWithRank(guid string, leaderboardName string) (*api.UserProfile, error) {
	profile, err := us.GetByID(guid)
	if err != nil {
		return nil, err
	}

	err = us.SetRank(profile, leaderboardName)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (us *UserService) SetRank(profile *api.UserProfile, leaderboardName string) error {
	rank, err := us.redisService.GetRank(leaderboardName, profile.UserId)
	if err != nil {
		return err
	}

	profile.Rank = rank
	go func() {
		_ = us.redisService.SetProfile(profile)
	}()

	return nil
}

func (us *UserService) GetAllByID(guid ...string) ([]*api.UserProfile, error) {
	if len(guid) == 0 {
		return []*api.UserProfile{}, nil
	}

	var profiles []*api.UserProfile
	for _, id := range guid {
		byID, err := us.GetByID(id)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, byID)
	}

	return profiles, nil
}
