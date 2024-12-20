package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// precompile insert query
	insert_stmt, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, err
	}

	// execute insert
	result, err := insert_stmt.ExecContext(ctx, subject, description)
	if err != nil {
		return nil, err
	}

	// get TODO ID last inserted
	todo_id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// get TODO from DB
	todo_row := s.db.QueryRowContext(ctx, confirm, todo_id)

	var created_at, updated_at time.Time

	// 値の取り出し
	err = todo_row.Scan(&subject, &description, &created_at, &updated_at)

	if err != nil {
		return nil, err
	}

	return &model.TODO{ID: todo_id, Subject: subject, Description: description, CreatedAt: created_at, UpdatedAt: updated_at}, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	var (
		res  = []*model.TODO{}
		rows *sql.Rows

		err error
	)

	if prevID == 0 {
		rows, err = s.db.QueryContext(ctx, read, int(size))

	} else {
		rows, err = s.db.QueryContext(ctx, readWithID, int(prevID), int(size))
	}

	if err != nil {
		return nil, err
	}

	// それぞれのRowに対して処理
	for rows.Next() {
		row := model.TODO{}
		err := rows.Scan(&row.ID, &row.Subject, &row.Description, &row.CreatedAt, &row.UpdatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, &row)
	}
	return res, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	update_stmt, err := s.db.PrepareContext(ctx, update)

	if err != nil {
		return nil, err
	}

	// 更新実行
	result, err := update_stmt.ExecContext(ctx, subject, description, id)

	if err != nil {
		return nil, err
	}

	updated_row, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if updated_row == 0 {
		return nil, &model.ErrNotFound{}
	}

	todo_id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	updated_todo := s.db.QueryRowContext(ctx, confirm, id)
	var created_at, updated_at time.Time

	// 値の取り出し
	err = updated_todo.Scan(&subject, &description, &created_at, &updated_at)
	return &model.TODO{ID: todo_id, Subject: subject, Description: description, CreatedAt: created_at, UpdatedAt: updated_at}, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	symbol := strings.Repeat(", ?", len(ids)-1)

	query := fmt.Sprintf(deleteFmt, symbol)

	log.Printf(query)
	stmt, err := s.db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	var ids_arr []interface{}

	for _, id := range ids {
		ids_arr = append(ids_arr, id)
	}

	if len(ids_arr) == 0 {
		return nil
	}

	result, err := stmt.ExecContext(ctx, ids_arr...)
	if err != nil {
		return err
	}

	affected_row, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected_row == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}
