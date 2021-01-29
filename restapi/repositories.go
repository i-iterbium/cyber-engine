package restapi

import (
	"database/sql"
	"errors"

	"github.com/i-iterbium/cyber-engine/internal/pkg/repository"
	"github.com/i-iterbium/cyber-engine/internal/repositories/sessions"
	"github.com/i-iterbium/cyber-engine/internal/repositories/userconfirmation"
	"github.com/i-iterbium/cyber-engine/internal/repositories/users"
)

// Repositories описывает список репозиториев
type Repositories struct {
	Users            repository.Users
	Sessions         repository.Sessions
	UserConfirmation repository.UserConfirmation
}

func configureRepositories(db *sql.DB, err error) (Repositories, error) {
	if err != nil {
		return Repositories{}, errors.New("Ошибка получения конфигурации репозиториев")
	}

	return Repositories{
		Users:            users.New(db),
		Sessions:         sessions.New(db),
		UserConfirmation: userconfirmation.New(db),
	}, nil
}
