package http

import (
	"payment-engine/internal/domain/application"
	"payment-engine/internal/http/handlers"
	"payment-engine/internal/http/responses/writers"
	"payment-engine/internal/http/validator"

	"github.com/recodextech/container"
)

type HTTP struct{}

// Init initializes the http module.
func (h *HTTP) Init(c container.Container) error {
	// Http validators
	c.Bind(validator.ModuleAccountIDValidator, new(validator.AccountValidator))
	c.Bind(validator.ModuleJobIDValidator, new(validator.JobValidator))
	c.Bind(validator.ModulePaymentIDValidator, new(validator.PaymentValidator))

	// Http handlers
	c.Bind(handlers.ModuleCreatePayment, new(handlers.CreatePaymentHandler))
	c.Bind(handlers.ModuleCancelPayment, new(handlers.CancelPaymentHandler))
	c.Bind(handlers.ModuleSuccessPayment, new(handlers.SuccessPaymentHandler))
	c.Bind(handlers.ModuleCreateWallet, new(handlers.CreateWalletHandler))
	c.Bind(handlers.ModuleGetWallets, new(handlers.GetWalletsHandler))

	// Http error handler
	c.Bind(handlers.ModuleErrorHandler, new(handlers.ErrorHandler))

	// HTTP writer
	c.Bind(writers.ModuleJobListWriter, new(writers.JobListWriter))
	c.Bind(writers.ModulePaymentWriter, new(writers.PaymentWriter))
	c.Bind(writers.ModuleNoContentWriter, new(writers.NoContentWriter))
	c.Bind(writers.ModuleWalletWriter, new(writers.WalletWriter))
	c.Bind(writers.ModuleInternalWalletWriter, new(writers.WalletsInternalWriter))
	c.Bind(writers.ModuleGetWalletsWriter, new(writers.GetWalletsWriter))

	c.Bind(application.ModuleHTTPRouter, new(Router))
	c.Bind(application.ModuleHTTPServer, new(Server))

	c.Init(
		// Http request validators
		validator.ModuleAccountIDValidator,
		validator.ModuleJobIDValidator,
		validator.ModulePaymentIDValidator,

		// Http handlers
		handlers.ModuleCreatePayment,
		handlers.ModuleCancelPayment,
		handlers.ModuleSuccessPayment,
		handlers.ModuleCreateWallet,
		handlers.ModuleGetWallets,

		// Http error handler
		handlers.ModuleErrorHandler,
		// HTTP writers
		writers.ModuleJobListWriter,
		writers.ModulePaymentWriter,
		writers.ModuleNoContentWriter,
		writers.ModuleWalletWriter,
		writers.ModuleInternalWalletWriter,
		writers.ModuleGetWalletsWriter,
		writers.ModuleNoContentWriter,

		// Http Server Init
		application.ModuleHTTPRouter,
		application.ModuleHTTPServer,
	)

	return nil
}
