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

	group.GET("", l.GetLeaderboard)
	group.GET("/:country_iso_code", l.GetLeaderboard)
}

// GetLeaderboard godoc
// @Summary Get leaderboard
// @Description Get leaderboard
// @Produce  json
// @Success 200 {array} api.LeaderboardRow
// @Failure 500
// @Tags leaderboard
// @Param page query int false "page number" minimum(1)
// @Param page_size query int false "number of records in a page" minimum(1)
// @Param page_size query int false "number of records in a page" minimum(1)
// @Router /leaderboard [get]
func (l *LeaderboardHandler) GetLeaderboard(c echo.Context) error {
	return l.handleLeaderboardRequest(c)
}

// GetLeaderboardByCountryCode godoc
// @Summary Get leaderboard
// @Description Get leaderboard
// @Produce  json
// @Success 200 {array} api.LeaderboardRow
// @Failure 500
// @Tags leaderboard
// @Param page query int false "page number" minimum(1)
// @Param page_size query int false "number of records in a page" minimum(1)
// @Param page_size query int false "number of records in a page" minimum(1)
// @Param country_iso_code path string false "ISO standard country code"
// @Router /leaderboard/{country_iso_code} [get]
func (l *LeaderboardHandler) GetLeaderboardByCountryCode(c echo.Context) error {
	return l.handleLeaderboardRequest(c)
}

func (l *LeaderboardHandler) handleLeaderboardRequest(c echo.Context) (err error) {
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
