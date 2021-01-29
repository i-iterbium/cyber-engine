package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/i-iterbium/cyber-engine/internal/pkg/repository"
	"github.com/i-iterbium/cyber-engine/models"
	"github.com/i-iterbium/cyber-engine/restapi/operations/user_confirmation"
)

// UserConfirmation подтверждает пользователя по коду
func UserConfirmation(repo repository.UserConfirmation) func(params user_confirmation.UserConfirmationParams) middleware.Responder {
	return func(params user_confirmation.UserConfirmationParams) middleware.Responder {
		status, _, err := repo.Confirm(params.HTTPRequest, params.UserConfirmation)
		if err != nil {
			if status == 400 {
				code, message := models.ErrCreateSession, err.Error()
				return user_confirmation.NewUserConfirmationBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 500 {
				code, message := models.ErrCreateSession, err.Error()
				return user_confirmation.NewUserConfirmationInternalServerError().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 472 {
				code, message := models.ErrSessionOutdate, err.Error()
				return user_confirmation.NewUserConfirmationBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
		}
		return user_confirmation.NewUserConfirmationOK()
	}
}

// UserConfirmationResendCode отправляет код подтверждения
func UserConfirmationResendCode(repo repository.UserConfirmation) func(params user_confirmation.ResendCodeParams) middleware.Responder {
	return func(params user_confirmation.ResendCodeParams) middleware.Responder {
		status, _, err := repo.ResendCode(params.HTTPRequest, params.ResendCode)
		if err != nil {
			if status == 400 {
				code, message := models.ErrCreateSession, err.Error()
				return user_confirmation.NewResendCodeBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 500 {
				code, message := models.ErrCreateSession, err.Error()
				return user_confirmation.NewResendCodeInternalServerError().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 472 {
				code, message := models.ErrSessionOutdate, err.Error()
				return user_confirmation.NewResendCodeBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
		}
		return user_confirmation.NewResendCodeOK()
	}
}
