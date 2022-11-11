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

type ToDoRepositoryMysql struct {
	db    *sql.DB
	table string
}

func NewToDoRepositoryMysql(db *sql.DB) *ToDoRepositoryMysql {
	return &ToDoRepositoryMysql{
		db:    db,
		table: "todos",
	}
}

func (q *ToDoRepositoryMysql) queryRow(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (datas []entity.ToDo, err error) {
	var rows *sql.Rows

	rows, err = cmd.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ToDoRepositoryMysql.queryRow",
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
					"at":    "ToDoRepositoryMysql.queryRow",
					"src":   "rows.Close",
					"query": query,
					"args":  args,
				}).
				Error(err)
		}
	}()

	datas = []entity.ToDo{}
	for rows.Next() {
		var data entity.ToDo

		var ca, ua, da, ia string

		err = rows.Scan(
			&data.ID,
			&data.ActivityGroupID,
			&data.Title,
			&ia,
			&data.Priority,
			&ca,
			&ua,
			&da,
		)
		if err != nil {
			logrus.
				WithFields(logrus.Fields{
					"at":  "ToDoRepositoryMysql.queryRow",
					"src": "rows.Scan",
				}).
				Error(err)
			return
		}

		data.IsActive = ia == "1"
		data.CreatedAt = dbhelper.ToSqlFormat(ca)
		data.UpdatedAt = dbhelper.ToSqlFormat(ua)

		if da != "" {
			tm := dbhelper.ToSqlFormat(da)
			if !tm.IsZero() {
				data.DeletedAt = &tm
			}
		}

		datas = append(datas, data)
	}

	return
}

func (a *ToDoRepositoryMysql) GetOne(ctx context.Context, id int64, tx *sql.Tx) (entity.ToDo, error) {
	var cmd sqlcommand.Command = a.db
	if tx != nil {
		cmd = tx
	}

	query := fmt.Sprintf("select id, activity_group_id, title, is_active, priority, created_at, updated_at, deleted_at from %s where id = ? and is_deleted = ?", a.table)

	acts, err := a.queryRow(ctx, cmd, query, id, 0)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ToDoRepositoryMysql.GetOne",
				"src": "q.queryRow",
				"id":  id,
			}).
			Error(err)

		return entity.ToDo{}, rapperror.FromMysqlError(
			err,
			err.Error(),
			"todo-repository",
			nil,
		)
	}

	if len(acts) == 0 {
		return entity.ToDo{}, rapperror.ErrNotFound(
			"",
			fmt.Sprintf("Todo with ID %d Not Found", id),
			"todo-repository",
			nil,
		)
	}

	return acts[0], nil
}

func (a *ToDoRepositoryMysql) Get(ctx context.Context, filter ToDoQueryRepositoryFilter, tx *sql.Tx) ([]entity.ToDo, error) {
	var cmd sqlcommand.Command = a.db
	if tx != nil {
		cmd = tx
	}

	query := fmt.Sprintf("select id, activity_group_id, title, is_active, priority, created_at, updated_at, deleted_at from %s where is_deleted = ?", a.table)

	args := []interface{}{
		0,
	}

	if filter.ActivityId != nil {
		query += " and activity_group_id = ?"
		args = append(args, *filter.ActivityId)
	}

	acts, err := a.queryRow(ctx, cmd, query, args...)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ToDoRepositoryMysql.Get",
				"src": "q.queryRow",
			}).
			Error(err)

		return nil, rapperror.FromMysqlError(
			err,
			err.Error(),
			"todo-repository",
			nil,
		)
	}

	return acts, nil
}

func (q *ToDoRepositoryMysql) exec(ctx context.Context, cmd sqlcommand.Command, query string, args ...interface{}) (result sql.Result, err error) {
	var stmt *sql.Stmt
	if stmt, err = cmd.PrepareContext(ctx, query); err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ToDoRepositoryMysql.exec",
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
					"at":    "ToDoRepositoryMysql.exec",
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
				"at":    "ToDoRepositoryMysql.exec",
				"src":   "stmt.ExecContext",
				"query": query,
				"args":  args,
			}).
			Error(err)

		return nil, err
	}

	return result, nil
}

func (a *ToDoRepositoryMysql) BeginTx(ctx context.Context) (tx *sql.Tx, err error) {
	return a.db.BeginTx(ctx, nil)
}

func (a *ToDoRepositoryMysql) CommitTx(ctx context.Context, tx *sql.Tx) (err error) {
	return tx.Commit()
}

func (a *ToDoRepositoryMysql) RollbackTx(ctx context.Context, tx *sql.Tx) (err error) {
	return tx.Rollback()
}

func (a *ToDoRepositoryMysql) GetOneForUpdate(ctx context.Context, id int64, tx *sql.Tx) (entity.ToDo, error) {
	if tx == nil {
		return entity.ToDo{}, rapperror.ErrInternalServerError(
			"",
			"Tx is required",
			"todo-repository",
			nil,
		)
	}

	var cmd sqlcommand.Command = tx

	query := fmt.Sprintf("select id, activity_group_id, title, is_active, priority, created_at, updated_at, deleted_at from %s where id = ? and is_deleted = ? FOR UPDATE", a.table)

	acts, err := a.queryRow(ctx, cmd, query, id, 0)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ToDoRepositoryMysql.GetOne",
				"src": "q.queryRow",
				"id":  id,
			}).
			Error(err)

		return entity.ToDo{}, rapperror.FromMysqlError(
			err,
			err.Error(),
			"todo-repository",
			nil,
		)
	}

	if len(acts) == 0 {
		return entity.ToDo{}, rapperror.ErrNotFound(
			"",
			fmt.Sprintf("Todo with ID %d Not Found", id),
			"todo-repository",
			nil,
		)
	}

	return acts[0], nil
}

func (a *ToDoRepositoryMysql) Create(ctx context.Context, assets entity.ToDo, tx *sql.Tx) (int64, error) {
	var cmd sqlcommand.Command = a.db
	if tx != nil {
		cmd = tx
	}

	query := fmt.Sprintf("insert into %s set activity_group_id = ?, title = ?, is_active = ?, priority = ?, created_at = ?, updated_at = ?", a.table)

	res, err := a.exec(ctx, cmd, query, assets.ActivityGroupID, assets.Title, assets.IsActive, assets.Priority, assets.CreatedAt, assets.UpdatedAt)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ToDoRepositoryMysql.Create",
				"src": "a.exec",
			}).
			Error(err)

		return 0, rapperror.FromMysqlError(
			err,
			"",
			"todo-repository",
			nil,
		)
	}

	return res.LastInsertId()
}

func (a *ToDoRepositoryMysql) Update(ctx context.Context, id int64, assets entity.ToDo, tx *sql.Tx) error {
	var cmd sqlcommand.Command = a.db
	if tx != nil {
		cmd = tx
	}

	query := fmt.Sprintf("update %s set activity_group_id = ?, title = ?, is_active = ?, priority = ?, created_at = ?, updated_at = ? where id = ?", a.table)

	_, err := a.exec(ctx, cmd, query, assets.ActivityGroupID, assets.Title, assets.IsActive, assets.Priority, assets.CreatedAt, assets.UpdatedAt, id)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ToDoRepositoryMysql.Update",
				"src": "a.exec",
			}).
			Error(err)

		return rapperror.FromMysqlError(
			err,
			err.Error(),
			"todo-repository",
			nil,
		)
	}

	return nil
}

func (a *ToDoRepositoryMysql) Delete(ctx context.Context, id int64, tx *sql.Tx) error {
	var cmd sqlcommand.Command = a.db
	if tx != nil {
		cmd = tx
	}

	query := fmt.Sprintf("update %s set deleted_at = ?, is_deleted = ? where is_deleted = 0 and id = ?", a.table)

	res, err := a.exec(ctx, cmd, query, time.Now().UTC(), 1, id)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ToDoRepositoryMysql.Delete",
				"src": "a.exec",
			}).
			Error(err)

		return rapperror.FromMysqlError(
			err,
			err.Error(),
			"todo-repository",
			nil,
		)
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		return rapperror.ErrNotFound(
			"",
			fmt.Sprintf("Todo with ID %d Not Found", id),
			"activity-repository",
			nil,
		)
	}

	return nil
}
