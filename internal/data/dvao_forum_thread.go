// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"context"
	"peanut/internal/logger"
	"time"
)

type ForumThreadSummaryViewRow struct {
	ThreadId         string
	ForumId          string
	ThreadVisibility string
	ThreadName       string
	AuthorName       string
	ReplyCount       int
	ThreadCreated    time.Time
	LastPostDate     time.Time
	UnreadPostCount  *int
}

type ForumThreadSummaryDvao interface {
	CreateDBObjects(ctx context.Context) error
}

func NewForumThreadSummaryDvao() ForumThreadSummaryDvao {
	return &forumThreadSummaryDvaoImpl{}
}

type forumThreadSummaryDvaoImpl struct{}

var sqlCreateViewForumThreadSummary = `
	CREATE VIEW view_forum_thread_summary AS
	SELECT
		ft.id AS thread_id,
		ft.forum_id AS forum_id,
		ft.visibility as thread_visibility,
		ft.title AS thread_name,
		'' AS author_name,
		0 AS reply_count,
		ft._created AS thread_created,
		ft._created AS last_post_date,
		NULL::INTEGER AS unread_post_count
	FROM forum_threads ft;
`

func (*forumThreadSummaryDvaoImpl) CreateDBObjects(ctx context.Context) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlCreateViewForumThreadSummary)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumThreadSummaryDvao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}
