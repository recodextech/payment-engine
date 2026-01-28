package db

import (
	"payment-engine/internal/domain/application"

	"github.com/recodextech/container"

	gopostgres "github.com/HADLakmal/go-postgres"
)

type PostgresDB struct {
	conf *DatabaseConfig
	gopostgres.DatabaseReporter
}

// New creates a new Postgres adapter instance.
func (a *PostgresDB) Init(c container.Container) (err error) {
	a.conf = c.GetGlobalConfig(application.ModuleSQL).(*DatabaseConfig)

	conf := gopostgres.DBConfig{
		User:     a.conf.User,
		Password: a.conf.Password,
		Host:     a.conf.Host,
		Port:     a.conf.Port,
		Database: a.conf.Database,
	}
	a.DatabaseReporter, err = gopostgres.NewDBPool(conf)
	if err != nil {
		return err
	}
	return
}

// Stop will close the Postgres adapter releasing all resources.
func (a *PostgresDB) Stop() error {
	return a.DatabaseReporter.Close()
}
