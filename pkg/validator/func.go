package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func DateValid(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		now := time.Now()
		if now.After(date) || now.Equal(date) {
			return true
		}
	}
	return false
}
