package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	CUSTOM_VALIDATION_PASSWORD = "strong_password"
)

type FieldNotValid struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

var listFunc = map[string]func(fl validator.FieldLevel) bool{
	CUSTOM_VALIDATION_PASSWORD: ValidatePassword,
}

func JsonValidator(s interface{}, addFunction ...string) *[]FieldNotValid {
	validation := validator.New()
	for _, key := range addFunction {
		validation.RegisterValidation(key, listFunc[key])
	}
	if err := validation.Struct(s); err != nil {
		fieldNotFailed := make([]FieldNotValid, 0)
		for _, err := range err.(validator.ValidationErrors) {
			fieldNotFailed = append(fieldNotFailed, FieldNotValid{
				Field:  err.Field(),
				Reason: err.ActualTag(),
			})
		}
		return &fieldNotFailed
	}
	return nil
}

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Panjang minimal 8
	if len(password) < 8 {
		return false
	}
	// Uppercase
	uppercase, _ := regexp.MatchString(`[A-Z]`, password)
	// Lowercase
	lowercase, _ := regexp.MatchString(`[a-z]`, password)
	// Angka
	number, _ := regexp.MatchString(`[0-9]`, password)
	// Special char
	special, _ := regexp.MatchString(`[\W_]`, password)

	return uppercase && lowercase && number && special
}
