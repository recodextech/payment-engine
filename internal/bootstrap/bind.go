package bootstrap

import (
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
	c.Bind(application.ModuleJobRepo, new(repositories.JobRepository))

	// Adapters
	c.Bind(application.ModuleSQL, new(db.PostgresDB))
	c.Bind(application.MoudleDBConector, new(postgres.DBConnector))

	// Services
	c.Bind(application.ModuleJobService, new(services.JobService))

	// http
	c.Bind(application.ModuleHTTP, new(http.HTTP))
}
