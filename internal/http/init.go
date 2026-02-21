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
	c.Bind(validator.ModuleAccountIDVaidator, new(validator.AccountValidator))

	// Http handlers
	c.Bind(handlers.ModuleCreatePayment, new(handlers.CreatePaymentHandler))
	c.Bind(handlers.ModuleCreateWallet, new(handlers.CreateWalletHandler))

	// Http error handler
	c.Bind(handlers.ModuleErrorHandler, new(handlers.ErrorHandler))

	// HTTP writer
	c.Bind(writers.ModuleJobListWriter, new(writers.JobListWriter))
	c.Bind(writers.ModulePaymentWriter, new(writers.PaymentWriter))
	c.Bind(writers.ModuleWalletWriter, new(writers.WalletWriter))

	c.Bind(application.ModuleHTTPRouter, new(Router))
	c.Bind(application.ModuleHTTPServer, new(Server))

	c.Init(
		// Http request validators
		validator.ModuleAccountIDVaidator,

		// Http handlers
		handlers.ModuleCreatePayment,
		handlers.ModuleCreateWallet,

		// Http error handler
		handlers.ModuleErrorHandler,
		// HTTP writers
		writers.ModuleJobListWriter,
		writers.ModulePaymentWriter,
		writers.ModuleWalletWriter,

		// Http Server Init
		application.ModuleHTTPRouter,
		application.ModuleHTTPServer,
	)

	return nil
}
