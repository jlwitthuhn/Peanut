// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"database/sql"
	"peanut/internal/data/datasource"
	"peanut/internal/logger"
)

type SessionRow struct {
	Id     string
	UserId string
}

type SessionDao interface {
	CreateDBObjects(tx *sql.Tx) error
	CountValidDedupeByUser(tx *sql.Tx) (int64, error)
	DeleteRowById(tx *sql.Tx, sessionId string) error
	InsertRow(tx *sql.Tx, sessionId string, userId string) error
	SelectValidRowBySessionId(tx *sql.Tx, sessionId string) (*SessionRow, error)
}

func NewSessionDao() SessionDao {
	return &sessionDaoImpl{}
}

type sessionDaoImpl struct{}

var sqlCreateTableSessions = `
	CREATE TABLE sessions (
		id VARCHAR(255) PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		valid_until TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() + INTERVAL '30 minutes'),
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

var sqlCountValidSessionsDedupeByUser = "SELECT COUNT(DISTINCT user_id) FROM sessions WHERE valid_until >= NOW();"

func (*sessionDaoImpl) CountValidDedupeByUser(tx *sql.Tx) (int64, error) {
	sqlh := selectExecutor(datasource.PostgresHandle(), tx)

	var count int64
	row := sqlh.QueryRow(sqlCountValidSessionsDedupeByUser)
	err := row.Scan(&count)
	if err != nil {
		logger.Error(nil, "Got error on SessionDao/CountValidDedupeByUser query: ", err)
		return 0, err
	}
	return count, nil
}

var sqlDeleteSessionsRowById = "DELETE FROM sessions WHERE id = $1"

func (*sessionDaoImpl) DeleteRowById(tx *sql.Tx, sessionId string) error {
	sqlh := selectExecutor(datasource.PostgresHandle(), tx)

	_, err := sqlh.Exec(sqlDeleteSessionsRowById, sessionId)
	if err != nil {
		logger.Error(nil, "Got error on SessionDao/DeleteRowById query: ", err)
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

var sqlSelectSessionsRowById = "SELECT id, user_id FROM sessions WHERE id = $1 AND valid_until >= NOW()"

func (*sessionDaoImpl) SelectValidRowBySessionId(tx *sql.Tx, sessionId string) (*SessionRow, error) {
	sqlh := selectExecutor(datasource.PostgresHandle(), tx)

	result := &SessionRow{}
	row := sqlh.QueryRow(sqlSelectSessionsRowById, sessionId)
	err := row.Scan(&result.Id, &result.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			logger.Error(nil, "Got error on SessionDao/SelectRowBySessionId query: ", err)
			return nil, err
		}
	}
	return result, nil
}
