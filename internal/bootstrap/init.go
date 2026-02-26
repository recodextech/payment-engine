package bootstrap

import (
	"payment-engine/internal/domain/application"

	"github.com/recodextech/container"
)

func initModules(c container.Container) {
	c.Init(
		// main application dependencies
		application.ModuleBaseLogger,
		application.ModuleLogger,
		application.ModuleMetricsReporter,
		application.ModuleReadyIndicator,
		application.ModulePprofIndicator,

		// Repositories
		application.ModuleAccountRepo,
		application.ModuleAccountWalletRepo,
		application.ModulePaymentRepo,
		application.ModuleJobRepo,

		// Adapters
		application.ModuleSQL,
		application.MoudleDBConector,

		// Services
		application.ModulePaymentService,
		application.ModuleWalletService,

		// Http
		application.ModuleHTTP,
	)
}
