package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/i-iterbium/cyber-engine/internal/pkg/repository"
	"github.com/i-iterbium/cyber-engine/models"
	"github.com/i-iterbium/cyber-engine/restapi/operations/users"
)

// FetchUserByID возвращает пользователя по ID
func FetchUserByID(repo repository.Users) func(params users.FetchUserByIDParams) middleware.Responder {
	return func(params users.FetchUserByIDParams) middleware.Responder {
		status, payload, err := repo.Fetch(params.HTTPRequest, params.ID)
		if err != nil {
			if status == 400 {
				code, message := models.ErrFetchUserByID, err.Error()
				return users.NewFetchUserByIDBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 500 {
				code, message := models.ErrFetchUserByID, err.Error()
				return users.NewFetchUserByIDInternalServerError().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 472 {
				code, message := models.ErrSessionOutdate, err.Error()
				return users.NewFetchUserByIDBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
		}
		return users.NewFetchUserByIDOK().WithPayload(payload)
	}
}

// CreateUser создает пользователя
func CreateUser(repo repository.Users) func(params users.CreateUserParams) middleware.Responder {
	return func(params users.CreateUserParams) middleware.Responder {
		status, payload, err := repo.Create(params.HTTPRequest, params.User)
		if err != nil {
			if status == 400 {
				code, message := models.ErrCreateUser, err.Error()
				return users.NewCreateUserBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 500 {
				code, message := models.ErrCreateUser, err.Error()
				return users.NewCreateUserInternalServerError().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 472 {
				code, message := models.ErrSessionOutdate, err.Error()
				return users.NewCreateUserBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
		}
		return users.NewCreateUserOK().WithPayload(payload)
	}
}

// UpdateUser обновляет данные пользователя
func UpdateUser(repo repository.Users) func(params users.UpdateUserByIDParams) middleware.Responder {
	return func(params users.UpdateUserByIDParams) middleware.Responder {
		status, payload, err := repo.Patch(params.HTTPRequest, params.ID, params.User)
		if err != nil {
			if status == 400 {
				code, message := models.ErrUpdateUser, err.Error()
				return users.NewUpdateUserByIDBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 500 {
				code, message := models.ErrUpdateUser, err.Error()
				return users.NewUpdateUserByIDInternalServerError().WithPayload(&models.Error{Code: &code, Message: &message})
			}
			if status == 472 {
				code, message := models.ErrSessionOutdate, err.Error()
				return users.NewUpdateUserByIDBadRequest().WithPayload(&models.Error{Code: &code, Message: &message})
			}
		}
		return users.NewUpdateUserByIDOK().WithPayload(payload)
	}
}
