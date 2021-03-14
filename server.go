package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kawabatas/m-bank/domain/model"
	"github.com/kawabatas/m-bank/gen/models"
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
	ctx := context.Background()
	api.BankGetBalanceHandler = bank.GetBalanceHandlerFunc(func(params bank.GetBalanceParams) middleware.Responder {
		balance, err := app.BalanceService.Get(ctx, uint(params.UserID))
		if err != nil {
			// TODO: 適切なエラーコードを返す
			return bank.NewGetBalanceDefault(500).WithPayload(toErrorResponse(500, err.Error()))
		}
		return bank.NewGetBalanceOK().WithPayload(&models.Balance{UserID: int32(balance.UserID), Amount: int32(balance.Amount)})
	})

	api.BankPaymentTryHandler = bank.PaymentTryHandlerFunc(func(params bank.PaymentTryParams) middleware.Responder {
		pt, balance, err := app.PaymentService.Try(ctx, *params.Body.IdempotencyKey, uint(*params.Body.UserID), int(params.Body.Amount))
		if err != nil {
			return bank.NewPaymentTryDefault(500).WithPayload(toErrorResponse(500, err.Error()))
		}
		return bank.NewPaymentTryOK().WithPayload(toPayResponse(pt, balance))
	})
	api.BankPaymentConfirmHandler = bank.PaymentConfirmHandlerFunc(func(params bank.PaymentConfirmParams) middleware.Responder {
		pt, balance, err := app.PaymentService.Confirm(ctx, *params.Body.IdempotencyKey, uint(*params.Body.UserID), int(params.Body.Amount))
		if err != nil {
			return bank.NewPaymentConfirmDefault(500).WithPayload(toErrorResponse(500, err.Error()))
		}
		return bank.NewPaymentConfirmOK().WithPayload(toPayResponse(pt, balance))
	})
	api.BankPaymentCancelHandler = bank.PaymentCancelHandlerFunc(func(params bank.PaymentCancelParams) middleware.Responder {
		pt, balance, err := app.PaymentService.Cancel(ctx, *params.Body.IdempotencyKey, uint(*params.Body.UserID), int(params.Body.Amount))
		if err != nil {
			return bank.NewPaymentCancelDefault(500).WithPayload(toErrorResponse(500, err.Error()))
		}
		return bank.NewPaymentCancelOK().WithPayload(toPayResponse(pt, balance))
	})

	api.BankPaymentAddToUsersHandler = bank.PaymentAddToUsersHandlerFunc(func(params bank.PaymentAddToUsersParams) middleware.Responder {
		return middleware.NotImplemented("operation bank.PaymentAddToUsers has not yet been implemented")
	})
}

func toPayResponse(pt *model.PaymentTransaction, balance *model.Balance) *models.PayResponse {
	return &models.PayResponse{
		IdempotencyKey: pt.UUID,
		TryTime:        strfmt.DateTime(pt.TryTime),
		ConfirmTime:    strfmt.DateTime(pt.ConfirmTime),
		CancelTime:     strfmt.DateTime(pt.CancelTime),
		Balance: &models.Balance{
			UserID: int32(balance.UserID),
			Amount: int32(balance.Amount),
		},
	}
}

func toErrorResponse(c int, m string) *models.ErrorResponse {
	code := int32(c)
	return &models.ErrorResponse{
		Code:    code,
		Message: m,
	}
}
