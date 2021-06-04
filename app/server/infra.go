package app

import (
	"database/sql/driver"

	"github.com/org39/webapp-tutorial-backend/pkg/db"
	"github.com/org39/webapp-tutorial-backend/pkg/log"

	"github.com/facebookgo/inject"
)

func newInfra(dbConnectorFn func(*Config) (driver.Connector, error)) error {
	// application config
	conf, err := NewConfig()
	if err != nil {
		return err
	}

	// set loglevel
	log.SetLevel(conf.LogLevel)

	// database
	dbConn, err := dbConnectorFn(conf)
	if err != nil {
		return err
	}
	database, err := db.New(dbConn)
	if err != nil {
		return err
	}

	// build depency graph
	err = DepencencyInjector.Provide(
		&inject.Object{Value: conf},
		&inject.Object{Value: database},
		&inject.Object{Name: "repo.user.table", Value: conf.UserTable},
		&inject.Object{Name: "repo.todo.table", Value: conf.TodoTable},
		&inject.Object{Name: "usecase.user.password_salt", Value: conf.UserPasswordSalt},
		&inject.Object{Name: "usecase.auth.secret", Value: conf.AuthSecret},
		&inject.Object{Name: "usecase.auth.access_token_duration", Value: conf.AuthAccessTokenDuration},
		&inject.Object{Name: "usecase.auth.refresh_token_duration", Value: conf.AuthRefreshTokenDuration},
		&inject.Object{Name: "rest.auth.secure_refresh_token", Value: conf.RestAuthSecureRefreshToken},
	)
	if err != nil {
		return err
	}

	return nil
}
