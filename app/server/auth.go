package app

import (
	"github.com/org39/webapp-tutorial-backend/usecase/auth"

	"github.com/facebookgo/inject"
)

func newAuthUsecase() error {
	u, err := auth.NewService()
	if err != nil {
		return err
	}

	err = DepencencyInjector.Provide(
		&inject.Object{Value: u},
	)
	if err != nil {
		return err
	}

	return nil
}
