// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"net/http"
	"peanut/internal/logger"
	"time"
)

type SessionStringRow struct {
	SessionId string
	Name      string
	Value     string
	Created   time.Time
	Updated   time.Time
}

type SessionStringDao interface {
	CreateDBObjects(req *http.Request) error
	SelectRow(req *http.Request, sessionId string, name string) (*SessionStringRow, error)
	UpsertString(req *http.Request, sessionId string, name string, value string) error
}

func NewSessionStringDao() SessionStringDao {
	return &sessionStringDaoImpl{}
}

type sessionStringDaoImpl struct{}

var sqlCreateTableSessionString = `
	CREATE TABLE session_string (
	    session_id VARCHAR REFERENCES sessions(id) ON DELETE CASCADE,
	    name VARCHAR(100) NOT NULL,
	    value VARCHAR(2000) NOT NULL,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL,
		PRIMARY KEY (session_id, name)
	);

	CREATE TRIGGER
		session_string_trigger_created_updated_before_insert
	BEFORE INSERT ON
		session_string
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		session_string_trigger_created_updated_before_update
	BEFORE UPDATE ON
		session_string
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func (*sessionStringDaoImpl) CreateDBObjects(req *http.Request) error {
	sqlh := getSqlExecutorFromRequest(req)
	_, err := sqlh.Exec(sqlCreateTableSessionString)
	if err != nil {
		logger.Error(req, "Got error on SessionStringDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlSelectSessionString = `
	SELECT
		session_id, name, value, _created, _updated
	FROM
		session_string
	WHERE
	    session_id = $1 AND name = $2;
`

func (*sessionStringDaoImpl) SelectRow(req *http.Request, sessionId string, name string) (*SessionStringRow, error) {
	sqlh := getSqlExecutorFromRequest(req)
	result := &SessionStringRow{}
	row := sqlh.QueryRow(sqlSelectSessionString, sessionId, name)
	err := row.Scan(&result.SessionId, &result.Name, &result.Value, &result.Created, &result.Updated)
	if err != nil {
		logger.Error(req, "Got error on SessionStringDao/SelectString query:", err)
		return nil, err
	}
	return result, nil
}

var sqlUpsertSessionStringByName = `
	INSERT INTO
		session_string (session_id, name, value)
	VALUES
		($1, $2, $3)
	ON CONFLICT (session_id, name)
		DO UPDATE SET value = EXCLUDED.value;
`

func (*sessionStringDaoImpl) UpsertString(req *http.Request, sessionId string, name string, value string) error {
	sqlh := getSqlExecutorFromRequest(req)
	_, err := sqlh.Exec(sqlUpsertSessionStringByName, sessionId, name, value)
	if err != nil {
		logger.Error(req, "Got error on SessionStringDao/UpsertString query:", err)
		return err
	}
	return nil
}
