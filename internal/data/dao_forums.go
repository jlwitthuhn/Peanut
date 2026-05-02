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
