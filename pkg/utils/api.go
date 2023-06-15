package utils

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/pkg/entities"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

func GetApiResponse(data interface{}, msg string, links entities.E, meta entities.E) *entities.ApiResponse {
	return &entities.ApiResponse{
		Message: msg,
		Data:    data,
		Links:   links,
		Meta:    meta,
	}
}

func FormatValidationError(errs validator.ValidationErrors) entities.E {
	fieldErrors := entities.E{}
	for _, fieldEr := range errs {
		name := fieldEr.Field()
		name = strcase.ToSnake(name)
		errorMessage := fieldEr.Translate(configs.Trans)
		if len(name) != 0 && len(errorMessage) != 0 {
			fieldErrors[name] = errorMessage
		}
	}
	return fieldErrors
}
