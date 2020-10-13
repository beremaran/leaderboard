package handlers

import (
	"github.com/labstack/echo/v4"
	"leaderboard/app/api"
	"net/http"
)

type ActuatorHandler struct {
	redisService api.RedisService
}

func NewActuatorHandler(redisService api.RedisService) *ActuatorHandler {
	return &ActuatorHandler{redisService: redisService}
}

func (a *ActuatorHandler) Register(echo *echo.Echo) {
	group := echo.Group("/_actuator")

	group.DELETE("/flush-all", a.flushAll)
}

func (a *ActuatorHandler) flushAll(c echo.Context) error {
	a.redisService.FlushAll()
	return c.NoContent(http.StatusOK)
}
