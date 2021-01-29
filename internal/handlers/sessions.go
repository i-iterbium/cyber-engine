package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/i-iterbium/cyber-engine/internal/pkg/repository"
	"github.com/i-iterbium/cyber-engine/models"
	"github.com/i-iterbium/cyber-engine/restapi/operations/sessions"
)

// CreateSession создает пользовательскую сессию
func CreateSession(repo repository.Sessions) func(params sessions.CreateSessionParams) middleware.Responder {
	return func(params sessions.CreateSessionParams) middleware.Responder {
		status, payload, err := repo.Create(params.HTTPRequest, params.AuthData)
		if err != nil {
			if status == 400 {
				code, message := models.ErrCreateSession, err.Error()
				return sessions.NewCreateSessionBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 500 {
				code, message := models.ErrCreateSession, err.Error()
				return sessions.NewCreateSessionInternalServerError().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 472 {
				code, message := models.ErrSessionOutdate, err.Error()
				return sessions.NewCreateSessionBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
		}
		return sessions.NewCreateSessionOK().WithPayload(payload)
	}
}

// UpdateSession обновляет пользовательскую сессию
func UpdateSession(repo repository.Sessions) func(params sessions.UpdateSessionParams) middleware.Responder {
	return func(params sessions.UpdateSessionParams) middleware.Responder {
		status, payload, err := repo.Update(params.HTTPRequest)
		if err != nil {
			if status == 400 {
				code, message := models.ErrUpdateSession, err.Error()
				return sessions.NewUpdateSessionBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 500 {
				code, message := models.ErrUpdateSession, err.Error()
				return sessions.NewUpdateSessionInternalServerError().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 472 {
				code, message := models.ErrSessionOutdate, err.Error()
				return sessions.NewUpdateSessionBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
		}
		return sessions.NewCreateSessionOK().WithPayload(payload)
	}
}

// DeleteSession удаляет пользовательскую сессию
func DeleteSession(repo repository.Sessions) func(params sessions.DeleteSessionParams) middleware.Responder {
	return func(params sessions.DeleteSessionParams) middleware.Responder {
		status, err := repo.Delete(params.HTTPRequest)
		if err != nil {
			if status == 400 {
				code, message := models.ErrDeleteSession, err.Error()
				return sessions.NewDeleteSessionBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 500 {
				code, message := models.ErrDeleteSession, err.Error()
				return sessions.NewDeleteSessionInternalServerError().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 472 {
				code, message := models.ErrSessionOutdate, err.Error()
				return sessions.NewDeleteSessionBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
		}
		return sessions.NewDeleteSessionOK()
	}
}
