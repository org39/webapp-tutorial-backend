package app

import (
	"github.com/org39/webapp-tutorial-backend/repo"
	"github.com/org39/webapp-tutorial-backend/usecase/user"

	"github.com/facebookgo/inject"
)

func newUserUsecase() error {
	r, err := repo.NewUserRepository()
	if err != nil {
		return err
	}

	u, err := user.NewService()
	if err != nil {
		return err
	}

	err = DepencencyInjector.Provide(
		&inject.Object{Value: r},
		&inject.Object{Value: u},
	)
	if err != nil {
		return err
	}

	return nil
}
