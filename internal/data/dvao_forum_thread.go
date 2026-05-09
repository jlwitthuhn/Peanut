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
	AuthorId         string
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
	SelectForumThreadSummaryViewRowByForumIdPublic(ctx context.Context, forumId string) ([]ForumThreadSummaryViewRow, error)
}

func NewForumThreadSummaryDvao() ForumThreadSummaryDvao {
	return &forumThreadSummaryDvaoImpl{}
}

type forumThreadSummaryDvaoImpl struct{}

var sqlCreateViewForumThreadSummary = `
	CREATE VIEW view_forum_thread_summary AS
	(
		WITH
		q_post_stats_by_thread_id AS (
			SELECT
				thread_id, COUNT(*) AS post_count, MAX(_created) AS last_post_date
			FROM
				forum_posts
			GROUP BY
				thread_id
		)

		SELECT
			ft.id AS thread_id,
			ft.forum_id AS forum_id,
			ft.author_id AS author_id,
			ft.visibility as thread_visibility,
			ft.title AS thread_name,
			us.display_name AS author_name,
			(ps.post_count - 1) AS reply_count,
			ft._created AS thread_created,
			ps.last_post_date AS last_post_date,
			NULL::INTEGER AS unread_post_count
		FROM forum_threads ft
			INNER JOIN users us ON ft.author_id = us.id
			INNER JOIN q_post_stats_by_thread_id ps ON ft.id = ps.thread_id
	)
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

var sqlSelectForumThreadSummaryViewRowByForumIdPublic = `
	SELECT
		thread_id, forum_id, author_id, thread_visibility, thread_name, author_name, reply_count, thread_created, last_post_date, unread_post_count
	FROM
		view_forum_thread_summary
	WHERE
	    forum_id = $1 AND thread_visibility = 'Public'
	ORDER BY
	    last_post_date DESC
`

func (*forumThreadSummaryDvaoImpl) SelectForumThreadSummaryViewRowByForumIdPublic(ctx context.Context, forumId string) ([]ForumThreadSummaryViewRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	rows, err := sqlh.Query(sqlSelectForumThreadSummaryViewRowByForumIdPublic, forumId)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumThreadSummaryDvao/SelectForumThreadSummaryViewRowByForumIdPublic query: ", err)
		return nil, err
	}
	defer rows.Close()

	var result []ForumThreadSummaryViewRow
	for rows.Next() {
		thisRow := ForumThreadSummaryViewRow{}
		err = rows.Scan(&thisRow.ThreadId, &thisRow.ForumId, &thisRow.AuthorId, &thisRow.ThreadVisibility, &thisRow.ThreadName, &thisRow.AuthorName, &thisRow.ReplyCount, &thisRow.ThreadCreated, &thisRow.LastPostDate, &thisRow.UnreadPostCount)
		if err != nil {
			logger.Error(ctx, "Got database error on ForumThreadSummaryDvao/SelectForumThreadSummaryViewRowByForumIdPublic scan: ", err)
			return nil, err
		}
		result = append(result, thisRow)
	}
	return result, nil
}
