// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"context"
	"database/sql"
	"errors"
	"peanut/internal/logger"
	"time"

	"github.com/lib/pq"
)

type ForumRow struct {
	Id          string
	SectionId   string
	Name        string
	Description string
	Ordering    float32
	Visibility  string
	Created     time.Time
	Updated     time.Time
}

type ForumsDao interface {
	CreateDBObjects(ctx context.Context) error
	InsertRow(ctx context.Context, sectionId string, name string, description string, ordering float32, visibility string) (string, error)
	SelectRowAll(ctx context.Context) ([]ForumRow, error)
	SelectRowById(ctx context.Context, id string) (*ForumRow, error)
	UpdateUserConfigById(ctx context.Context, id string, sectionId string, name string, description string, ordering float32, visibility string) error
	UpdateVisibilityById(ctx context.Context, id string, visibility string) error
}

func NewForumsDao() ForumsDao {
	return &forumsDaoImpl{}
}

type forumsDaoImpl struct{}

var sqlCreateTableForums = `
	CREATE TABLE forums (
		id UUID PRIMARY KEY DEFAULT uuidv7(),
		section_id UUID NOT NULL REFERENCES forum_sections(id) ON DELETE RESTRICT,
		name VARCHAR(150) NOT NULL,
		description VARCHAR(250) NOT NULL,
		ordering REAL NOT NULL,
		visibility visibility_enum NOT NULL DEFAULT 'Public',
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE TRIGGER
		forums_trigger_created_updated_before_insert
	BEFORE INSERT ON
		forums
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		forums_trigger_created_updated_before_update
	BEFORE UPDATE ON
		forums
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func (*forumsDaoImpl) CreateDBObjects(ctx context.Context) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlCreateTableForums)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumsDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlInsertForumsRow = "INSERT INTO forums(section_id, name, description, ordering, visibility) VALUES ($1, $2, $3, $4, $5) RETURNING id"

func (*forumsDaoImpl) InsertRow(ctx context.Context, sectionId string, name string, description string, ordering float32, visibility string) (string, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	row := sqlh.QueryRow(sqlInsertForumsRow, sectionId, name, description, ordering, visibility)
	newId := ""
	err := row.Scan(&newId)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumsDao/InsertRow query: ", err)
		return "", err
	}
	return newId, nil
}

var sqlUpdateForumsVisibilityById = "UPDATE forums SET visibility = $1 WHERE id = $2"

func (*forumsDaoImpl) UpdateVisibilityById(ctx context.Context, id string, visibility string) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlUpdateForumsVisibilityById, visibility, id)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumsDao/UpdateVisibilityById query: ", err)
		return err
	}
	return nil
}

var sqlSelectForumsRowAll = "SELECT id, section_id, name, description, ordering, visibility, _created, _updated FROM forums ORDER BY ordering"

func (*forumsDaoImpl) SelectRowAll(ctx context.Context) ([]ForumRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	rows, err := sqlh.Query(sqlSelectForumsRowAll)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumsDao/SelectRowAll query: ", err)
		return nil, err
	}
	defer rows.Close()

	var result []ForumRow
	for rows.Next() {
		thisRow := ForumRow{}
		err = rows.Scan(&thisRow.Id, &thisRow.SectionId, &thisRow.Name, &thisRow.Description, &thisRow.Ordering, &thisRow.Visibility, &thisRow.Created, &thisRow.Updated)
		if err != nil {
			logger.Error(ctx, "Got database error on ForumsDao/SelectRowAll query: ", err)
			return nil, err
		}
		result = append(result, thisRow)
	}
	return result, nil
}

var sqlSelectForumsRowById = "SELECT id, section_id, name, description, ordering, visibility, _created, _updated FROM forums WHERE id = $1"

func (*forumsDaoImpl) SelectRowById(ctx context.Context, id string) (*ForumRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	result := &ForumRow{}
	row := sqlh.QueryRow(sqlSelectForumsRowById, id)
	err := row.Scan(&result.Id, &result.SectionId, &result.Name, &result.Description, &result.Ordering, &result.Visibility, &result.Created, &result.Updated)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		var pqErr *pq.Error
		ok := errors.As(err, &pqErr)
		if ok {
			if pqErr.Code == "22P02" {
				return nil, nil
			}
		}
		logger.Error(ctx, "Got database error on ForumsDao/SelectRowById query: ", err)
		return nil, err
	}
	return result, nil
}

var sqlUpdateForumsUserConfigById = "UPDATE forums SET section_id = $1, name = $2, description = $3, ordering = $4, visibility = $5 WHERE id = $6"

func (*forumsDaoImpl) UpdateUserConfigById(ctx context.Context, id string, sectionId string, name string, description string, ordering float32, visibility string) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlUpdateForumsUserConfigById, sectionId, name, description, ordering, visibility, id)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumsDao/UpdateUserConfigById query: ", err)
		return err
	}
	return nil
}
