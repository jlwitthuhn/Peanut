// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"database/sql"
	"peanut/internal/data/datasource"
	"peanut/internal/logger"
	"sync"
)

type SessionDao interface {
	CreateDBObjects(tx *sql.Tx) error
	InsertRow(tx *sql.Tx, sessionId string, userId string) error
}

var sessionDaoInstance SessionDao
var sessionDaoInstanceOnce sync.Once

func SessionDaoInst() SessionDao {
	sessionDaoInstanceOnce.Do(func() {
		sessionDaoInstance = &sessionDaoImpl{}
	})
	return sessionDaoInstance
}

type sessionDaoImpl struct{}

var sqlCreateTableSessions = `
	CREATE TABLE sessions (
		id VARCHAR(63) PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE TRIGGER
		sessions_trigger_created_updated_before_insert
	BEFORE INSERT ON
		sessions
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		sessions_trigger_created_updated_before_update
	BEFORE UPDATE ON
		sessions
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func (*sessionDaoImpl) CreateDBObjects(tx *sql.Tx) error {
	sqlh := selectExecutor(datasource.PostgresHandle(), tx)
	_, err := sqlh.Exec(sqlCreateTableSessions)
	if err != nil {
		logger.Error(nil, "Got error on SessionDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlInsertSessionsRow = "INSERT INTO sessions (id, user_id) VALUES ($1, $2::uuid)"

func (*sessionDaoImpl) InsertRow(tx *sql.Tx, sessionId string, userId string) error {
	sqlh := selectExecutor(datasource.PostgresHandle(), tx)

	_, err := sqlh.Exec(sqlInsertSessionsRow, sessionId, userId)
	if err != nil {
		logger.Error(nil, "Got error on SessionDao/InsertRow query: ", err)
		return err
	}
	return nil
}
