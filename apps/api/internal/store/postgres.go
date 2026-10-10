package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mimustafa/research-workflow-platform/apps/api/internal/workflow"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *Postgres {
	return &Postgres{pool: pool}
}

func (s *Postgres) Create(ctx context.Context, item workflow.Workflow) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO workflows (id, brief, status, created_at)
		VALUES ($1, $2, $3, $4)
	`, item.ID, item.Brief, item.Status, item.CreatedAt)
	return err
}

func (s *Postgres) List(ctx context.Context) ([]workflow.Workflow, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, brief, status, COALESCE(result, ''), created_at
		FROM workflows
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []workflow.Workflow
	for rows.Next() {
		var item workflow.Workflow
		if err := rows.Scan(&item.ID, &item.Brief, &item.Status, &item.Result, &item.CreatedAt); err != nil {
			return nil, err
		}
		workflows = append(workflows, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return workflows, nil
}

func (s *Postgres) MarkFailed(ctx context.Context, id string, message string) error {
	command, err := s.pool.Exec(ctx, `
		UPDATE workflows
		SET status = 'failed', result = $2
		WHERE id = $1
	`, id, message)
	if err != nil {
		return err
	}
	if command.RowsAffected() != 1 {
		return fmt.Errorf("workflow %s was not found", id)
	}
	return nil
}
