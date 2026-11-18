package bootstrap

import (
	"payment-engine/internal/adaptors/clients"
	"payment-engine/internal/adaptors/database/postgres"
	"payment-engine/internal/adaptors/repositories"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/http"
	"payment-engine/internal/services"
	db "payment-engine/pkg/db"
	log2 "payment-engine/pkg/log"
	"payment-engine/pkg/metrics"

	"github.com/recodextech/container"
)

func bind(c container.Container) {
	c.Bind(application.ModuleBaseLogger, new(log2.BaseLogger))
	c.Bind(application.ModuleLogger, new(log2.Logger))
	c.Bind(application.ModuleMetricsReporter, new(metrics.Reporter))

	// Repositories
	c.Bind(application.ModuleAccountRepo, new(repositories.AccountRepository))
	c.Bind(application.ModuleAccountWalletRepo, new(repositories.AccountWalletRepository))
	c.Bind(application.ModulePaymentRepo, new(repositories.PaymentRepository))
	c.Bind(application.ModuleJobRepo, new(repositories.JobRepository))
	c.Bind(application.ModulePlanRepo, new(repositories.PlanRepository))
	c.Bind(application.ModuleSubscriptionRepo, new(repositories.SubscriptionRepository))
	c.Bind(application.ModuleContractorAccountRepo, new(repositories.ContractorAccountRepository))

	// Adapters
	c.Bind(application.ModuleSQL, new(db.PostgresDB))
	c.Bind(application.MoudleDBConector, new(postgres.DBConnector))
	c.Bind(application.ModuleNotificationClient, new(clients.NotificationServiceClient))

	// Services
	c.Bind(application.ModulePaymentService, new(services.PaymentService))
	c.Bind(application.ModuleWalletService, new(services.WalletService))
	c.Bind(application.ModuleSubscriptionNotifier, new(services.SubscriptionNotifier))
	c.Bind(application.ModuleSubscriptionCore, new(services.SubscriptionCore))
	c.Bind(application.ModuleSubscriptionService, new(services.SubscriptionService))
	c.Bind(application.ModuleQuotaService, new(services.QuotaService))

	// http
	c.Bind(application.ModuleHTTP, new(http.HTTP))
}
