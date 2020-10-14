package tasks

import (
	"fmt"
	"github.com/google/uuid"
	"leaderboard/app/api"
	"leaderboard/app/leaderboard/services"
	"log"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"
)

const KeyTask = "TASK_GU"
const KeyRunning = "TASK_GU_RUN_STATE"
const KeyStopFlag = "FLAG_STOP_TASK_GU"
const FieldStatus = "STATUS"
const FieldCompleted = "COMPLETED"
const FieldConcurrency = "CONCURRENCY"
const FieldStartedAt = "STARTED_AT"
const FieldUserCount = "USER_COUNT"

type GenerateUsersSingletonTask struct {
	userService  *services.UserService
	redisService api.RedisService
}

func NewGenerateUsersSingletonTask(userService *services.UserService, redisService api.RedisService) *GenerateUsersSingletonTask {
	return &GenerateUsersSingletonTask{userService: userService, redisService: redisService}
}

func (g *GenerateUsersSingletonTask) Start(nUsers uint64, maxConcurrency int64) (*api.GenerateUserTaskStatus, error) {
	running, _ := g.redisService.Get(KeyRunning)
	if running == "t" {
		log.Printf("already running")
		return g.Status()
	}

	log.Printf("starting generator ...")
	redisInit := make(chan bool, 1)
	go g.generate(nUsers, maxConcurrency, redisInit)

	<-redisInit
	return g.Status()
}

func (g *GenerateUsersSingletonTask) Stop() error {
	g.redisService.Set(KeyStopFlag, "1")
	g.redisService.Set(KeyRunning, "f")

	return nil
}

func (g *GenerateUsersSingletonTask) Status() (*api.GenerateUserTaskStatus, error) {
	statusMap, err := g.redisService.HGetAll(KeyTask).Result()
	if err != nil {
		return nil, err
	}

	completedRatio, err := strconv.ParseFloat(statusMap[FieldCompleted], 64)
	if err != nil {
		return nil, err
	}

	concurrency, err := strconv.ParseInt(statusMap[FieldConcurrency], 10, 64)
	if err != nil {
		return nil, err
	}

	userCount, err := strconv.ParseInt(statusMap[FieldUserCount], 10, 64)
	if err != nil {
		return nil, err
	}

	return &api.GenerateUserTaskStatus{
		Status:      statusMap[FieldStatus],
		Completed:   completedRatio,
		Concurrency: concurrency,
		StartedAt:   statusMap[FieldStartedAt],
		UserCount:   userCount,
	}, nil
}

func (g *GenerateUsersSingletonTask) generate(n uint64, maxConcurrency int64, redisInit chan bool) {
	g.redisService.Set(KeyStopFlag, "0")
	g.redisService.Set(KeyRunning, "t")
	countries := []string{"TR", "US", "GB", "CN", "JP", "AU", "NZ"}

	var stop int32 = 0
	var generated uint64 = 0

	g.redisService.HSet(
		KeyTask,
		FieldCompleted,
		strconv.FormatFloat(float64(generated)*100/float64(n), 'f', 2, 64),
		FieldConcurrency,
		strconv.FormatInt(maxConcurrency, 10),
		FieldUserCount,
		strconv.FormatInt(int64(generated), 10),
		FieldStartedAt,
		time.Now().String(),
		FieldStatus,
		"STARTING",
	)
	redisInit <- true

	var cpu int64
	for cpu = 0; cpu < maxConcurrency; cpu++ {
		go func() {
			for atomic.LoadInt32(&stop) == 0 && atomic.LoadUint64(&generated) < n {
				id := uuid.New().String()
				_, err := g.userService.Create(&api.UserProfile{
					UserId:      id,
					DisplayName: fmt.Sprintf("user_%d_%s", atomic.LoadUint64(&generated), id),
					Points:      rand.Float64() * 100_000,
					Rank:        0,
					Country:     countries[rand.Intn(len(countries))],
				})

				if err != nil {
					panic(err)
				}
				atomic.AddUint64(&generated, 1)
			}

			log.Printf("generator goroutine is done.")
		}()
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		updateStatus := func() {
			log.Printf("updating task status")
			var taskStatus = "RUNNING"
			if atomic.LoadInt32(&stop) == 1 {
				taskStatus = "STOPPED"
			}

			// update stats
			g.redisService.HSet(
				KeyTask,
				FieldCompleted,
				strconv.FormatFloat(float64(generated)*100/float64(n), 'f', 2, 64),
				FieldConcurrency,
				strconv.FormatInt(maxConcurrency, 10),
				FieldUserCount,
				strconv.FormatInt(int64(generated), 10),
				FieldStatus,
				taskStatus,
			)
		}
		defer updateStatus()

		for atomic.LoadInt32(&stop) == 0 {
			<-ticker.C
			stopFlagStr, err := g.redisService.Get(KeyStopFlag)
			if err != nil {
				panic(err)
			}

			stopFlag, err := strconv.ParseInt(stopFlagStr, 10, 32)
			atomic.StoreInt32(&stop, int32(stopFlag))

			updateStatus()
			if atomic.LoadUint64(&generated) >= n {
				g.redisService.Set(KeyStopFlag, "1")
				g.redisService.Set(KeyRunning, "f")
				atomic.StoreInt32(&stop, 1)
				break
			}
		}

		log.Printf("status update goroutine is done.")
	}()
}
