package api

import (
	"regexp"
	"strings"
)

type LeaderboardRow struct {
	Rank        int64  `json:"rank"`
	Points      int64  `json:"points"`
	DisplayName string `json:"display_name"`
	Country     string `json:"country"`
}

type ScoreSubmission struct {
	Score     float64 `json:"score" validate:"required"`
	UserId    string  `json:"user_id" validate:"required"`
	Timestamp int64   `json:"timestamp" validate:"required"`
}

type UserProfile struct {
	UserId      string  `json:"user_id"`
	DisplayName string  `json:"display_name" validate:"required"`
	Points      float64 `json:"points"`
	Rank        int64   `json:"rank"`
	Country     string  `json:"country" validate:"required"`
}

type LeaderboardQuery struct {
	Country  string `json:"country" query:"country"`
	Page     int64  `json:"page" query:"page"`
	PageSize int64  `json:"page_size" query:"page_size"`
}

type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

func NewValidationErrorResponse(rawMessage string) *ValidationErrorResponse {
	rows := strings.Split(rawMessage, "\n")
	r, err := regexp.Compile("Key:\\s+'(?P<Key>.+)'\\s+Error:(?P<Message>.+)")
	if err != nil {
		panic(err)
	}

	var errors []ValidationError
	for _, row := range rows {
		matches := r.FindStringSubmatch(row)
		errors = append(errors, ValidationError{
			Path:    matches[1],
			Message: matches[2],
		})
	}

	return &ValidationErrorResponse{Errors: errors}
}

type UserNotFound struct {
	Message string `json:"message"`
}

type GenerateUserTaskStatus struct {
	Status         string `json:"status"`
	Concurrency    uint64 `json:"concurrency"`
	StartedAt      string `json:"started_at"`
	RemainingUsers uint64 `json:"remaining_users"`
}

type GenerateUserTaskConfiguration struct {
	NumberOfUsers uint64 `json:"nUsers" validate:"required"`
	Concurrency   uint64 `json:"concurrency" validate:"required"`
}
