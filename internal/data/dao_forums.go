// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"context"
	"peanut/internal/logger"
	"time"
)

type ForumRow struct {
	Id         string
	SectionId  string
	Name       string
	Ordering   float32
	Visibility string
	Created    time.Time
	Updated    time.Time
}

type ForumsDao interface {
	CreateDBObjects(ctx context.Context) error
	InsertRow(ctx context.Context, sectionId string, name string, ordering float32, visibility string) (string, error)
	SelectRowAll(ctx context.Context) ([]ForumRow, error)
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

var sqlInsertForumsRow = "INSERT INTO forums(section_id, name, ordering, visibility) VALUES ($1, $2, $3, $4) RETURNING id"

func (*forumsDaoImpl) InsertRow(ctx context.Context, sectionId string, name string, ordering float32, visibility string) (string, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	row := sqlh.QueryRow(sqlInsertForumsRow, sectionId, name, ordering, visibility)
	newId := ""
	err := row.Scan(&newId)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumsDao/InsertRow query: ", err)
		return "", err
	}
	return newId, nil
}

var sqlSelectForumsRowAll = "SELECT id, section_id, name, ordering, visibility, _created, _updated FROM forums ORDER BY ordering"

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
		err = rows.Scan(&thisRow.Id, &thisRow.SectionId, &thisRow.Name, &thisRow.Ordering, &thisRow.Visibility, &thisRow.Created, &thisRow.Updated)
		if err != nil {
			logger.Error(ctx, "Got database error on ForumsDao/SelectRowAll query: ", err)
			return nil, err
		}
		result = append(result, thisRow)
	}
	return result, nil
}
