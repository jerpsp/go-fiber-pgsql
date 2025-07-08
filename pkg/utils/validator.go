package utils

import (
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
