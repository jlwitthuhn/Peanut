// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"context"
	"peanut/internal/logger"
	"time"
)

type ForumPostListingViewRow struct {
	PostId        string
	ThreadId      string
	AuthorName    string
	PostMessage   string
	PostTimestamp time.Time
}

type ForumPostListingDvao interface {
	CreateDBObjects(ctx context.Context) error
}

func NewForumPostListingDvao() ForumPostListingDvao {
	return &forumPostListingDvaoImpl{}
}

type forumPostListingDvaoImpl struct{}

var sqlCreateViewForumPostListing = `
	CREATE VIEW
		view_forum_post_listing
	AS
	(
		SELECT
			fp.id AS post_id,
			fp.thread_id AS thread_id,
			users.display_name AS author_name,
			fp.message AS post_message,
			fp._created AS post_timestamp
		FROM
			forum_posts fp
			INNER JOIN users on fp.author_id = users.id
	)
`

func (*forumPostListingDvaoImpl) CreateDBObjects(ctx context.Context) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlCreateViewForumPostListing)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumPostListingDvao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}
