package db

import (
	pg "github.com/lib/pq"
)

const (
	uniqueViolation pg.ErrorCode = "23505"
)

func IsBadRequestErr(err error) bool {
	if pgErr, isPGErr := err.(pg.Error); isPGErr {
		return pgErr.Code != uniqueViolation
	}
	return false
}
