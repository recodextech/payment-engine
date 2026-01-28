package bootstrap

import (
	"payment-engine/internal/domain/application"
	"os"
	"os/signal"

	"github.com/recodextech/container"
)

func start(con container.AppContainer) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		con.Shutdown(
			application.ModuleHTTPServer,
			application.ModuleHTTPRouter,
			application.ModuleMetricsReporter,
		)
	}()

	con.Start(
		application.ModuleMetricsReporter,
		application.ModuleHTTPRouter,
		application.ModuleHTTPServer,
	)
}
