package repository

import (
	"net/http"

	"github.com/i-iterbium/cyber-engine/models"
)

// UserConfirmation описывает методы для работы с подтверждением пользователей
type UserConfirmation interface {
	Confirm(*http.Request, *models.UserConfirmation) (int, interface{}, error)
	ResendCode(*http.Request, *models.UserConfirmation) (int, interface{}, error)
}
