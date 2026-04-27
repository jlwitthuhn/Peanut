// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"net/http"
	"peanut/internal/logger"
	"time"
)

type ForumSectionRow struct {
	Id         string
	Name       string
	Ordering   int
	Visibility string
	Created    time.Time
	Updated    time.Time
}

type ForumSectionsDao interface {
	CreateDBObjects(req *http.Request) error
	InsertRow(req *http.Request, name string, ordering int) error
	SelectRowAll(req *http.Request) ([]ForumSectionRow, error)
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

func (*forumSectionsDaoImpl) CreateDBObjects(req *http.Request) error {
	sqlh := getSqlExecutorFromRequest(req)
	_, err := sqlh.Exec(sqlCreateTableForumSections)
	if err != nil {
		logger.Error(nil, "Got error on ForumSectionsDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlInsertForumSectionsRow = "INSERT INTO forum_sections(name, ordering) VALUES ($1, $2)"

func (*forumSectionsDaoImpl) InsertRow(req *http.Request, name string, ordering int) error {
	sqlh := getSqlExecutorFromRequest(req)
	_, err := sqlh.Exec(sqlInsertForumSectionsRow, name, ordering)
	if err != nil {
		logger.Error(nil, "Got error on ForumSectionsDao/InsertRow query: ", err)
		return err
	}
	return nil
}

var sqlSelectForumSectionsRowAll = "SELECT id, name, ordering, visibility, _created, _updated FROM forum_sections ORDER BY ordering"

func (*forumSectionsDaoImpl) SelectRowAll(req *http.Request) ([]ForumSectionRow, error) {
	sqlh := getSqlExecutorFromRequest(req)
	rows, err := sqlh.Query(sqlSelectForumSectionsRowAll)
	if err != nil {
		logger.Error(nil, "Got error on ForumSectionsDao/SelectRowAll query: ", err)
		return nil, err
	}
	defer rows.Close()

	var result []ForumSectionRow
	for rows.Next() {
		thisRow := ForumSectionRow{}
		err = rows.Scan(&thisRow.Id, &thisRow.Name, &thisRow.Ordering, &thisRow.Visibility, &thisRow.Created, &thisRow.Updated)
		if err != nil {
			return nil, err
		}
		result = append(result, thisRow)
	}
	return result, nil
}
