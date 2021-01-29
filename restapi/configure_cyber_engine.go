// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/i-iterbium/cyber-engine/restapi/operations/user_confirmation"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/i-iterbium/cyber-engine/internal/handlers"
	"github.com/i-iterbium/cyber-engine/internal/pkg/database"
	"github.com/i-iterbium/cyber-engine/restapi/operations"
	"github.com/i-iterbium/cyber-engine/restapi/operations/sessions"
	"github.com/i-iterbium/cyber-engine/restapi/operations/users"
)

//go:generate swagger generate server --target ..\..\cyber-engine --name CyberEngine --spec ..\api\openapi-spec\swagger.yml

func configureFlags(api *operations.CyberEngineAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.CyberEngineAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	repos, err := configureRepositories(database.Open("pgCyberEngine"))
	if err != nil {
		return nil
	}

	api.SessionsCreateSessionHandler = sessions.CreateSessionHandlerFunc(handlers.CreateSession(repos.Sessions))
	api.SessionsUpdateSessionHandler = sessions.UpdateSessionHandlerFunc(handlers.UpdateSession(repos.Sessions))
	api.SessionsDeleteSessionHandler = sessions.DeleteSessionHandlerFunc(handlers.DeleteSession(repos.Sessions))
	api.UsersCreateUserHandler = users.CreateUserHandlerFunc(handlers.CreateUser(repos.Users))
	api.UsersFetchUserByIDHandler = users.FetchUserByIDHandlerFunc(handlers.FetchUserByID(repos.Users))
	api.UsersUpdateUserByIDHandler = users.UpdateUserByIDHandlerFunc(handlers.UpdateUser(repos.Users))
	api.UserConfirmationUserConfirmationHandler = user_confirmation.UserConfirmationHandlerFunc(handlers.UserConfirmation(repos.UserConfirmation))
	api.UserConfirmationResendCodeHandler = user_confirmation.ResendCodeHandlerFunc(handlers.UserConfirmationResendCode(repos.UserConfirmation))

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
