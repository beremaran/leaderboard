package services

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"leaderboard/app/api"
	"log"
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

	userProfile, err := api.SerializeUserProfile(profile)
	if err != nil {
		return "", err
	}

	us.redisService.Set(profile.UserId, userProfile)
	return profile.UserId, nil
}

func (us *UserService) GetByID(guid string) (*api.UserProfile, error) {
	profileEncoded, err := us.redisService.Get(guid)
	if err != nil {
		return nil, err
	}

	profile, err := api.DeserializeUserProfile(profileEncoded)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (us *UserService) GetByIDWithRank(guid string, leaderboardName string) (*api.UserProfile, error) {
	profile, err := us.GetByID(guid)
	if err != nil {
		return nil, err
	}

	us.SetRank(profile,leaderboardName)
	return profile, nil
}

func (us *UserService) SetRank(profile *api.UserProfile, leaderboardName string) {
	rank, err := us.redisService.GetRank(leaderboardName, profile.UserId)
	if err != nil {
		log.Fatal(err)
	}

	profile.Rank = rank
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
