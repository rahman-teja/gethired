package repository

import (
	"context"
	"database/sql"

	"github.com/rahman-teja/gethired/internal/entity"
)

type ToDoQueryRepositoryFilter struct {
	ActivityId *int64
}

type ToDoQueryRepository interface {
	GetOne(ctx context.Context, id int64, tx *sql.Tx) (entity.ToDo, error)
	Get(ctx context.Context, filter ToDoQueryRepositoryFilter, tx *sql.Tx) ([]entity.ToDo, error)
}

type ToDoCommandRepository interface {
	BeginTx(ctx context.Context) (tx *sql.Tx, err error)
	CommitTx(ctx context.Context, tx *sql.Tx) (err error)
	RollbackTx(ctx context.Context, tx *sql.Tx) (err error)
	GetOneForUpdate(ctx context.Context, id int64, tx *sql.Tx) (entity.ToDo, error)
	Create(ctx context.Context, assets entity.ToDo, tx *sql.Tx) (int64, error)
	Update(ctx context.Context, id int64, assets entity.ToDo, tx *sql.Tx) error
	Delete(ctx context.Context, id int64, tx *sql.Tx) error
}

type ToDoRepository interface {
	ToDoQueryRepository
	ToDoCommandRepository
}
