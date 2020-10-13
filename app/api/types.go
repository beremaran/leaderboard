package api

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
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

func SerializeUserProfile(profile *UserProfile) (string, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(profile); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func DeserializeUserProfile(b64string string) (*UserProfile, error) {
	raw, err := base64.StdEncoding.DecodeString(b64string)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buf)

	profile := new(UserProfile)
	if err := dec.Decode(profile); err != nil {
		return nil, err
	}

	return profile, nil
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
