package repository

import (
	"net/http"

	"github.com/i-iterbium/cyber-engine/models"
)

// Users описывает методы для работы с пользователями
type Users interface {
	Fetch(r *http.Request, ID int64) (int, *models.User, error)
	Create(*http.Request, *models.NewUser) (int, *models.User, error)
	Patch(r *http.Request, ID int64, editUser *models.EditUser) (int, *models.User, error)
}
