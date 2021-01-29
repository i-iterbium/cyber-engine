package sessions

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/i-iterbium/cyber-engine/internal/pkg/errorutil"
	"github.com/i-iterbium/cyber-engine/internal/pkg/session"
	"github.com/i-iterbium/cyber-engine/internal/pkg/transaction"
	"github.com/i-iterbium/cyber-engine/internal/pkg/validation"
	"github.com/i-iterbium/cyber-engine/models"
	"github.com/kisielk/sqlstruct"
)

// Repository описывает подключения к БД
type Repository struct {
	db *sql.DB
}

// New создает новый экземпляр репозитория БД
func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create создает пользовательскую сессию
func (repo *Repository) Create(req *http.Request, auth *models.AuthData) (int, *models.Session, error) {
	if !validation.Password(auth.Password) {
		return http.StatusBadRequest, nil, errors.New("Некорректный формат пароля")
	}

	if auth.Phone != nil {
		return authByPhone(repo, auth.Phone, auth.Password)
	}
	if auth.Email != nil {
		return authByEmail(repo, auth.Email, auth.Password)
	}

	return http.StatusBadRequest, nil, errors.New("Переданы неполные данные для авторизации")
}

func authByPhone(repo *Repository, phone *int64, password *string) (int, *models.Session, error) {
	if !validation.Phone(phone, false) {
		return http.StatusBadRequest, nil, errors.New("Некорректный формат номера телефона")
	}

	tx, txErr := repo.db.Begin()
	if txErr != nil {
		return http.StatusInternalServerError, nil, txErr
	}
	defer func() { transaction.CompleteTx(tx, txErr) }()

	var rows *sql.Rows
	rows, txErr = tx.Query(`select * from fn_session_by_phone_ins($1, $2)`,
		&phone,
		&password,
	)
	if txErr != nil {
		_, txErr = errorutil.HandleDBError(txErr)
		return http.StatusInternalServerError, nil, txErr
	}
	defer rows.Close()

	var s models.Session
	if rows.Next() {
		if txErr = sqlstruct.Scan(&s, rows); txErr != nil {
			_, txErr = errorutil.HandleDBError(txErr)
			return http.StatusInternalServerError, nil, txErr
		}
	}

	return http.StatusOK, &s, nil
}

func authByEmail(repo *Repository, email *strfmt.Email, password *string) (int, *models.Session, error) {
	if !validation.Email(string(*email), false) {
		return http.StatusBadRequest, nil, errors.New("Некорректный формат адреса электронной почты")
	}

	tx, txErr := repo.db.Begin()
	if txErr != nil {
		return http.StatusInternalServerError, nil, txErr
	}
	defer func() { transaction.CompleteTx(tx, txErr) }()

	var rows *sql.Rows
	rows, txErr = tx.Query(`select * from fn_session_by_email_ins($1, $2)`,
		&email,
		&password,
	)
	if txErr != nil {
		_, txErr = errorutil.HandleDBError(txErr)
		return http.StatusInternalServerError, nil, txErr
	}
	defer rows.Close()

	var s models.Session
	if rows.Next() {
		if txErr = sqlstruct.Scan(&s, rows); txErr != nil {
			_, txErr = errorutil.HandleDBError(txErr)
			return http.StatusInternalServerError, nil, txErr
		}
	}

	return http.StatusOK, &s, nil
}

// Update обновляет пользовательскую сессию
func (repo *Repository) Update(req *http.Request) (int, *models.Session, error) {
	status, s, err := session.Get(repo.db, req)
	if err != nil {
		return status, nil, err
	}
	if s == nil || s.UserID == nil {
		return http.StatusBadRequest, nil, errors.New("Пользователь не авторизован")
	}

	tx, txErr := repo.db.Begin()
	if txErr != nil {
		return http.StatusInternalServerError, nil, txErr
	}
	defer func() { transaction.CompleteTx(tx, txErr) }()

	var rows *sql.Rows
	rows, txErr = tx.Query(`select * from fn_session_upd($1, $2)`, s.Session, s.CSRFToken)
	if txErr != nil {
		_, txErr = errorutil.HandleDBError(txErr)
		return http.StatusInternalServerError, nil, txErr
	}
	defer rows.Close()

	var session models.Session
	if rows.Next() {
		if txErr = sqlstruct.Scan(&session, rows); txErr != nil {
			_, txErr = errorutil.HandleDBError(txErr)
			return http.StatusInternalServerError, nil, txErr
		}
	}

	return http.StatusOK, &session, nil
}

// Delete удаляет пользовательскую сессию
func (repo *Repository) Delete(req *http.Request) (int, error) {
	status, s, err := session.Get(repo.db, req)
	if err != nil {
		return status, err
	}
	if s == nil || s.UserID == nil {
		return http.StatusBadRequest, errors.New("Пользователь не авторизован")
	}

	tx, txErr := repo.db.Begin()
	if txErr != nil {
		return http.StatusInternalServerError, txErr
	}
	defer func() { transaction.CompleteTx(tx, txErr) }()

	_, txErr = tx.Exec(`select * from fn_session_del($1, $2)`, s.Session, s.CSRFToken)
	if txErr != nil {
		_, txErr = errorutil.HandleDBError(txErr)
		return http.StatusInternalServerError, txErr
	}

	return http.StatusOK, nil
}
