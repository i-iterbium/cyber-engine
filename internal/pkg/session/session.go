package session

import (
	"database/sql"
	"errors"
	"net/http"
)

// UserSession описывает модель пользвоательской сессии
type UserSession struct {
	UserID    *int64 `json:"userID,omitempty" sql:"user_id"`
	UserRole  string `json:"userRole,omitempty" sql:"user_role"`
	Session   string `json:"session,omitempty" sql:"session"`
	CSRFToken string `json:"CRSFToken,omitempty" sql:"crsf_token"`
}

// Get проверяет данные сессии, привязанной к переданному идентификатору Session, и возвращает userID
func Get(db *sql.DB, r *http.Request) (int, *UserSession, error) {
	sessionID, err := r.Cookie("session")
	if err != nil || sessionID.Value == "" {
		return 0, nil, nil
	}
	CSRFToken := r.Header.Get("CSRF-Token")
	if CSRFToken == "" {
		return 0, nil, nil
	}

	var s UserSession
	err = db.QueryRow(`select * from fn_check_user_session($1, $2)`, sessionID.Value, CSRFToken).Scan(&s.UserID, &s.UserRole)
	if err != nil {
		if err.Error() == "pq: Время жизни сессии пользователя истекло" {
			return 472, nil, err
		}
		return http.StatusBadRequest, nil, err
	}
	if s.UserID == nil {
		return http.StatusBadRequest, nil, errors.New("Пользователь не авторизован")
	}
	// if s.UserRole == "" {
	// 	return http.StatusBadRequest, UserSession{}, errors.New("Ошибка получения пользовательской роли")
	// }
	s.Session = sessionID.Value
	s.CSRFToken = CSRFToken

	return http.StatusOK, &s, nil
}
