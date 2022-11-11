package usecase

import (
	"context"

	"github.com/rahman-teja/gethired/internal/entity"
	"github.com/rahman-teja/gethired/internal/repository"
)

type ActivityQueryUsecase struct {
	repo repository.ActivityQueryRepository
}

func NewActivityQueryUsecase(props ActivityUsecaseProps) *ActivityQueryUsecase {
	return &ActivityQueryUsecase{
		repo: props.Repository,
	}
}

func (a *ActivityQueryUsecase) GetOne(ctx context.Context, id int64) (entity.Activity, error) {
	return a.repo.GetOne(ctx, id, nil)
}

func (a *ActivityQueryUsecase) Get(ctx context.Context) ([]entity.Activity, error) {
	return a.repo.Get(ctx, nil)
}
