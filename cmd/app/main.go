package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/d-kolpakov/logger"
	"github.com/dhnikolas/configo"
	"github.com/pressly/goose"
	_ "rebrain-location/internal/migrations"
	"rebrain-location/internal/routes"
	"rebrain-location/pkg/helpers/pg"
)

const ServiceName = "rebrain-location"

var AppVersion string

func main() {
	fmt.Println("Starting service " + ServiceName)

	//Logger initialization
	logDriver := &logger.STDOUTDriver{}
	loggerConfig := logger.LoggerConfig{
		ServiceName: ServiceName,
		Level:       configo.EnvInt("logging-level", 2),
		Buffer:      configo.EnvInt("app-logger-buffer-size", 10000),
		Output:      []logger.LogDriver{logDriver},
	}
	l, err := logger.GetLogger(loggerConfig)
	if err != nil {
		panic(err)
	}

	cfg := &pg.Config{}
	cfg.Host = configo.EnvString("db-host", "")
	cfg.Username = configo.EnvString("db-username", "")
	cfg.Password = configo.EnvString("db-password", "")
	cfg.Port = configo.EnvString("db-port", "")
	cfg.DbName = configo.EnvString("db-name", "")
	cfg.Timeout = 5
	config, err := pg.NewPoolConfig(cfg)
	if err != nil {
		panic(err)
	}
	config.MaxConns = 500
	c, err := pg.NewConnection(config)
	if err != nil {
		panic(err)
	}

	mdb, err := sql.Open("postgres", config.ConnString())
	if err != nil {
		panic(err)
	}
	err = goose.Up(mdb, "/var")
	if err != nil {
		l.ErrorWithContext(context.Background(), err)
	}

	routes.Start(ServiceName, AppVersion, l, c)
}
