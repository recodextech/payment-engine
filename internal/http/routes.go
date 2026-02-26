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
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleAccountIDValidator).(krouter.Validator)),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleJobIDValidator).(krouter.Validator)),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModulePaymentWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPost)

	// Cancel payment endpoint
	r.router.Handle(fmt.Sprintf("/jobs/{%s}/payment/{%s}/cancelled", request.PathParamJobID.String(), request.PathParamPaymentID.String()),
		r.krouter.NewHandler(
			handlers.ModuleCancelPayment,
			request.NoOp{},
			c.Resolve(handlers.ModuleCancelPayment).(*handlers.CancelPaymentHandler).Handle,
			krouter.HandlerWithParameter(request.PathParamJobID.String(), request.ParamTypeAppUUID),
			krouter.HandlerWithParameter(request.PathParamPaymentID.String(), request.ParamTypeAppUUID),
			krouter.HandlerWithHeader(request.HeaderAccountID.String(), request.ParamTypeAppUUID, nil),
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleAccountIDValidator).(krouter.Validator)),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleJobIDValidator).(krouter.Validator)),
			krouter.HandlerWithValidator(c.Resolve(validator.ModulePaymentIDValidator).(krouter.Validator)),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleNoContentWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPatch)

	// Success payment endpoint
	r.router.Handle(fmt.Sprintf("/jobs/{%s}/payment/{%s}/success", request.PathParamJobID.String(), request.PathParamPaymentID.String()),
		r.krouter.NewHandler(
			handlers.ModuleSuccessPayment,
			request.NoOp{},
			c.Resolve(handlers.ModuleSuccessPayment).(*handlers.SuccessPaymentHandler).Handle,
			krouter.HandlerWithParameter(request.PathParamJobID.String(), request.ParamTypeAppUUID),
			krouter.HandlerWithParameter(request.PathParamPaymentID.String(), request.ParamTypeAppUUID),
			krouter.HandlerWithHeader(request.HeaderAccountID.String(), request.ParamTypeAppUUID, nil),
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleAccountIDValidator).(krouter.Validator)),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleJobIDValidator).(krouter.Validator)),
			krouter.HandlerWithValidator(c.Resolve(validator.ModulePaymentIDValidator).(krouter.Validator)),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleNoContentWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPatch)

	// Create wallet endpoint
	r.router.Handle("/wallet/create",
		r.krouter.NewHandler(
			handlers.ModuleCreateWallet,
			request.CreateWalletRequest{},
			c.Resolve(handlers.ModuleCreateWallet).(*handlers.CreateWalletHandler).Handle,
			krouter.HandlerWithHeader(request.HeaderAccountID.String(), request.ParamTypeAppUUID, nil),
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleAccountIDValidator).(krouter.Validator)),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleWalletWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPost)

	// Create internal wallet endpoint
	r.router.Handle("/internal-wallet/create",
		r.krouter.NewHandler(
			handlers.ModuleCreateInternalWallet,
			request.NoOp{},
			c.Resolve(handlers.ModuleCreateWallet).(*handlers.CreateWalletHandler).HandleInternalWallet,
			krouter.HandlerWithHeader(request.HeaderAccountID.String(), request.ParamTypeAppUUID, nil),
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleAccountIDValidator).(krouter.Validator)),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleInternalWalletWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPost)

	// Get wallets endpoint
	r.router.Handle("/wallets",
		r.krouter.NewHandler(
			handlers.ModuleGetWallets,
			request.NoOp{},
			c.Resolve(handlers.ModuleGetWallets).(*handlers.GetWalletsHandler).Handle,
			krouter.HandlerWithHeader(request.HeaderAccountID.String(), request.ParamTypeAppUUID, nil),
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleAccountIDValidator).(krouter.Validator)),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleGetWalletsWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodGet)
}
