package error

import (
	"encoding/json"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type PostgresError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
	/*
		Other possible fields:
			Severity,
			Hint,
			Position,
			InternalPosition,
			InternalQuery,
			Where,
			Schema,
			TableName,
			ColumnName,
			DataTypeName,
			ConstraintName,
			File,
			Line,
			Routine
	*/
}

// Convert gorm error to app error, this support only postgres error
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

	var pgErr PostgresError
	buf, err := json.Marshal(gormError)
	if err == nil {
		if err := json.Unmarshal(buf, &pgErr); err == nil {
			switch pgErr.Code {
			case "23505":
				return New(http.StatusConflict, gormError, "duplicated key")
			case "23503":
				return New(http.StatusNotFound, gormError, "record not found")
			default:
				return New(http.StatusInternalServerError, gormError, "internal server error")
			}
		}
	}

	return New(http.StatusInternalServerError, gormError, "internal server error")
}
