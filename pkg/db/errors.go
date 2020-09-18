package db

import (
	"errors"

	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

func IsBadRequestErr(err error) bool {

	sqliteErr, ok := err.(sqlite3.Error)
	if ok && sqliteErr.Code == 19 {
		return true
	}

	possibleErrors := []error{gorm.ErrInvalidData, gorm.ErrInvalidField, gorm.ErrModelValueRequired, gorm.ErrPrimaryKeyRequired}
	for _, e := range possibleErrors {
		if errors.Is(err, e) {
			return true
		}
	}

	return false
}
