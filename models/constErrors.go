package models

// Основы ошибок модулей
const (
	ErrUsers int64 = (iota + 1) * 1000
	ErrSessions
)

// Ошибки модуля пользовательских сессий
const (
	ErrSessionOutdate = ErrSessions + iota + 1
	ErrCreateSession
	ErrUpdateSession
	ErrDeleteSession
)

// Ошибки модуля работы с пользователями
const (
	ErrFetchUserByID = ErrUsers + iota + 1
	ErrCreateUser
	ErrUpdateUser
	ErrCreateCode
	ErrUpdateConfirmedUser
	ErrFetchUpdatedUsers
)
