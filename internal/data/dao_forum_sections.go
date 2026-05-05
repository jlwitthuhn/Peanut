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

type ForumSectionRow struct {
	Id         string
	Name       string
	Ordering   float32
	Visibility string
	Created    time.Time
	Updated    time.Time
}

type ForumSectionsDao interface {
	CreateDBObjects(ctx context.Context) error
	InsertRow(ctx context.Context, name string, ordering float32) (string, error)
	SelectRowById(ctx context.Context, id string) (*ForumSectionRow, error)
	SelectRowAll(ctx context.Context) ([]ForumSectionRow, error)
	UpdateUserConfigById(ctx context.Context, id string, name string, ordering float32, visibility string) error
	UpdateVisibilityById(ctx context.Context, id string, visibility string) error
}

func NewForumSectionsDao() ForumSectionsDao {
	return &forumSectionsDaoImpl{}
}

type forumSectionsDaoImpl struct{}

var sqlCreateTableForumSections = `
	CREATE TABLE forum_sections (
		id UUID PRIMARY KEY DEFAULT uuidv7(),
		name VARCHAR(100) UNIQUE NOT NULL,
		ordering REAL NOT NULL,
		visibility visibility_enum NOT NULL DEFAULT 'Public',
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE TRIGGER
		forum_sections_trigger_created_updated_before_insert
	BEFORE INSERT ON
		forum_sections
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		forum_sections_trigger_created_updated_before_update
	BEFORE UPDATE ON
		forum_sections
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func (*forumSectionsDaoImpl) CreateDBObjects(ctx context.Context) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlCreateTableForumSections)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumSectionsDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlInsertForumSectionsRow = "INSERT INTO forum_sections(name, ordering) VALUES ($1, $2) RETURNING id"

func (*forumSectionsDaoImpl) InsertRow(ctx context.Context, name string, ordering float32) (string, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	row := sqlh.QueryRow(sqlInsertForumSectionsRow, name, ordering)
	newId := ""
	err := row.Scan(&newId)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumSectionsDao/InsertRow query: ", err)
		return "", err
	}
	return newId, nil
}

var sqlUpdateForumSectionsUserConfigById = "UPDATE forum_sections SET name = $1, ordering = $2, visibility = $3 WHERE id = $4"

func (*forumSectionsDaoImpl) UpdateUserConfigById(ctx context.Context, id string, name string, ordering float32, visibility string) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlUpdateForumSectionsUserConfigById, name, ordering, visibility, id)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumSectionsDao/UpdateUserConfigById query: ", err)
		return err
	}
	return nil
}

var sqlUpdateForumSectionsVisibilityById = "UPDATE forum_sections SET visibility = $1 WHERE id = $2"

func (*forumSectionsDaoImpl) UpdateVisibilityById(ctx context.Context, id string, visibility string) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlUpdateForumSectionsVisibilityById, visibility, id)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumSectionsDao/UpdateVisibilityById query: ", err)
		return err
	}
	return nil
}

var sqlSelectForumSectionsRowById = "SELECT id, name, ordering, visibility, _created, _updated FROM forum_sections WHERE id = $1"

func (*forumSectionsDaoImpl) SelectRowById(ctx context.Context, id string) (*ForumSectionRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	result := &ForumSectionRow{}
	row := sqlh.QueryRow(sqlSelectForumSectionsRowById, id)
	err := row.Scan(&result.Id, &result.Name, &result.Ordering, &result.Visibility, &result.Created, &result.Updated)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		var pqErr *pq.Error
		ok := errors.As(err, &pqErr)
		if ok {
			if pqErr.Code == "22P02" {
				// INVALID_TEXT_REPRESENTATION
				// This means the input string is not a UUID, so no match exists
				return nil, nil
			}
		}
		logger.Error(ctx, "Got database error on ForumSectionsDao/SelectRowById query: ", err)
		return nil, err
	}
	return result, nil
}

var sqlSelectForumSectionsRowAll = "SELECT id, name, ordering, visibility, _created, _updated FROM forum_sections ORDER BY ordering"

func (*forumSectionsDaoImpl) SelectRowAll(ctx context.Context) ([]ForumSectionRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	rows, err := sqlh.Query(sqlSelectForumSectionsRowAll)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumSectionsDao/SelectRowAll query: ", err)
		return nil, err
	}
	defer rows.Close()

	var result []ForumSectionRow
	for rows.Next() {
		thisRow := ForumSectionRow{}
		err = rows.Scan(&thisRow.Id, &thisRow.Name, &thisRow.Ordering, &thisRow.Visibility, &thisRow.Created, &thisRow.Updated)
		if err != nil {
			logger.Error(ctx, "Got database error on ForumSectionsDao/SelectRowAll query: ", err)
			return nil, err
		}
		result = append(result, thisRow)
	}
	return result, nil
}
