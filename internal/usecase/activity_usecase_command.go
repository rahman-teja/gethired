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

type ActivityCommandUsecase struct {
	repo repository.ActivityCommandRepository
}

func NewActivityCommandUsecase(props ActivityUsecaseProps) *ActivityCommandUsecase {
	return &ActivityCommandUsecase{
		repo: props.Repository,
	}
}

func (a *ActivityCommandUsecase) Create(ctx context.Context, act model.Activity) (entity.Activity, error) {
	tmNow := time.Now().UTC()

	insertedActivity := entity.Activity{
		Title:     *act.Title,
		CreatedAt: tmNow,
		UpdatedAt: tmNow,
		DeletedAt: nil,
	}

	if act.Email != nil {
		insertedActivity.Email = *act.Email
	}

	insertId, err := a.repo.Create(ctx, insertedActivity, nil)
	if err != nil {
		return entity.Activity{}, err
	}

	insertedActivity.ID = insertId

	return insertedActivity, nil
}

func (a *ActivityCommandUsecase) Update(ctx context.Context, id int64, act model.Activity) (entity.Activity, error) {
	tx, err := a.repo.BeginTx(ctx)
	if err != nil {
		return entity.Activity{}, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when beginx tx",
			"",
			nil,
		)
	}

	insertedActivity, err := a.repo.GetOneForUpdate(ctx, id, tx)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ActivityCommandUsecase.Update",
				"src": "repo.Update",
			}).
			Error(err)

		tx.Rollback()

		return entity.Activity{}, err
	}

	if act.Email != nil {
		insertedActivity.Email = *act.Email
	}

	if act.Title != nil {
		insertedActivity.Title = *act.Title
	}

	insertedActivity.UpdatedAt = time.Now().UTC()

	err = a.repo.Update(ctx, id, insertedActivity, tx)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ActivityCommandUsecase.Update",
				"src": "repo.Update",
			}).
			Error(err)

		tx.Rollback()

		return entity.Activity{}, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()

		return entity.Activity{}, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"",
			nil,
		)
	}

	insertedActivity.ID = id

	return insertedActivity, nil
}

func (a *ActivityCommandUsecase) Delete(ctx context.Context, id int64) error {
	return a.repo.Delete(ctx, id, nil)
}
