package services

import (
	"strings"

	"vocabulary/models"

	"github.com/go-playground/validator/v10"
)

var registerValidator = func() *validator.Validate {
	v := validator.New()
	v.SetTagName("binding")
	return v
}()

type RegistrationValidationError struct {
	Message string
}

func (e *RegistrationValidationError) Error() string { return e.Message }

// ValidateRegisterRequest normalizes identity fields without changing passwords.
func ValidateRegisterRequest(request *models.RegisterRequest) error {
	request.DisplayName = strings.TrimSpace(request.DisplayName)
	request.Email = normalizeEmail(request.Email)
	if err := registerValidator.Struct(request); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}
		field := validationErrors[0]
		message := "invalid registration data"
		switch field.StructField() {
		case "DisplayName":
			message = "display_name is required and must be at most 100 characters"
		case "Email":
			message = "email must be a valid email address of at most 255 characters"
		case "Password":
			message = "password must be at least 8 characters"
		case "ConfirmPassword":
			message = "confirm_password is required"
			if field.Tag() == "eqfield" {
				message = "confirm_password must match password"
			}
		}
		return &RegistrationValidationError{Message: message}
	}
	// bcrypt limits passwords by bytes, whereas min=8 above counts characters.
	if len(request.Password) > 72 {
		return &RegistrationValidationError{Message: "password must be at most 72 bytes"}
	}
	return nil
}
