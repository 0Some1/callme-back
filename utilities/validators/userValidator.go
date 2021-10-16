package validators

import (
	"callme/models"
	"github.com/go-playground/validator/v10"
)

func ValidateStruct(user models.User) []*models.ErrorResponse {
	var errors []*models.ErrorResponse
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element models.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
