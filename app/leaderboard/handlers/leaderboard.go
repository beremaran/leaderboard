package handlers

import (
	"github.com/labstack/echo/v4"
	api2 "leaderboard/app/api"
	"leaderboard/app/leaderboard/services"
	"net/http"
	"strings"
)

type LeaderboardHandler struct {
	leaderboardService *services.LeaderboardService
}

func NewLeaderboardHandler(leaderboardService *services.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{leaderboardService: leaderboardService}
}

func (l *LeaderboardHandler) Register(echo *echo.Echo) {
	group := echo.Group("/leaderboard")

	group.GET("", l.getLeaderboard)
	group.GET("/:country_iso_code", l.getLeaderboard)
}

func (l *LeaderboardHandler) getLeaderboard(c echo.Context) (err error) {
	q := new(api2.LeaderboardQuery)
	if err = c.Bind(q); err != nil {
		return
	}

	if q.Page <= 0 {
		q.Page = 1
	}

	if q.PageSize <= 0 {
		q.PageSize = 10
	}

	if err = c.Validate(q); err != nil {
		return c.JSON(http.StatusBadRequest, api2.NewValidationErrorResponse(err.Error()))
	}

	countryParam := c.Param("country_iso_code")
	if len(countryParam) > 0 {
		q.Country = strings.ToUpper(countryParam)
	}

	if len(q.Country) == 0 {
		q.Country = "GLOBAL"
	}

	page, err := l.leaderboardService.GetPage(q.Country, q.Page, q.PageSize)
	if err != nil || page == nil {
		page = []*api2.LeaderboardRow{}
	}

	return c.JSON(http.StatusOK, page)
}
