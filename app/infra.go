package app

import (
	"database/sql/driver"
	"net"
	"time"

	"github.com/org39/webapp-tutorial-backend/pkg/db"
	"github.com/org39/webapp-tutorial-backend/pkg/log"

	"github.com/facebookgo/inject"
	"github.com/go-sql-driver/mysql"
)

func newInfra() error {
	// application config
	conf, err := NewConfig()
	if err != nil {
		return err
	}

	// logger
	logger := log.Wrap(log.Log.WithField("service_name", conf.ServieName))

	// database
	mysqlConn, err := newMysqlConn(conf)
	if err != nil {
		return err
	}
	database, err := db.New(mysqlConn)
	if err != nil {
		return err
	}

	// build depency graph
	err = DepencencyInjector.Provide(
		&inject.Object{Value: conf},
		&inject.Object{Value: logger},
		&inject.Object{Value: database},
		&inject.Object{Name: "repo.user.table", Value: conf.UserTable},
		&inject.Object{Name: "usecase.auth.secret", Value: conf.AuthSecret},
		&inject.Object{Name: "usecase.auth.access_token_duration", Value: conf.AuthAccessTokenDuration},
		&inject.Object{Name: "usecase.auth.refresh_token_duration", Value: conf.AuthRefereshTokenDuration},
	)
	if err != nil {
		return err
	}

	return nil
}

func newMysqlConn(conf *Config) (driver.Connector, error) {
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}

	databaseHost := net.JoinHostPort(conf.DatabaseHost, conf.DatabasePort)
	dsn := &mysql.Config{
		Addr:                 databaseHost,
		Net:                  "tcp",
		User:                 conf.DatabaseUser,
		Passwd:               conf.DatabasePass,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  utc,
		ParseTime:            true,
		DBName:               conf.DatabaseName,
		AllowNativePasswords: true,
	}

	return mysql.NewConnector(dsn)
}
