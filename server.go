package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kawabatas/m-bank/gen/restapi"
	"github.com/kawabatas/m-bank/gen/restapi/operations"
	"github.com/kawabatas/m-bank/gen/restapi/operations/bank"
)

func newServer(db *sql.DB) (*restapi.Server, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	api := operations.NewBankAPI(swaggerSpec)
	server := restapi.NewServer(api)
	api.Logger = log.Printf

	app := newApp(db)
	setHandler(api, app)
	server.SetAPI(api)

	api.Middleware = func(middleware.Builder) http.Handler {
		return recoveryMiddleware(corsMiddleware(accessLogMiddleware(server.GetHandler())))
	}
	server.ConfigureAPI()

	return server, nil
}

func setHandler(api *operations.BankAPI, app *application) {
	// ctx := context.Background()
	api.BankGetBalanceHandler = bank.GetBalanceHandlerFunc(func(params bank.GetBalanceParams) middleware.Responder {
		return middleware.NotImplemented("operation bank.GetBalance has not yet been implemented")
	})

	api.BankPaymentTryHandler = bank.PaymentTryHandlerFunc(func(params bank.PaymentTryParams) middleware.Responder {
		return middleware.NotImplemented("operation bank.PaymentTry has not yet been implemented")
	})
	api.BankPaymentConfirmHandler = bank.PaymentConfirmHandlerFunc(func(params bank.PaymentConfirmParams) middleware.Responder {
		return middleware.NotImplemented("operation bank.PaymentConfirm has not yet been implemented")
	})
	api.BankPaymentCancelHandler = bank.PaymentCancelHandlerFunc(func(params bank.PaymentCancelParams) middleware.Responder {
		return middleware.NotImplemented("operation bank.PaymentCancel has not yet been implemented")
	})

	api.BankPaymentAllHandler = bank.PaymentAllHandlerFunc(func(params bank.PaymentAllParams) middleware.Responder {
		return middleware.NotImplemented("operation bank.PaymentAll has not yet been implemented")
	})
}
