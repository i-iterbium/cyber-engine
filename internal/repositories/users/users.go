package users

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/i-iterbium/cyber-engine/internal/pkg/confirmationcode"
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

// Fetch возвращает данные пользователя по его ID
func (repo *Repository) Fetch(r *http.Request, ID int64) (int, *models.User, error) {
	status, s, err := session.Get(repo.db, r)
	if err != nil {
		return status, nil, err
	}
	if s == nil || s.UserID == nil {
		return http.StatusBadRequest, nil, errors.New("Пользователь не авторизован")
	}
	if *s.UserID != ID {
		return http.StatusBadRequest, nil, errors.New("Недопустимый запрос данных другого пользователя")
	}

	rows, err := repo.db.Query(`
		select
			id,
			name,
			sname,
			pname,
			phone,
			email,
			birthday
		from v_user 
		where id = $1`, ID)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer rows.Close()

	var u models.User
	if rows.Next() {
		if err = sqlstruct.Scan(&u, rows); err != nil {
			_, err = errorutil.HandleDBError(err)
			return http.StatusInternalServerError, nil, err
		}
	}

	if u.ID == nil {
		return http.StatusInternalServerError, nil, errors.New("Пользователь с данным кодом не существует")
	}

	return http.StatusOK, &u, nil
}

// Create создает пользователя
func (repo *Repository) Create(r *http.Request, user *models.NewUser) (int, *models.User, error) {

	status, s, err := session.Get(repo.db, r)
	if err != nil {
		return status, nil, err
	}

	if s != nil {
		return http.StatusForbidden, nil, errors.New("Пользователь уже авторизован")
	}

	if !validation.Password(user.Password) {
		return http.StatusBadRequest, nil, errors.New("Некорректный формат пароля")
	}

	if *user.RegisterType == "phone" {
		return regByPhone(repo, user)
	}
	if *user.RegisterType == "email" {
		return regByEmail(repo, user)
	}

	return http.StatusBadRequest, nil, errors.New("Некорректный тип регистрации")

}

func regByPhone(repo *Repository, user *models.NewUser) (int, *models.User, error) {

	if !validation.Phone(user.Phone, false) {
		return http.StatusBadRequest, nil, errors.New("Некорректный формат номера телефона")
	}

	tx, txErr := repo.db.Begin()
	if txErr != nil {
		return 0, nil, txErr
	}
	defer func() { transaction.CompleteTx(tx, txErr) }()

	u, txErr := createUser(tx, user)
	if txErr != nil {
		return http.StatusInternalServerError, nil, txErr
	}

	status, timeout, txErr := confirmationcode.Send(tx, u.ID, *user.RegisterType, "userConfirm")
	if txErr != nil {
		return status, nil, txErr
	}

	u.Timeout = *timeout

	return http.StatusOK, u, nil

}

func regByEmail(repo *Repository, user *models.NewUser) (int, *models.User, error) {

	if !validation.Email(string(*user.Email), false) {
		return http.StatusBadRequest, nil, errors.New("Некорректный формат адреса электронной почты")
	}

	tx, txErr := repo.db.Begin()
	if txErr != nil {
		return 0, nil, txErr
	}
	defer func() { transaction.CompleteTx(tx, txErr) }()

	u, txErr := createUser(tx, user)
	if txErr != nil {
		return http.StatusInternalServerError, nil, txErr
	}

	status, timeout, txErr := confirmationcode.Send(tx, u.ID, *user.RegisterType, "userConfirm")
	if txErr != nil {
		return status, nil, txErr
	}

	u.Timeout = *timeout

	return http.StatusOK, u, nil
}

func createUser(tx *sql.Tx, user *models.NewUser) (*models.User, error) {

	rows, txErr := tx.Query(`select * from fn_user_ins($1, $2, $3, $4, $5, $6, $7)`,
		&user.Phone,
		&user.Password,
		&user.Name,
		&user.Sname,
		&user.Pname,
		&user.Email,
		&user.Birthday,
	)
	if txErr != nil {
		_, txErr = errorutil.HandleDBError(txErr)
		return nil, txErr
	}
	defer rows.Close()

	var u models.User
	if rows.Next() {
		if txErr = sqlstruct.Scan(&u, rows); txErr != nil {
			_, txErr = errorutil.HandleDBError(txErr)
			return nil, txErr
		}
	}

	return &u, nil
}

// Patch обновляет данные пользователя
func (repo *Repository) Patch(r *http.Request, ID int64, user *models.EditUser) (int, *models.User, error) {

	status, s, err := session.Get(repo.db, r)
	if err != nil {
		return status, nil, err
	}

	if *s.UserID != ID {
		return http.StatusBadRequest, nil, errors.New("Недопустимый запрос данных другого пользователя")
	}

	tx, txErr := repo.db.Begin()
	if txErr != nil {
		return http.StatusInternalServerError, nil, txErr
	}
	defer func() { transaction.CompleteTx(tx, txErr) }()

	var rows *sql.Rows
	rows, txErr = tx.Query(`select * from fn_user_upd($1, $2, $3, $4, $5, $6)`,
		&ID,
		&user.Name,
		&user.Sname,
		&user.Pname,
		&user.Birthday,
		&user.Email,
	)
	if txErr != nil {
		_, txErr = errorutil.HandleDBError(txErr)
		return http.StatusInternalServerError, nil, txErr
	}
	defer rows.Close()

	var u *models.User
	if rows.Next() {
		if txErr = sqlstruct.Scan(&u, rows); txErr != nil {
			_, txErr = errorutil.HandleDBError(txErr)
			return http.StatusInternalServerError, nil, txErr
		}
	}

	return http.StatusOK, u, nil
}

// // Delete удаляет пользователя
// func (repo *Repository) Delete(ID int64) (int, error) {
// 	tx, txErr := repo.db.Begin()
// 	if txErr != nil {
// 		return http.StatusInternalServerError, txErr
// 	}
// 	defer func() { transaction.CompleteTx(tx, txErr) }()

// 	_, txErr = tx.Exec(`select * from fn_user_del($1)`, ID)
// 	if txErr != nil {
// 		_, txErr = errorutil.HandleDBError(txErr)
// 		return http.StatusInternalServerError, txErr
// 	}

// 	return http.StatusOK, nil
// }
