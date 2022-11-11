package usecase

import (
	"context"
	"time"

	"github.com/rahman-teja/gethired/internal/entity"
	"github.com/rahman-teja/gethired/internal/model"
	"github.com/rahman-teja/gethired/internal/repository"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type ToDoCommandUsecase struct {
	repo repository.ToDoCommandRepository
}

func NewToDoCommandUsecase(props ToDoUsecaseProps) *ToDoCommandUsecase {
	return &ToDoCommandUsecase{
		repo: props.Repository,
	}
}

func (a *ToDoCommandUsecase) Create(ctx context.Context, act model.ToDo) (entity.ToDo, error) {
	tmNow := time.Now().UTC()

	insertedToDo := entity.ToDo{
		ActivityGroupID: *act.ActivityGroupID,
		Title:           *act.Title,
		IsActive:        true,
		Priority:        "very-high",
		CreatedAt:       tmNow,
		UpdatedAt:       tmNow,
		DeletedAt:       nil,
	}

	if act.IsActive != nil {
		insertedToDo.IsActive = *act.IsActive
	}

	if act.Priority != nil {
		insertedToDo.Priority = *act.Priority
	}

	insertId, err := a.repo.Create(ctx, insertedToDo, nil)
	if err != nil {
		return entity.ToDo{}, err
	}

	insertedToDo.ID = insertId

	return insertedToDo, nil
}

func (a *ToDoCommandUsecase) Update(ctx context.Context, id int64, act model.ToDo) (entity.ToDo, error) {
	tx, err := a.repo.BeginTx(ctx)
	if err != nil {
		return entity.ToDo{}, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when beginx tx",
			"",
			nil,
		)
	}

	insertedToDo, err := a.repo.GetOneForUpdate(ctx, id, tx)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ToDoCommandUsecase.Update",
				"src": "repo.Update",
			}).
			Error(err)

		tx.Rollback()

		return entity.ToDo{}, err
	}

	if act.ActivityGroupID != nil {
		insertedToDo.ActivityGroupID = *act.ActivityGroupID
	}

	if act.Title != nil {
		insertedToDo.Title = *act.Title
	}

	if act.IsActive != nil {
		isAct := *act.IsActive
		insertedToDo.IsActive = isAct
	}

	if act.Priority != nil {
		insertedToDo.Priority = *act.Priority
	}

	insertedToDo.UpdatedAt = time.Now().UTC()

	err = a.repo.Update(ctx, id, insertedToDo, tx)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ToDoCommandUsecase.Update",
				"src": "repo.Update",
			}).
			Error(err)

		tx.Rollback()

		return entity.ToDo{}, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()

		return entity.ToDo{}, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"",
			nil,
		)
	}

	insertedToDo.ID = id

	return insertedToDo, nil
}

func (a *ToDoCommandUsecase) Delete(ctx context.Context, id int64) error {
	return a.repo.Delete(ctx, id, nil)
}
