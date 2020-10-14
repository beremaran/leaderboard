package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"leaderboard/app/api"
	"leaderboard/app/leaderboard/services"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type ActuatorHandler struct {
	redisService api.RedisService
	userService  *services.UserService
}

func NewActuatorHandler(redisService api.RedisService, userService *services.UserService) *ActuatorHandler {
	return &ActuatorHandler{redisService: redisService, userService: userService}
}

func (a *ActuatorHandler) Register(echo *echo.Echo) {
	group := echo.Group("/_actuator")

	group.DELETE("/flush-all", a.FlushAll)
	group.GET("/bulk-generate", a.GenerateBulk)
}

// FlushAll godoc
// @Summary Flush Redis Cache
// @Description Remove all data
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 500
// @Tags actuator
// @Router /_actuator/flush-all [delete]
func (a *ActuatorHandler) FlushAll(c echo.Context) error {
	a.redisService.FlushAll()
	return c.NoContent(http.StatusOK)
}

// GenerateBulk godoc
// @Summary Generate users
// @Description Generate users
// @Produce  json
// @Success 200
// @Failure 500
// @Tags actuator
// @Param n query int true "how many users to generate" minimum(1)
// @Router /_actuator/bulk-generate [get]
func (a *ActuatorHandler) GenerateBulk(c echo.Context) error {
	n, err := strconv.ParseUint(c.QueryParam("n"), 10, 64)
	if err != nil {
		return err
	}

	countries := []string{"TR", "US", "GB", "CN", "JP", "AU", "NZ"}

	var wg sync.WaitGroup

	workload := n / uint64(runtime.NumCPU())
	c.Logger().Printf("gonna generate %d per goroutine", workload)
	ticker := time.NewTicker(time.Second)
	for cpu := 0; cpu < runtime.NumCPU()*2; cpu++ {
		go func() {
			wg.Add(1)
			defer wg.Done()

			var i uint64 = 0
			for ; i < workload; i++ {
				select {
				case <-ticker.C:
					c.Logger().Printf("generated %d", i+1)
					i--
				default:
					id := uuid.New().String()
					_, err = a.userService.Create(&api.UserProfile{
						UserId:      id,
						DisplayName: fmt.Sprintf("user_%d_%s", i, id),
						Points:      rand.Float64() * 100_000,
						Rank:        0,
						Country:     countries[rand.Intn(len(countries))],
					})

					if err != nil {
						panic(err)
					}
				}
			}
		}()
	}

	wg.Wait()
	ticker.Stop()
	return c.NoContent(http.StatusOK)
}
