package metrics

import (
	"context"
	"errors"
	"fmt"
	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/application"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/recodextech/container"
	"github.com/tryfix/metrics"
)

type Reporter struct {
	reporter metrics.Reporter
	http     *http.Server
	logger   adaptors.Logger
	conf     *Config
}

func (m *Reporter) Init(c container.Container) error {
	conf := c.GetGlobalConfig(application.ModuleMetricsReporter).(*Config)

	m.Reporter(adaptors.ReporterConf{
		System:      conf.System,
		Subsystem:   conf.SubSystem,
		ConstLabels: map[string]string{},
	})

	m.logger = c.Resolve(application.ModuleLogger).(adaptors.Logger).
		NewLog(adaptors.LoggerPrefixed("metrics.Reporter"))
	m.conf = conf
	m.http = &http.Server{
		Addr:              conf.HTTP.Host,
		ReadHeaderTimeout: conf.Timeouts.ReadHeader,
	}

	router := mux.NewRouter()
	router.Handle(fmt.Sprintf(`/%s`, conf.HTTP.Path), promhttp.Handler())
	m.http.Handler = router

	c.Bind(application.ModuleBaseReporter, m.reporter)

	return nil
}

func (m *Reporter) Run() error {
	m.logger.Info(fmt.Sprintf(`http server starting on %s/%s`, m.conf.HTTP.Host, m.conf.HTTP.Path))

	if err := m.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (m *Reporter) Ready() chan bool {
	return nil
}

func (m *Reporter) Stop() error {
	c, fn := context.WithTimeout(context.Background(), m.conf.HTTP.ShutdownWait)
	defer fn()

	return m.http.Shutdown(c)
}

func (m *Reporter) Reporter(conf adaptors.ReporterConf) adaptors.MetricsReporter {
	m.reporter = metrics.PrometheusReporter(metrics.ReporterConf{
		System:      conf.System,
		Subsystem:   conf.Subsystem,
		ConstLabels: conf.ConstLabels,
	})

	return m
}

func (m *Reporter) Counter(conf adaptors.MetricConf) adaptors.Counter {
	return m.reporter.Counter(metrics.MetricConf{
		Path:        conf.Path,
		Labels:      conf.Labels,
		ConstLabels: conf.ConstLabels,
	})
}

func (m *Reporter) Observer(conf adaptors.MetricConf) adaptors.Observer {
	return m.reporter.Observer(metrics.MetricConf{
		Path:        conf.Path,
		Labels:      conf.Labels,
		ConstLabels: conf.ConstLabels,
	})
}

func (m *Reporter) Gauge(conf adaptors.MetricConf) adaptors.Gauge {
	return m.reporter.Gauge(metrics.MetricConf{
		Path:        conf.Path,
		Labels:      conf.Labels,
		ConstLabels: conf.ConstLabels,
	})
}

func (m *Reporter) GaugeFunc(conf adaptors.MetricConf, f func() float64) adaptors.GaugeFunc {
	return m.reporter.GaugeFunc(metrics.MetricConf{
		Path:        conf.Path,
		Labels:      conf.Labels,
		ConstLabels: conf.ConstLabels,
	}, f)
}

func (m *Reporter) Info() string {
	return m.reporter.Info()
}

func (m *Reporter) UnRegister(metrics string) {
	m.reporter.UnRegister(metrics)
}
