// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"context"
	"peanut/internal/logger"
)

type SystemLogDao interface {
	CreateDBObjects(ctx context.Context) error
}

func NewSystemLogDao() SystemLogDao {
	return &systemLogDaoImpl{}
}

type systemLogDaoImpl struct{}

var sqlCreateTableSystemLog = `
	CREATE TABLE system_log (
		id UUID PRIMARY KEY DEFAULT uuidv7(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
		message TEXT NOT NULL,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE TRIGGER
		system_log_trigger_created_updated_before_insert
	BEFORE INSERT ON
		system_log
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		system_log_trigger_created_updated_before_update
	BEFORE UPDATE ON
		system_log
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func (*systemLogDaoImpl) CreateDBObjects(ctx context.Context) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlCreateTableSystemLog)
	if err != nil {
		logger.Error(ctx, "Got database error on SystemLogDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}
