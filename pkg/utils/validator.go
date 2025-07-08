package utils

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validatorInstance *validator.Validate
)

func GetValidator() *validator.Validate {
	sync.OnceFunc(func() {
		validatorInstance = validator.New()
	})()

	return validatorInstance
}

// Validate validates a struct based on the validator tags
func Validate(s interface{}) error {
	if err := GetValidator().Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				return fmt.Errorf("validation error: field '%s' failed on the '%s' tag", e.Field(), e.Tag())
			}
		}
		return err
	}
	return nil
}
