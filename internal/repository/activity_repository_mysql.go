package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rahman-teja/gethired/internal/entity"
	"github.com/rahman-teja/gethired/pkg/dbhelper"
	"github.com/rahman-teja/gethired/pkg/sqlcommand"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type ActivityRepositoryMysql struct {
	db    *sql.DB
	table string
}

func NewActivityRepositoryMysql(db *sql.DB) *ActivityRepositoryMysql {
	return &ActivityRepositoryMysql{
		db:    db,
		table: "activities",
	}
}

func (q *ActivityRepositoryMysql) queryRow(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (datas []entity.Activity, err error) {
	var rows *sql.Rows

	rows, err = cmd.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ActivityRepositoryMysql.queryRow",
				"src":   "cmd.QueryContext",
				"query": query,
				"args":  args,
			}).
			Error(query, err)
		return
	}

	defer func() {
		if err := rows.Close(); err != nil {
			logrus.
				WithFields(logrus.Fields{
					"at":    "ActivityRepositoryMysql.queryRow",
					"src":   "rows.Close",
					"query": query,
					"args":  args,
				}).
				Error(err)
		}
	}()

	datas = []entity.Activity{}
	for rows.Next() {
		var data entity.Activity

		var ca, ua, da int64

		err = rows.Scan(
			&data.ID,
			&data.Email,
			&data.Title,
			&ca,
			&ua,
			&da,
		)
		if err != nil {
			logrus.
				WithFields(logrus.Fields{
					"at":  "ActivityRepositoryMysql.queryRow",
					"src": "rows.Scan",
				}).
				Error(err)
			return
		}

		data.CreatedAt = dbhelper.ToSqlFormatFromEpoch(ca)
		data.UpdatedAt = dbhelper.ToSqlFormatFromEpoch(ua)

		datas = append(datas, data)
	}

	return
}

func (a *ActivityRepositoryMysql) GetOne(ctx context.Context, id int64, tx *sql.Tx) (entity.Activity, error) {
	var cmd sqlcommand.Command = a.db
	if tx != nil {
		cmd = tx
	}

	query := fmt.Sprintf("select activity_group_id, email, title, created_at, updated_at, deleted_at from %s where activity_group_id = ?", a.table)

	acts, err := a.queryRow(ctx, cmd, query, id)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ActivityRepositoryMysql.GetOne",
				"src": "q.queryRow",
				"id":  id,
			}).
			Error(err)

		return entity.Activity{}, rapperror.FromMysqlError(
			err,
			err.Error(),
			"activity-repository",
			nil,
		)
	}

	if len(acts) == 0 {
		return entity.Activity{}, rapperror.ErrNotFound(
			"",
			fmt.Sprintf("Activity with ID %d Not Found", id),
			"activity-repository",
			nil,
		)
	}

	return acts[0], nil
}

func (a *ActivityRepositoryMysql) Get(ctx context.Context, tx *sql.Tx) ([]entity.Activity, error) {
	var cmd sqlcommand.Command = a.db
	if tx != nil {
		cmd = tx
	}

	query := fmt.Sprintf("select activity_group_id, email, title, created_at, updated_at, deleted_at from %s", a.table)

	acts, err := a.queryRow(ctx, cmd, query)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ActivityRepositoryMysql.Get",
				"src": "q.queryRow",
			}).
			Error(err)

		return nil, rapperror.FromMysqlError(
			err,
			err.Error(),
			"activity-repository",
			nil,
		)
	}

	return acts, nil
}

func (q *ActivityRepositoryMysql) exec(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (result sql.Result, err error) {
	var stmt *sql.Stmt
	if stmt, err = cmd.PrepareContext(ctx, query); err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ActivityRepositoryMysql.exec",
				"src":   "cmd.PrepareContext",
				"query": query,
				"args":  args,
			}).
			Error(query, err)

		return nil, err
	}

	defer func() {
		if e := stmt.Close(); e != nil {
			logrus.
				WithFields(logrus.Fields{
					"at":    "ActivityRepositoryMysql.exec",
					"src":   "stmt.Close",
					"query": query,
					"args":  args,
				}).
				Error(err)
		}
	}()

	result, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ActivityRepositoryMysql.exec",
				"src":   "stmt.ExecContext",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return nil, err
	}

	return result, nil
}

func (a *ActivityRepositoryMysql) BeginTx(ctx context.Context) (tx *sql.Tx, err error) {
	return a.db.BeginTx(ctx, nil)
}

func (a *ActivityRepositoryMysql) CommitTx(ctx context.Context, tx *sql.Tx) (err error) {
	return tx.Commit()
}

func (a *ActivityRepositoryMysql) RollbackTx(ctx context.Context, tx *sql.Tx) (err error) {
	return tx.Rollback()
}

func (a *ActivityRepositoryMysql) GetOneForUpdate(ctx context.Context, id int64, tx *sql.Tx) (entity.Activity, error) {
	if tx == nil {
		return entity.Activity{}, rapperror.ErrInternalServerError(
			"",
			"Tx is required",
			"todo-repository",
			nil,
		)
	}

	var cmd sqlcommand.Command = tx

	query := fmt.Sprintf("select activity_group_id, email, title, created_at, updated_at, deleted_at from %s where activity_group_id = ? FOR UPDATE", a.table)

	acts, err := a.queryRow(ctx, cmd, query, id)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ActivityRepositoryMysql.GetOne",
				"src": "q.queryRow",
				"id":  id,
			}).
			Error(err)

		return entity.Activity{}, rapperror.FromMysqlError(
			err,
			err.Error(),
			"activity-repository",
			nil,
		)
	}

	if len(acts) == 0 {
		return entity.Activity{}, rapperror.ErrNotFound(
			"",
			fmt.Sprintf("Activity with ID %d Not Found", id),
			"activity-repository",
			nil,
		)
	}

	return acts[0], nil
}

func (a *ActivityRepositoryMysql) Create(ctx context.Context, assets entity.Activity, tx *sql.Tx) (int64, error) {
	var cmd sqlcommand.Command = a.db
	if tx != nil {
		cmd = tx
	}

	query := fmt.Sprintf("insert into %s set email = ?, title = ?, created_at = ?, updated_at = ?", a.table)

	res, err := a.exec(ctx, cmd, query,
		assets.Email,
		assets.Title,
		assets.CreatedAt.UnixNano()/1000000,
		assets.UpdatedAt.UnixNano()/1000000,
	)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ActivityRepositoryMysql.Create",
				"src": "a.exec",
			}).
			Error(err)

		return 0, rapperror.FromMysqlError(
			err,
			err.Error(),
			"activity-repository",
			nil,
		)
	}

	return res.LastInsertId()
}

func (a *ActivityRepositoryMysql) Update(ctx context.Context, id int64, assets entity.Activity, tx *sql.Tx) error {
	var cmd sqlcommand.Command = a.db
	if tx != nil {
		cmd = tx
	}

	query := fmt.Sprintf("update %s set email = ?, title = ?, updated_at = ? where activity_group_id = ?", a.table)

	_, err := a.exec(ctx, cmd, query,
		assets.Email,
		assets.Title,
		time.Now().UnixNano()/1000000,
		id,
	)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ActivityRepositoryMysql.Update",
				"src": "a.exec",
			}).
			Error(err)

		return rapperror.FromMysqlError(
			err,
			err.Error(),
			"activity-repository",
			nil,
		)
	}

	return nil
}

func (a *ActivityRepositoryMysql) Delete(ctx context.Context, id int64, tx *sql.Tx) error {
	var cmd sqlcommand.Command = a.db
	if tx != nil {
		cmd = tx
	}

	query := fmt.Sprintf("delete from %s where activity_group_id = ?", a.table)

	res, err := a.exec(ctx, cmd, query, id)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ActivityRepositoryMysql.Delete",
				"src": "a.exec",
			}).
			Error(err)

		return rapperror.FromMysqlError(
			err,
			err.Error(),
			"activity-repository",
			nil,
		)
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		return rapperror.ErrNotFound(
			"",
			fmt.Sprintf("Activity with ID %d Not Found", id),
			"activity-repository",
			nil,
		)
	}

	return nil
}
