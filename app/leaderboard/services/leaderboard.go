package services

import (
	"leaderboard/app/api"
)

type LeaderboardService struct {
	userService          *UserService
	redisService         api.RedisService
	leaderboardKeyPrefix string
}

func NewLeaderboardService(userService *UserService, redisService api.RedisService, leaderboardKeyPrefix string) *LeaderboardService {
	return &LeaderboardService{userService: userService, redisService: redisService, leaderboardKeyPrefix: leaderboardKeyPrefix}
}

func (ls *LeaderboardService) GetPage(boardName string, page int64, pageSize int64) ([]*api.LeaderboardRow, error) {
	rankingTuples, err := ls.redisService.GetPage(
		boardName,
		(page-1)*pageSize, page*pageSize-1,
	)
	if err != nil {
		return nil, err
	}

	var userIds []string
	scoreMap := make(map[string]float64)
	for _, t := range rankingTuples {
		scoreMap[t.Member.(string)] = t.Score
		userIds = append(
			userIds,
			t.Member.(string),
		)
	}

	profiles, err := ls.userService.GetAllByID(userIds...)
	if err != nil {
		return nil, err
	}

	var rows []*api.LeaderboardRow
	for _, profile := range profiles {
		rank, err := ls.redisService.GetRank(boardName, profile.UserId)
		if err != nil {
			return nil, err
		}

		rows = append(rows, &api.LeaderboardRow{
			Rank:        rank,
			Points:      int64(scoreMap[profile.UserId]),
			DisplayName: profile.DisplayName,
			Country:     profile.Country,
		})
	}

	return rows, nil
}
