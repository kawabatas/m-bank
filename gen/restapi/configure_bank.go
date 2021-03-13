// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/kawabatas/m-bank/gen/restapi/operations"
	"github.com/kawabatas/m-bank/gen/restapi/operations/bank"
)

//go:generate swagger generate server --target ../../gen --name Bank --spec ../../swagger.yml --principal interface{} --exclude-main

func configureFlags(api *operations.BankAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.BankAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.BankGetBalanceHandler == nil {
		api.BankGetBalanceHandler = bank.GetBalanceHandlerFunc(func(params bank.GetBalanceParams) middleware.Responder {
			return middleware.NotImplemented("operation bank.GetBalance has not yet been implemented")
		})
	}
	if api.BankPaymentAllHandler == nil {
		api.BankPaymentAllHandler = bank.PaymentAllHandlerFunc(func(params bank.PaymentAllParams) middleware.Responder {
			return middleware.NotImplemented("operation bank.PaymentAll has not yet been implemented")
		})
	}
	if api.BankPaymentCancelHandler == nil {
		api.BankPaymentCancelHandler = bank.PaymentCancelHandlerFunc(func(params bank.PaymentCancelParams) middleware.Responder {
			return middleware.NotImplemented("operation bank.PaymentCancel has not yet been implemented")
		})
	}
	if api.BankPaymentConfirmHandler == nil {
		api.BankPaymentConfirmHandler = bank.PaymentConfirmHandlerFunc(func(params bank.PaymentConfirmParams) middleware.Responder {
			return middleware.NotImplemented("operation bank.PaymentConfirm has not yet been implemented")
		})
	}
	if api.BankPaymentTryHandler == nil {
		api.BankPaymentTryHandler = bank.PaymentTryHandlerFunc(func(params bank.PaymentTryParams) middleware.Responder {
			return middleware.NotImplemented("operation bank.PaymentTry has not yet been implemented")
		})
	}

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
