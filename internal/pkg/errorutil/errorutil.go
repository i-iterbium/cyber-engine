package errorutil

import (
	"errors"
	"strings"
)

// HandleDBError отделяет ошибки базы данных
func HandleDBError(err error) (bool, error) {
	return trimPrefix(err, "pq: ")
}

func trimPrefix(err error, prefix string) (bool, error) {
	errorMsg := err.Error()
	if !strings.HasPrefix(errorMsg, prefix) {
		return false, err
	}
	errorMsg = strings.TrimSpace(errorMsg[len(prefix):])
	return true, errors.New(errorMsg)
}
