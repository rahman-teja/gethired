package usecase

import (
	"context"

	"github.com/rahman-teja/gethired/internal/entity"
	"github.com/rahman-teja/gethired/internal/model"
	"github.com/rahman-teja/gethired/internal/repository"
)

type ActivityUsecaseProps struct {
	Repository repository.ActivityRepository
}

type ActivityUsecaseCommand interface {
	Create(ctx context.Context, act model.Activity) (entity.Activity, error)
	Update(ctx context.Context, id int64, act model.Activity) (entity.Activity, error)
	Delete(ctx context.Context, id int64) error
}

type ActivityUsecaseQuery interface {
	GetOne(ctx context.Context, id int64) (entity.Activity, error)
	Get(ctx context.Context) ([]entity.Activity, error)
}
