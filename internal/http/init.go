package http

import (
	"payment-engine/internal/domain/application"
	"payment-engine/internal/http/handlers"
	"payment-engine/internal/http/responses/writers"

	"github.com/recodextech/container"
)

type HTTP struct{}

// Init initializes the http module.
func (h *HTTP) Init(c container.Container) error {
	// Http validators

	// Http handlers
	c.Bind(handlers.ModuleGetJobsByWorker, new(handlers.GetJobsByWorkerHandler))

	// Http error handler
	c.Bind(handlers.ModuleErrorHandler, new(handlers.ErrorHandler))

	// HTTP writer
	c.Bind(writers.ModuleJobListWriter, new(writers.JobListWriter))

	c.Bind(application.ModuleHTTPRouter, new(Router))
	c.Bind(application.ModuleHTTPServer, new(Server))

	c.Init(
		// Http request validators

		// Http handlers
		handlers.ModuleGetJobsByWorker,

		// Http error handler
		handlers.ModuleErrorHandler,
		// HTTP writers
		writers.ModuleJobListWriter,

		// Http Server Init
		application.ModuleHTTPRouter,
		application.ModuleHTTPServer,
	)

	return nil
}
