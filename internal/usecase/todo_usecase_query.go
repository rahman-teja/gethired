package usecase

import (
	"context"

	"github.com/rahman-teja/gethired/internal/entity"
	"github.com/rahman-teja/gethired/internal/repository"
)

type ToDoQueryUsecase struct {
	repo repository.ToDoQueryRepository
}

func NewToDoQueryUsecase(props ToDoUsecaseProps) *ToDoQueryUsecase {
	return &ToDoQueryUsecase{
		repo: props.Repository,
	}
}

func (a *ToDoQueryUsecase) GetOne(ctx context.Context, id int64) (entity.ToDo, error) {
	return a.repo.GetOne(ctx, id, nil)
}

func (a *ToDoQueryUsecase) Get(ctx context.Context, filter repository.ToDoQueryRepositoryFilter) ([]entity.ToDo, error) {
	return a.repo.Get(ctx, filter, nil)
}
