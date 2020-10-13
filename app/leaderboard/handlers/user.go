package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	api2 "leaderboard/app/api"
	"leaderboard/app/leaderboard/services"
	"net/http"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(e *echo.Echo) {
	group := e.Group("/user")

	group.POST("/create", h.create)
	group.GET("/profile/:guid", h.get)
}

func (h *UserHandler) create(c echo.Context) (err error) {
	// TODO: Handle duplicate display name

	u := new(api2.UserProfile)
	if err = c.Bind(u); err != nil {
		return
	}

	if err = c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, api2.NewValidationErrorResponse(err.Error()))
	}

	guid, err := h.userService.Create(u)
	if err != nil {
		return err
	}

	ranked, err := h.userService.GetByIDWithRank(guid, "GLOBAL")
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, ranked)
}

func (h *UserHandler) get(c echo.Context) (err error) {
	guid := c.Param("guid")
	profile, err := h.userService.GetByIDWithRank(guid, "GLOBAL")
	if profile == nil || err != nil {
		return c.JSON(http.StatusNotFound, api2.UserNotFound{Message: fmt.Sprintf("User with ID(%s) is not found.", guid)})
	}

	return c.JSON(http.StatusOK, profile)
}
