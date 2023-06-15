package utils

import (
	"axis/ecommerce-backend/pkg/entities"
)

func FormatApiError(errorMessage string, code int, errs entities.E) *entities.ApiError {
	return &entities.ApiError{Code: code, Errors: errs, Message: errorMessage}
}

func GeneralError(msg string, code int) *entities.ApiError {
	return FormatApiError(msg, code, entities.E{})
}
