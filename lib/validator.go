package lib

import (
	"github.com/go-playground/validator/v10"
)

func ValidateStruct(s interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Status = 403
			element.Title = "Validation failed!"
			element.Description = err.StructField() + " is not valid!"
			errors = append(errors, &element)
		}
	}
	return errors
}
