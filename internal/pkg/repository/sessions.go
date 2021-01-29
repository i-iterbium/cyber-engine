package repository

import (
	"net/http"

	"github.com/i-iterbium/cyber-engine/models"
)

// Sessions описывает методы для работы с сессиями пользователей
type Sessions interface {
	Create(*http.Request, *models.AuthData) (int, *models.Session, error)
	Update(*http.Request) (int, *models.Session, error)
	Delete(*http.Request) (int, error)
}
