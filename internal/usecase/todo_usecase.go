package usecase

import (
	"context"

	"github.com/rahman-teja/gethired/internal/entity"
	"github.com/rahman-teja/gethired/internal/model"
	"github.com/rahman-teja/gethired/internal/repository"
)

type ToDoUsecaseProps struct {
	Repository repository.ToDoRepository
}

type ToDoUsecaseCommand interface {
	Create(ctx context.Context, act model.ToDo) (entity.ToDo, error)
	Update(ctx context.Context, id int64, act model.ToDo) (entity.ToDo, error)
	Delete(ctx context.Context, id int64) error
}

type ToDoUsecaseQuery interface {
	GetOne(ctx context.Context, id int64) (entity.ToDo, error)
	Get(ctx context.Context, filter repository.ToDoQueryRepositoryFilter) ([]entity.ToDo, error)
}
