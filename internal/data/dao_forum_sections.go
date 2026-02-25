// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"net/http"
	"peanut/internal/logger"
)

type ForumSectionsDao interface {
	CreateDBObjects(req *http.Request) error
}

func NewForumSectionsDao() ForumSectionsDao {
	return &forumSectionsDaoImpl{}
}

type forumSectionsDaoImpl struct{}

var sqlCreateTableForumSections = `
	CREATE TABLE forum_sections (
		id UUID PRIMARY KEY DEFAULT uuidv7(),
		name VARCHAR(100) UNIQUE NOT NULL,
		ordering INTEGER NOT NULL,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL
	);
`

func (*forumSectionsDaoImpl) CreateDBObjects(req *http.Request) error {
	sqlh := getSqlExecutorFromRequest(req)
	_, err := sqlh.Exec(sqlCreateTableForumSections)
	if err != nil {
		logger.Error(nil, "Got error on ForumSectionsDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}
