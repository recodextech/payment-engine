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
		application.ModulePlanRepo,
		application.ModuleSubscriptionRepo,
		application.ModuleContractorAccountRepo,
		application.ModuleNotificationClient,

		// Services
		application.ModulePaymentService,
		application.ModuleWalletService,
		application.ModuleSubscriptionNotifier,
		application.ModuleSubscriptionCore,
		application.ModuleSubscriptionService,
		application.ModuleQuotaService,

		// Http
		application.ModuleHTTP,
	)
}
