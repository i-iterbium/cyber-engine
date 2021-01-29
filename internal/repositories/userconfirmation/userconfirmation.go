package userconfirmation

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/i-iterbium/cyber-engine/internal/pkg/confirmationcode"
	"github.com/i-iterbium/cyber-engine/internal/pkg/errorutil"
	"github.com/i-iterbium/cyber-engine/internal/pkg/session"
	"github.com/i-iterbium/cyber-engine/internal/pkg/transaction"
	"github.com/i-iterbium/cyber-engine/models"
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

// Confirm производит подтверждение пользователя по коду
func (repo *Repository) Confirm(r *http.Request, uc *models.UserConfirmation) (int, interface{}, error) {

	if uc.Code == nil {
		return http.StatusBadRequest, nil, errors.New("Не указан код подтверждения")
	}

	status, s, err := session.Get(repo.db, r)
	if err != nil {
		return status, nil, err
	}

	if s != nil {
		return http.StatusForbidden, nil, errors.New("Пользователь уже авторизован")
	}

	var userID int64
	switch *uc.ConfirmationType {
	case "phone":
		err := repo.db.QueryRow(`
			select
				id
			from v_user 
			where phone = $1`, *uc.Phone).Scan(&userID)
		if err != nil {
			_, err = errorutil.HandleDBError(err)
			return http.StatusInternalServerError, nil, err
		}
	case "email":
		err := repo.db.QueryRow(`
			select
				id
			from v_user 
			where email = $1`, *uc.Email).Scan(&userID)
		if err != nil {
			_, err = errorutil.HandleDBError(err)
			return http.StatusInternalServerError, nil, err
		}
	default:
		return http.StatusBadRequest, nil, errors.New("Некорректный тип подтверждения")
	}

	if &userID == nil {
		return http.StatusInternalServerError, nil, errors.New("Пользователь не существует")
	}

	return confirmUser(repo, &userID, uc.Code)
}

func confirmUser(repo *Repository, userID *int64, code *string) (int, interface{}, error) {

	tx, txErr := repo.db.Begin()
	if txErr != nil {
		return http.StatusInternalServerError, nil, txErr
	}
	defer func() { transaction.CompleteTx(tx, txErr) }()

	_, txErr = tx.Exec(`select * from fn_confirm_user($1, $2)`, *userID, *code)
	if txErr != nil {
		_, txErr = errorutil.HandleDBError(txErr)
		return http.StatusInternalServerError, nil, txErr
	}

	return http.StatusOK, nil, nil
}

// ResendCode производит отправку кода подтверждения пользователю
func (repo *Repository) ResendCode(r *http.Request, rc *models.UserConfirmation) (int, interface{}, error) {

	status, s, err := session.Get(repo.db, r)
	if err != nil {
		return status, nil, err
	}

	if s != nil {
		return http.StatusForbidden, nil, errors.New("Пользователь уже авторизован")
	}

	var userID int64
	switch *rc.ConfirmationType {
	case "phone":
		err := repo.db.QueryRow(`
			select
				id
			from v_user 
			where phone = $1`, *rc.Phone).Scan(&userID)
		if err != nil {
			_, err = errorutil.HandleDBError(err)
			return http.StatusInternalServerError, nil, err
		}
	case "email":
		err := repo.db.QueryRow(`
			select
				id
			from v_user 
			where email = $1`, *rc.Email).Scan(&userID)
		if err != nil {
			_, err = errorutil.HandleDBError(err)
			return http.StatusInternalServerError, nil, err
		}
	default:
		return http.StatusBadRequest, nil, errors.New("Некорректный тип подтверждения")
	}

	if &userID == nil {
		return http.StatusInternalServerError, nil, errors.New("Пользователь не существует")
	}

	tx, txErr := repo.db.Begin()
	if txErr != nil {
		return 0, nil, txErr
	}
	defer func() { transaction.CompleteTx(tx, txErr) }()

	status, _, txErr = confirmationcode.Send(tx, &userID, *rc.ConfirmationType, "userConfirm")
	if txErr != nil {
		return status, nil, txErr
	}

	return http.StatusOK, nil, nil
}
