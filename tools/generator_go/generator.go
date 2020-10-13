package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	api2 "leaderboard/app/api"
	"leaderboard/app/leaderboard/services"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	const usernamePrefix = "user_"
	countries := []string{"GB", "US", "TR", "AU", "NZ"}

	properties, err := LoadProperties()
	if err != nil {
		log.Fatal(err)
	}

	redisService := buildRedisService(properties)
	userService := services.NewUserService(redisService, properties.LeaderboardKeyPrefix)
	numCPU := runtime.NumCPU() * 8

	var nUser uint64 = 0
	var wg sync.WaitGroup
	wg.Add(1_000_000)

	for i := 0; i < numCPU; i++ {
		go func() {
			for {
				profile := &api2.UserProfile{
					UserId:      "",
					DisplayName: fmt.Sprintf("%s%d", usernamePrefix, atomic.LoadUint64(&nUser)),
					Points:      float64(rand.Intn(10_000)),
					Rank:        0,
					Country:     countries[rand.Intn(len(countries))],
				}

				userId, err := userService.Create(profile)
				if err != nil {
					panic(err)
				}

				nPlays := rand.Intn(5)
				for j := 0; j < nPlays; j++ {
					err := submitScore(&api2.ScoreSubmission{
						Score:     float64(rand.Intn(10_000)),
						UserId:    userId,
						Timestamp: time.Now().Unix(),
					}, userService, redisService, properties.LeaderboardKeyPrefix)
					if err != nil {
						panic(err)
					}
				}

				atomic.AddUint64(&nUser, 1)
				if atomic.LoadUint64(&nUser)%1000 == 0 {
					log.Printf("@%d", atomic.LoadUint64(&nUser))
				}

				wg.Done()
			}
		}()
	}

	wg.Wait()
}

func submitScore(submission *api2.ScoreSubmission, userService *services.UserService, redisService *services.SingleRedisService, leaderboardKeyPrefix string) error {
	// send to redis
	user, err := userService.GetByID(submission.UserId)
	if err != nil {
		return err
	}

	globalBoard := fmt.Sprintf("%s%s", leaderboardKeyPrefix, "GLOBAL")
	localBoard := fmt.Sprintf("%s%s", leaderboardKeyPrefix, user.Country)

	redisService.Add(globalBoard, &redis.Z{
		Score:  submission.Score,
		Member: submission.UserId,
	})

	redisService.Add(localBoard, &redis.Z{
		Score:  submission.Score,
		Member: submission.UserId,
	})

	return nil
}

func buildRedisService(properties *Properties) *services.SingleRedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     properties.RedisHost,
		Password: properties.RedisPassword,
		DB:       properties.RedisDB,
	})

	return services.NewRedisService(client)
}
