// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"net/http"
	"peanut/internal/logger"
)

type SessionStringDao interface {
	CreateDBObjects(req *http.Request) error
	UpsertString(req *http.Request, sessionId string, name string, value string) error
}

func NewSessionStringDao() SessionStringDao {
	return &sessionStringDaoImpl{}
}

type sessionStringDaoImpl struct{}

var sqlCreateTableSessionString = `
	CREATE TABLE session_string (
	    session_id VARCHAR(255) REFERENCES sessions(id) ON DELETE CASCADE,
	    name VARCHAR(255) NOT NULL,
	    value VARCHAR(255) NOT NULL,
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
		logger.Error(req, "Got error on SessionStringDao/UpsertString query: ", err)
		return err
	}
	return nil
}
