package tasks

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"leaderboard/app/api"
	"leaderboard/app/leaderboard/services"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const KeyTask = "TASK_GU"
const FieldStatus = "STATUS"
const FieldCompleted = "COMPLETED"
const FieldConcurrency = "CONCURRENCY"
const FieldStartedAt = "STARTED_AT"
const FieldUserCount = "USER_COUNT"

type GenerateUsersSingletonTask struct {
	userService  *services.UserService
	redisService api.RedisService
	stateMux     sync.Mutex
}

func NewGenerateUsersSingletonTask(userService *services.UserService, redisService api.RedisService) *GenerateUsersSingletonTask {
	return &GenerateUsersSingletonTask{userService: userService, redisService: redisService}
}

func (g *GenerateUsersSingletonTask) Start(nUsers uint64, maxConcurrency int64) (*api.GenerateUserTaskStatus, error) {
	if g.getStatusStr(true) == "RUNNING" {
		return g.Status()
	}

	initResult := make(chan bool, 1)
	go g.generate(nUsers, maxConcurrency, initResult)

	result := <-initResult
	if !result {
		return nil, fmt.Errorf("task initialization failed")
	}

	return g.Status()
}

func (g *GenerateUsersSingletonTask) Stop() error {
	g.updateStatus(&api.GenerateUserTaskStatus{
		Status: "CANCELLED",
	}, true)

	return nil
}

func (g *GenerateUsersSingletonTask) Status() (*api.GenerateUserTaskStatus, error) {
	return g.status(true)
}

func (g *GenerateUsersSingletonTask) status(withMux bool) (*api.GenerateUserTaskStatus, error) {
	if withMux {
		g.stateMux.Lock()
		defer g.stateMux.Unlock()
	}

	exists, err := g.redisService.Exists(KeyTask)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &api.GenerateUserTaskStatus{
			Status: "IDLE",
		}, nil
	}

	statusMap, err := g.redisService.HGetAll(KeyTask).Result()
	if err != nil {
		return nil, err
	}

	var completedRatio float64 = 0
	if statusMap[FieldCompleted] != "" {
		completedRatio, err = strconv.ParseFloat(statusMap[FieldCompleted], 64)
	}
	if err != nil {
		return nil, err
	}

	var concurrency uint64 = 0
	if statusMap[FieldConcurrency] != "" {
		concurrency, err = strconv.ParseUint(statusMap[FieldConcurrency], 10, 64)
	}
	if err != nil {
		return nil, err
	}

	var userCount = ^uint64(0)
	if statusMap[FieldUserCount] != "" {
		userCount, err = strconv.ParseUint(statusMap[FieldUserCount], 10, 64)
	}
	if err != nil {
		return nil, err
	}

	return &api.GenerateUserTaskStatus{
		Status:         statusMap[FieldStatus],
		Completed:      completedRatio,
		Concurrency:    concurrency,
		StartedAt:      statusMap[FieldStartedAt],
		RemainingUsers: userCount,
	}, nil
}

func (g *GenerateUsersSingletonTask) generate(n uint64, maxConcurrency int64, redisInit chan bool) {
	status, err := g.status(true)
	if err != nil {
		g.handleError(err, true)
		redisInit <- false
		return
	}
	status.RemainingUsers = n
	status.Status = "RUNNING"
	g.updateStatus(status, true)

	countries := []string{"TR", "US", "GB", "CN", "JP", "AU", "NZ"}

	redisInit <- true

	var cpu int64
	for cpu = 0; cpu < maxConcurrency; cpu++ {
		go func() {
			statusStr := g.getStatusStr(true)
			if statusStr == "CANCELLED" || statusStr == "DONE" || statusStr == "ERROR" {
				return
			}

			for g.getUserCount(true) > 0 {
				_, err := g.userService.Create(&api.UserProfile{
					DisplayName: fmt.Sprintf("user_%d", time.Now().UnixNano()),
					Points:      rand.Float64() * 100_000,
					Country:     countries[rand.Intn(len(countries))],
				})

				if err != nil {
					g.handleError(err, true)
					return
				}

				g.decrementToGenerate()
			}
		}()
	}
}

func (g *GenerateUsersSingletonTask) handleError(err error, withMux bool) {
	if withMux {
		g.stateMux.Lock()
		defer g.stateMux.Unlock()
	}

	g.updateStatus(&api.GenerateUserTaskStatus{
		Status:    "ERROR",
		StartedAt: time.Now().String(),
	}, !withMux)

	log.Error(err)
}

func (g *GenerateUsersSingletonTask) updateStatus(status *api.GenerateUserTaskStatus, withMux bool) {
	if withMux {
		g.stateMux.Lock()
		defer g.stateMux.Unlock()
	}

	g.redisService.HSet(
		KeyTask,
		FieldCompleted,
		strconv.FormatFloat(status.Completed, 'f', 2, 64),
		FieldConcurrency,
		strconv.FormatUint(status.Concurrency, 10),
		FieldUserCount,
		strconv.FormatUint(status.RemainingUsers, 10),
		FieldStartedAt,
		status.StartedAt,
		FieldStatus,
		status.Status,
	)
}

func (g *GenerateUsersSingletonTask) getUserCount(withMux bool) uint64 {
	if withMux {
		g.stateMux.Lock()
		defer g.stateMux.Unlock()
	}

	status, err := g.status(!withMux)
	if err != nil {
		g.handleError(err, !withMux)
		return 0
	}

	return status.RemainingUsers
}

func (g *GenerateUsersSingletonTask) getStatusStr(withMux bool) string {
	if withMux {
		g.stateMux.Lock()
		defer g.stateMux.Unlock()
	}

	status, err := g.status(!withMux)
	if err != nil {
		g.handleError(err, !withMux)
		return ""
	}

	return status.Status
}

func (g *GenerateUsersSingletonTask) decrementToGenerate() {
	g.stateMux.Lock()

	status, err := g.status(false)
	if err != nil {
		g.stateMux.Unlock()
		g.handleError(err, true)
		return
	}
	defer g.stateMux.Unlock()

	status.RemainingUsers--
	if status.RemainingUsers == 0 {
		status.Status = "DONE"
	}

	g.updateStatus(status, false)
}
