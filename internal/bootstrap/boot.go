package bootstrap

import (
	"payment-engine/internal/adaptors/billing"
	"payment-engine/internal/adaptors/clients"
	"payment-engine/internal/adaptors/encoders/avro"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/http"
	"payment-engine/pkg/db"
	"payment-engine/pkg/log"
	"payment-engine/pkg/metrics"

	"github.com/recodextech/container"
)

func Boot() {
	con := container.NewContainer()

	// Application config bindings
	err := con.SetModuleGlobalConfig(
		container.ModuleConfig{Key: application.ModuleBaseLogger, Value: new(log.LoggerConf)},
		container.ModuleConfig{Key: application.ModuleMetricsReporter, Value: new(metrics.Config)},
		container.ModuleConfig{Key: application.ModuleSchemaRegistry, Value: new(avro.SchemaRegistryConfig)},
		container.ModuleConfig{Key: application.ModuleHTTPRouter, Value: new(http.KRouterConf)},
		container.ModuleConfig{Key: application.ModuleHTTPServer, Value: new(http.Conf)},
		container.ModuleConfig{Key: application.ModuleSQL, Value: new(db.DatabaseConfig)},
		container.ModuleConfig{Key: application.ModuleBillingConfig, Value: new(billing.Config)},
		container.ModuleConfig{Key: application.ModuleHTTPClient, Value: new(clients.Config)},
	)
	if err != nil {
		panic(err)
	}

	bind(con)
	initModules(con)
	start(con)
}
