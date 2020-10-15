package handlers

import (
	"github.com/labstack/echo/v4"
	"leaderboard/app/api"
	"leaderboard/app/leaderboard/services"
	"leaderboard/app/leaderboard/tasks"
	"net/http"
	"strconv"
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
	group.DELETE("/bulk-generate", a.StopGenerateBulk)
	group.GET("/user-count", a.GetUserCount)
}

// GetUserCount godoc
// @Summary Get total number of users
// @Description Get total number of users
// @Produce  plain
// @Success 200
// @Failure 500
// @Tags actuator
// @Router /_actuator/user-count [get]
func (a *ActuatorHandler) GetUserCount(c echo.Context) (err error) {
	size, err := a.redisService.GetSortedSetSize("GLOBAL")

	return c.String(http.StatusOK, strconv.FormatInt(size, 10))
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
// @Param n body int true "how many users to generate" minimum(1)
// @Param concurrency body int true "generate with concurrency" minimum(1)
// @Router /_actuator/bulk-generate [post]
func (a *ActuatorHandler) GenerateBulk(c echo.Context) error {
	n, err := strconv.ParseUint(c.QueryParam("n"), 10, 64)
	if err != nil {
		return err
	}
	concurrency, err := strconv.ParseInt(c.QueryParam("concurrency"), 10, 64)
	if err != nil {
		return err
	}

	task := a.getUserGenerateTask()
	taskStatus, err := task.Start(n, concurrency)
	return c.JSON(http.StatusOK, taskStatus)
}

// QueryBulkGeneration godoc
// @Summary Query user generation task status
// @Description Query user generation task status
// @Produce  json
// @Success 200
// @Failure 500
// @Tags actuator
// @Router /_actuator/bulk-generate [get]
func (a *ActuatorHandler) QueryBulkGeneration(c echo.Context) error {
	status, err := a.getUserGenerateTask().Status()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, status)
}

func (a *ActuatorHandler) getUserGenerateTask() *tasks.GenerateUsersSingletonTask {
	return tasks.NewGenerateUsersSingletonTask(a.userService, a.redisService)
}

// StopGenerateBulk godoc
// @Summary Stop user generation
// @Description Stop user generation
// @Success 200
// @Failure 500
// @Tags actuator
// @Router /_actuator/bulk-generate [delete]
func (a *ActuatorHandler) StopGenerateBulk(c echo.Context) error {
	err := a.getUserGenerateTask().Stop()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, err)
}
