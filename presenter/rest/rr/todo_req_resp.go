package rr

import (
	"time"

	"github.com/org39/webapp-tutorial-backend/entity"

	"github.com/labstack/echo/v4"
)

func (f *Factory) NewTodoCreatRequest(c echo.Context) (*TodoCreatRequest, error) {
	req := &TodoCreatRequest{}
	err := c.Bind(req)
	return req, err
}

func (f *Factory) NewTodoResponse(todo *entity.Todo) *TodoResponse {
	return &TodoResponse{
		ID:        todo.ID,
		Content:   todo.Content,
		Completed: todo.Completed,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
		Deleted:   todo.Deleted,
	}
}

func (f *Factory) NewTodoUpdateRequest(c echo.Context) (*TodoUpdateRequest, error) {
	req := &TodoUpdateRequest{}
	err := c.Bind(req)
	return req, err
}

func (f *Factory) NewTodosResponse(todos []*entity.Todo) []*TodoResponse {
	resp := make([]*TodoResponse, len(todos))
	for i, todo := range todos {
		resp[i] = &TodoResponse{
			ID:        todo.ID,
			Content:   todo.Content,
			Completed: todo.Completed,
			CreatedAt: todo.CreatedAt,
			UpdatedAt: todo.UpdatedAt,
			Deleted:   todo.Deleted,
		}
	}
	return resp
}

// ------------------------------------------------------------------

type TodoCreatRequest struct {
	Content string `json:"content"`
}

type TodoResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   bool      `json:"deleted"`
}

type TodoUpdateRequest struct {
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
	Deleted   bool   `json:"deleted"`
}
