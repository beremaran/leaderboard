package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"leaderboard/app/api"
	"leaderboard/app/leaderboard/services"
	"net/http"
)

type ScoreHandler struct {
	userService          *services.UserService
	redisService         api.RedisService
	leaderboardKeyPrefix string
}

func NewScoreHandler(userService *services.UserService, redisService api.RedisService, leaderboardKeyPrefix string) *ScoreHandler {
	return &ScoreHandler{userService: userService, redisService: redisService, leaderboardKeyPrefix: leaderboardKeyPrefix}
}

func (s *ScoreHandler) Register(echo *echo.Echo) {
	group := echo.Group("/score")

	group.POST("/submit", s.submit)
}

func (s *ScoreHandler) submit(c echo.Context) (err error) {
	submission := new(api.ScoreSubmission)
	if err = c.Bind(submission); err != nil {
		return
	}

	if err = c.Validate(submission); err != nil {
		return c.JSON(http.StatusBadRequest, api.NewValidationErrorResponse(err.Error()))
	}

	// send to redis
	user, err := s.userService.GetByID(submission.UserId)
	if err != nil {
		return echo.ErrNotFound
	}

	s.redisService.Add("GLOBAL", &redis.Z{
		Score:  submission.Score,
		Member: submission.UserId,
	})

	s.redisService.Add(user.Country, &redis.Z{
		Score:  submission.Score,
		Member: submission.UserId,
	})

	return c.NoContent(http.StatusCreated)
}
