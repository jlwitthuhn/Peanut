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

type ForumThreadRow struct {
	Id         string
	ForumId    string
	AuthorId   string
	Title      string
	Visibility string
	Pinned     bool
	Locked     bool
	Created    time.Time
	Updated    time.Time
}

type ForumThreadsDao interface {
	CreateDBObjects(ctx context.Context) error
	InsertRow(ctx context.Context, forumId string, authorId string, title string, visibility string) (string, error)
	SelectRowById(ctx context.Context, id string) (*ForumThreadRow, error)
	UpdateModerationById(ctx context.Context, id string, title string, pinned bool, locked bool) error
}

func NewForumThreadsDao() ForumThreadsDao {
	return &forumThreadsDaoImpl{}
}

type forumThreadsDaoImpl struct{}

var sqlCreateTableForumThreads = `
	CREATE TABLE forum_threads (
		id UUID PRIMARY KEY DEFAULT uuidv7(),
		forum_id UUID NOT NULL REFERENCES forums(id) ON DELETE RESTRICT,
		author_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
		title VARCHAR(150) NOT NULL,
		visibility visibility_enum NOT NULL,
		pinned BOOLEAN NOT NULL DEFAULT FALSE,
		locked BOOLEAN NOT NULL DEFAULT FALSE,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE TRIGGER
		forum_threads_trigger_created_updated_before_insert
	BEFORE INSERT ON
		forum_threads
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		forum_threads_trigger_created_updated_before_update
	BEFORE UPDATE ON
		forum_threads
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func (*forumThreadsDaoImpl) CreateDBObjects(ctx context.Context) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlCreateTableForumThreads)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumThreadsDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlInsertForumThreadsRow = `
	INSERT INTO
		forum_threads(forum_id, author_id, title, visibility)
	VALUES
	    ($1, $2, $3, $4)
	RETURNING
	    id
`

func (*forumThreadsDaoImpl) InsertRow(ctx context.Context, forumId string, authorId string, title string, visibility string) (string, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	row := sqlh.QueryRow(sqlInsertForumThreadsRow, forumId, authorId, title, visibility)
	newId := ""
	err := row.Scan(&newId)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumThreadsDao/InsertRow query: ", err)
		return "", err
	}
	return newId, nil
}

var sqlSelectForumThreadsRowById = "SELECT id, forum_id, author_id, title, visibility, pinned, locked, _created, _updated FROM forum_threads WHERE id = $1"

func (*forumThreadsDaoImpl) SelectRowById(ctx context.Context, id string) (*ForumThreadRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	result := &ForumThreadRow{}
	row := sqlh.QueryRow(sqlSelectForumThreadsRowById, id)
	err := row.Scan(&result.Id, &result.ForumId, &result.AuthorId, &result.Title, &result.Visibility, &result.Pinned, &result.Locked, &result.Created, &result.Updated)
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
		logger.Error(ctx, "Got database error on ForumThreadsDao/SelectRowById query: ", err)
		return nil, err
	}
	return result, nil
}

var sqlUpdateForumThreadsModerationById = "UPDATE forum_threads SET title = $1, pinned = $2, locked = $3 WHERE id = $4"

func (*forumThreadsDaoImpl) UpdateModerationById(ctx context.Context, id string, title string, pinned bool, locked bool) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlUpdateForumThreadsModerationById, title, pinned, locked, id)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumThreadsDao/UpdateModerationById query: ", err)
		return err
	}
	return nil
}
