package repository

import (
	"context"
	"database/sql"

	"github.com/rahman-teja/gethired/internal/entity"
)

type ActivityQueryRepository interface {
	GetOne(ctx context.Context, id int64, tx *sql.Tx) (entity.Activity, error)
	Get(ctx context.Context, tx *sql.Tx) ([]entity.Activity, error)
}

type ActivityCommandRepository interface {
	BeginTx(ctx context.Context) (tx *sql.Tx, err error)
	CommitTx(ctx context.Context, tx *sql.Tx) (err error)
	RollbackTx(ctx context.Context, tx *sql.Tx) (err error)
	GetOneForUpdate(ctx context.Context, id int64, tx *sql.Tx) (entity.Activity, error)
	Create(ctx context.Context, assets entity.Activity, tx *sql.Tx) (int64, error)
	Update(ctx context.Context, id int64, assets entity.Activity, tx *sql.Tx) error
	Delete(ctx context.Context, id int64, tx *sql.Tx) error
}

type ActivityRepository interface {
	ActivityQueryRepository
	ActivityCommandRepository
}
