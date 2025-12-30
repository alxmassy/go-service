package store

import (
	"context"
)

type Task struct {
	ID int
	Title string
	Done bool
}

type TaskStore interface {
	List(ctx context.Context) ([]Task, error)
	Create(ctx context.Context, title string) (Task, error)
}