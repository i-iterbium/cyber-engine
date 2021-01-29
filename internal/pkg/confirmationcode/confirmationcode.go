package confirmationcode

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/i-iterbium/cyber-engine/internal/pkg/errorutil"
)

type codeResponse struct {
	Code    *string `json:"code,omitempty"    sql:"code"`
	Timeout *int64  `json:"timeout,omitempty" sql:"timeout"`
}

// Send отправляет код подтверждения
func Send(tx *sql.Tx, userID *int64, sendingType string, codeType string) (int, *int64, error) {

	if userID == nil {
		return http.StatusBadRequest, nil, errors.New("Не передан идентификатор пользователя")
	}

	var resp codeResponse
	switch codeType {
	case "userConfirm":
		txErr := tx.QueryRow(`select * from fn_user_confirmation_code_get($1)`, *userID).Scan(&resp.Code, &resp.Timeout)
		if txErr != nil {
			_, txErr = errorutil.HandleDBError(txErr)
			return http.StatusInternalServerError, nil, txErr
		}
	default:
		return http.StatusInternalServerError, nil, errors.New("Некорректный тип кода")
	}

	if resp.Code == nil && resp.Timeout == nil {
		return http.StatusBadRequest, nil, errors.New("Пользователь заблокирован")
	}

	if resp.Code == nil && *resp.Timeout > 0 {
		return http.StatusOK, resp.Timeout, errors.New("Таймаут ещё не вышел")
	}

	msgText := "Код подтверждения: " + *resp.Code
	switch sendingType {
	case "phone":
		print(msgText)
	case "email":
		print(msgText)
	}

	return http.StatusOK, resp.Timeout, nil
}
