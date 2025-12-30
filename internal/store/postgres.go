package store

import (
	"context"
	"database/sql"
)

type PostgresTaskStore struct {
	db *sql.DB
}

func NewPostgresTaskStore(db *sql.DB) *PostgresTaskStore {
	return &PostgresTaskStore{db: db}
}

func (s *PostgresTaskStore) List(ctx context.Context) ([]Task, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, title, done FROM tasks ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, rows.Err()
}

func (s *PostgresTaskStore) Create(ctx context.Context, title string) (Task, error) {
	var t Task
	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO tasks (title) VALUES ($1)
		 RETURNING id, title, done`,
		title,
	).Scan(&t.ID, &t.Title, &t.Done)

	return t, err
}