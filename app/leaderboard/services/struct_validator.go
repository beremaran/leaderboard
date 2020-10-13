package services

import "github.com/go-playground/validator/v10"

type StructValidator struct {
	validator *validator.Validate
}

func NewStructValidator(validator *validator.Validate) *StructValidator {
	return &StructValidator{validator: validator}
}

func (cv *StructValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
