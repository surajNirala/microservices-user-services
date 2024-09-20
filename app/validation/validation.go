package validation

import (
	"fmt"
	"strings"

	custom_validator "github.com/go-playground/validator"
	"github.com/go-playground/validator/v10"
)

// translateValidationErrors translates validator.ValidationErrors to a map of field errors
func TranslateValidationErrors(err error) map[string]string { //[]
	customErrors := make(map[string]string) //[]
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			field := strings.ToLower(err.Field())
			var message string
			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", field)
			case "email":
				message = fmt.Sprintf("%s is not a valid email address", field)
			default:
				message = fmt.Sprintf("%s is not valid", field)
			}
			customErrors[field] = message
		}
	}
	return customErrors
}

// Custom login Validation manage
func LoginValidationErrors(err error) map[string]string {
	customErrors := make(map[string]string)
	if validationErrors, ok := err.(custom_validator.ValidationErrors); ok {
		fmt.Println("ok : ", ok)
		fmt.Println("validationErrors : ", validationErrors)
		for _, err := range validationErrors {
			field := strings.ToLower(err.Field())
			var message string
			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", field)
			case "email":
				message = fmt.Sprintf("%s is not a valid email address", field)
			default:
				message = fmt.Sprintf("%s is not valid", field)
			}
			customErrors[field] = message
		}
	}
	fmt.Println("customErrors :", customErrors)
	return customErrors
}
