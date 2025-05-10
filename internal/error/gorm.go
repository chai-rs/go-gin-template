package error

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

func FromGorm(gormError error) error {
	if gormError == nil {
		return nil
	}

	if errors.Is(gormError, gorm.ErrRecordNotFound) {
		return New(http.StatusNotFound, gormError, "record not found")
	}

	if errors.Is(gormError, gorm.ErrInvalidData) {
		return New(http.StatusBadRequest, gormError, "invalid data")
	}

	if errors.Is(gormError, gorm.ErrDuplicatedKey) {
		return New(http.StatusConflict, gormError, "duplicated key")
	}

	return New(http.StatusInternalServerError, gormError, "internal server error")
}
