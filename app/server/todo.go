package app

import (
	"github.com/org39/webapp-tutorial-backend/repo"
	"github.com/org39/webapp-tutorial-backend/usecase/todo"

	"github.com/facebookgo/inject"
)

func newTodoUsecase() error {
	r, err := repo.NewTodoRepository()
	if err != nil {
		return err
	}

	u, err := todo.NewService()
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
