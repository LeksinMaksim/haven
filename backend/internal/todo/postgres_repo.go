package todo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(pool *pgxpool.Pool) Repository {
	return &postgresRepo{pool: pool}
}

func (r *postgresRepo) Create(ctx context.Context, t *Todo) error {
	query := `
		INSERT INTO todos (title, description, done, priority, due_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		t.Title, t.Description, t.Done, t.Priority, t.DueDate,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create todo: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetByID(ctx context.Context, id int64) (*Todo, error) {
	query := `
		SELECT id, title, description, done, priority, due_date, created_at, updated_at
		FROM todos WHERE id = $1`

	t := &Todo{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.Done,
		&t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get todo by id: %w", err)
	}
	return t, nil
}

func (r *postgresRepo) List(ctx context.Context, filter ListFilter) ([]Todo, error) {
	var (
		conditions []string
		args       []any
	)

	if filter.Done != nil {
		args = append(args, *filter.Done)
		conditions = append(conditions, fmt.Sprintf("done = $%d", len(args)))
	}
	if filter.Priority != nil {
		args = append(args, *filter.Priority)
		conditions = append(conditions, fmt.Sprintf("priority = $%d", len(args)))
	}

	query := "SELECT id, title, description, done, priority, due_date, created_at, updated_at FROM todos"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}

	todos, err := pgx.CollectRows(rows, pgx.RowToStructByName[Todo])
	if err != nil {
		return nil, fmt.Errorf("collecting todos: %w", err)
	}

	return todos, nil
}

func (r *postgresRepo) Update(ctx context.Context, t *Todo) error {
	query := `
		UPDATE todos
		SET title=$1, description=$2, done=$3, priority=$4, due_date=$5, update_at=$6
		WHERE id=$7`

	t.UpdatedAt = time.Now()
	res, err := r.pool.Exec(ctx, query,
		t.Title, t.Description, t.Done, t.Priority, t.DueDate, t.UpdatedAt, t.ID,
	)
	if err != nil {
		return fmt.Errorf("update todo: %w", err)
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *postgresRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.pool.Exec(ctx, "DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
