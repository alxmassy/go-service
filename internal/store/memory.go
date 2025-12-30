package store

import (
	"sync"
	"context"
)

// type Task struct {
// 	ID    int `json:"id"`
// 	Title string `json:"title"`
// 	Done  bool   `json:"done"`
// }

type MemoryTaskStore struct {
	mu   sync.Mutex
	tasks []Task
	nextID int
}

func NewMemoryTaskStore() *MemoryTaskStore {
	return &MemoryTaskStore{
		tasks: []Task{},
		nextID: 1,
	}
}

func (s *MemoryTaskStore) List(ctx context.Context) ([]Task, error) {
	select {
	case <-ctx.Done():
		return []Task{}, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]Task{}, s.tasks...), nil
}

func (s *MemoryTaskStore) Create(ctx context.Context, title string) (Task, error) {
	select {
	case <-ctx.Done():
		return Task{}, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}

	s.nextID++
	s.tasks = append(s.tasks, task)

	return task, nil
}