// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"context"
	"peanut/internal/logger"
)

type ForumHomeViewRow struct {
	SectionId       string
	SectionName     string
	SectionOrdering float32
	ForumId         string
	ForumName       string
	ForumOrdering   float32
}

type ForumHomeDvao interface {
	CreateDBObjects(ctx context.Context) error
	SelectRowAll(ctx context.Context) ([]ForumHomeViewRow, error)
}

func NewForumHomeDvao() ForumHomeDvao {
	return &forumHomeDvaoImpl{}
}

type forumHomeDvaoImpl struct{}

var sqlCreateViewForumHomePublic = `
	CREATE VIEW view_forum_home_public AS
	SELECT
		fs.id AS section_id,
		fs.name AS section_name,
		fs.ordering AS section_ordering,
		f.id AS forum_id,
		f.name AS forum_name,
		f.ordering AS forum_ordering
	FROM forums f
	INNER JOIN forum_sections fs ON f.section_id = fs.id
	WHERE f.visibility = 'Public'
	AND fs.visibility = 'Public'
	ORDER BY fs.ordering ASC, f.ordering ASC;
`

func (*forumHomeDvaoImpl) CreateDBObjects(ctx context.Context) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlCreateViewForumHomePublic)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumHomeDvao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlSelectForumHomeViewRowAll = "SELECT section_id, section_name, section_ordering, forum_id, forum_name, forum_ordering FROM view_forum_home_public"

func (*forumHomeDvaoImpl) SelectRowAll(ctx context.Context) ([]ForumHomeViewRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	rows, err := sqlh.Query(sqlSelectForumHomeViewRowAll)
	if err != nil {
		logger.Error(ctx, "Got database error on ForumHomeDvao/SelectRowAll query: ", err)
		return nil, err
	}
	defer rows.Close()

	var result []ForumHomeViewRow
	for rows.Next() {
		thisRow := ForumHomeViewRow{}
		err = rows.Scan(&thisRow.SectionId, &thisRow.SectionName, &thisRow.SectionOrdering, &thisRow.ForumId, &thisRow.ForumName, &thisRow.ForumOrdering)
		if err != nil {
			logger.Error(ctx, "Got database error on ForumHomeDvao/SelectRowAll query: ", err)
			return nil, err
		}
		result = append(result, thisRow)
	}
	return result, nil
}
