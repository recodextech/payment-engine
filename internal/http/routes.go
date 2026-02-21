package http

import (
	"fmt"
	"net/http"

	"payment-engine/internal/http/handlers"
	"payment-engine/internal/http/request"
	"payment-engine/internal/http/responses"
	"payment-engine/internal/http/responses/writers"
	"payment-engine/internal/http/validator"

	"github.com/recodextech/container"

	"github.com/recodextech/krouter"
)

func initRoutes(r *Router, c container.Container) {
	// Create payment endpoint
	r.router.Handle(fmt.Sprintf("/jobs/{%s}/payment/create", request.PathParamJobID.String()),
		r.krouter.NewHandler(
			handlers.ModuleCreatePayment,
			request.CreatePaymentRequest{},
			c.Resolve(handlers.ModuleCreatePayment).(*handlers.CreatePaymentHandler).Handle,
			krouter.HandlerWithParameter(request.PathParamJobID.String(), request.ParamTypeAppUUID),
			krouter.HandlerWithHeader(request.HeaderAccountID.String(), request.ParamTypeAppUUID, nil),
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleAccountIDVaidator).(krouter.Validator)),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModulePaymentWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPost)

	// Create wallet endpoint
	r.router.Handle("/wallet/create",
		r.krouter.NewHandler(
			handlers.ModuleCreateWallet,
			request.NoOp{},
			c.Resolve(handlers.ModuleCreateWallet).(*handlers.CreateWalletHandler).Handle,
			krouter.HandlerWithHeader(request.HeaderAccountID.String(), request.ParamTypeAppUUID, nil),
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleAccountIDVaidator).(krouter.Validator)),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleWalletWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPost)
}
