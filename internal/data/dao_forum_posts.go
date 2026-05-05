// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"context"
	"peanut/internal/logger"
	"time"
)

type ForumPostRow struct {
	Id       string
	ThreadId string
	AuthorId string
	Message  string
	Created  time.Time
	Updated  time.Time
}

type ForumPostsDao interface {
	CreateDBObjects(ctx context.Context) error
	InsertRow(ctx context.Context, threadId string, authorId string, message string) (string, error)
}

func NewForumPostsDao() ForumPostsDao {
	return &forumPostsDaoImpl{}
}

type forumPostsDaoImpl struct{}

var sqlCreateTableForumPosts = `
	CREATE TABLE forum_posts (
		id UUID PRIMARY KEY DEFAULT uuidv7(),
		thread_id UUID NOT NULL REFERENCES forum_threads(id) ON DELETE RESTRICT,
		author_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
		message TEXT NOT NULL,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE TRIGGER
		forum_posts_trigger_created_updated_before_insert
	BEFORE INSERT ON
		forum_posts
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		forum_posts_trigger_created_updated_before_update
	BEFORE UPDATE ON
		forum_posts
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func (*forumPostsDaoImpl) CreateDBObjects(ctx context.Context) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlCreateTableForumPosts)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumPostsDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlInsertForumPostsRow = `
	INSERT INTO
		forum_posts(thread_id, author_id, message)
	VALUES
	    ($1, $2, $3)
	RETURNING
	    id
`

func (*forumPostsDaoImpl) InsertRow(ctx context.Context, threadId string, authorId string, message string) (string, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	row := sqlh.QueryRow(sqlInsertForumPostsRow, threadId, authorId, message)
	newId := ""
	err := row.Scan(&newId)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumPostsDao/InsertRow query: ", err)
		return "", err
	}
	return newId, nil
}
