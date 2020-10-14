package leaderboard

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"leaderboard/app/api"
	"leaderboard/app/leaderboard/handlers"
	"leaderboard/app/leaderboard/services"
	_ "leaderboard/docs"
	"log"
)

func Run() {
	properties, err := LoadProperties()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// validation
	e.Validator = services.NewStructValidator(validator.New())

	// services
	redisService := buildRedisService(properties)
	userService := services.NewUserService(redisService, properties.LeaderboardKeyPrefix)
	leaderboardService := services.NewLeaderboardService(userService, redisService, properties.LeaderboardKeyPrefix)

	// handlers
	userHandler := handlers.NewUserHandler(userService)
	userHandler.Register(e)

	leaderboardHandler := handlers.NewLeaderboardHandler(leaderboardService)
	leaderboardHandler.Register(e)

	scoreHandler := handlers.NewScoreHandler(userService, redisService, properties.LeaderboardKeyPrefix)
	scoreHandler.Register(e)

	actuator := handlers.NewActuatorHandler(redisService, userService)
	actuator.Register(e)

	e.Logger.Fatal(e.Start(":1323"))
}

func buildRedisService(properties *Properties) api.RedisService {
	var client redis.UniversalClient
	if properties.RedisCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{properties.RedisHost},
			PoolSize: 64,
			Password: properties.RedisPassword,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     properties.RedisHost,
			PoolSize: 64,
			Password: properties.RedisPassword,
		})
	}

	return services.NewRedisService(client, properties.LeaderboardKeyPrefix)
}
