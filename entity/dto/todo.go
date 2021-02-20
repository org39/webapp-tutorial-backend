package dto

import (
	"time"
)

type Todo struct {
	ID        string
	UserID    string
	Content   string
	Completed bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   bool
}
