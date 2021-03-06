package app

import (
	"database/sql/driver"

	"github.com/org39/webapp-tutorial-backend/usecase/auth"
	"github.com/org39/webapp-tutorial-backend/usecase/todo"
	"github.com/org39/webapp-tutorial-backend/usecase/user"

	"github.com/facebookgo/inject"
	"github.com/org39/webapp-tutorial-backend/pkg/db"
)

var DepencencyInjector inject.Graph

type App struct {
	// infra
	Config *Config `inject:""`
	DB     *db.DB  `inject:""`

	// application usecase
	AuthUsecase auth.Usecase `inject:""`
	UserUsecase user.Usecase `inject:""`
	TodoUsecase todo.Usecase `inject:""`
}

func New(dbConnectorFn func(*Config) (driver.Connector, error)) (*App, error) {
	if err := newInfra(dbConnectorFn); err != nil {
		return nil, err
	}

	if err := newUserUsecase(); err != nil {
		return nil, err
	}

	if err := newAuthUsecase(); err != nil {
		return nil, err
	}

	if err := newTodoUsecase(); err != nil {
		return nil, err
	}

	app := new(App)
	err := DepencencyInjector.Provide(
		&inject.Object{Value: app},
	)
	if err != nil {
		return nil, err
	}

	err = DepencencyInjector.Populate()
	if err != nil {
		return nil, err
	}

	return app, nil
}

func ClearDepencencyGraph() {
	DepencencyInjector = inject.Graph{}
}
