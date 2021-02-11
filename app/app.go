package app

import (
	"github.com/org39/webapp-tutorial-backend/usecase/auth"
	"github.com/org39/webapp-tutorial-backend/usecase/user"

	"github.com/facebookgo/inject"
	"github.com/org39/webapp-tutorial-backend/pkg/db"
	"github.com/org39/webapp-tutorial-backend/pkg/log"
)

var DepencencyInjector inject.Graph

type App struct {
	// infra
	Config     *Config     `inject:""`
	RootLogger *log.Logger `inject:""`
	DB         *db.DB      `inject:""`

	// application usecase
	AuthUsecase auth.Usecase `inject:""`
	UserUsecase user.Usecase `inject:""`
}

func New() (*App, error) {
	if err := newInfra(); err != nil {
		return nil, err
	}

	if err := newUserUsecase(); err != nil {
		return nil, err
	}

	if err := newAuthUsecase(); err != nil {
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
