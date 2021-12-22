package lib

import (
	"callme/models"

	"github.com/go-playground/validator/v10"
)

func ValidateRegister(user models.User) ErrorResponse {
	var errors ErrorResponse
	validate := validator.New()
	err := validate.StructPartial(user, "Password", "Email", "Username")
	if err != nil {
		errors.Status = 400
		for i, error := range err.(validator.ValidationErrors) {
			errors.Description += error.Field() + " is not valid"
			if i != len(err.(validator.ValidationErrors))-1 {
				errors.Description += " && "
			}
		}
	}
	return errors
}
