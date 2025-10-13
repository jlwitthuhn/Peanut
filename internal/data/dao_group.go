// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"net/http"
	"peanut/internal/logger"
	"time"
)

type GroupRow struct {
	Id          string
	Name        string
	Description string
	SystemOwned bool
	Created     time.Time
	Updated     time.Time
}

type GroupDao interface {
	CreateDBObjects(req *http.Request) error
	InsertRow(req *http.Request, name string, desc string, systemOwned bool) error
	SelectRowAll(req *http.Request) ([]GroupRow, error)
	SelectRowByName(req *http.Request, name string) (*GroupRow, error)
}

func NewGroupDao() GroupDao {
	return &groupDaoImpl{}
}

type groupDaoImpl struct{}

var sqlCreateTableGroups = `
	CREATE TABLE groups (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(100) UNIQUE NOT NULL,
		description VARCHAR(500) NOT NULL,
		system_owned BOOLEAN NOT NULL,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE TRIGGER
		groups_trigger_created_updated_before_insert
	BEFORE INSERT ON
		groups
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		groups_trigger_created_updated_before_update
	BEFORE UPDATE ON
		groups
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func (*groupDaoImpl) CreateDBObjects(req *http.Request) error {
	sqlh := getSqlExecutorFromRequest(req)
	_, err := sqlh.Exec(sqlCreateTableGroups)
	if err != nil {
		logger.Error(nil, "Got error on GroupDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlInsertGroupsRow = "INSERT INTO groups(name, description, system_owned) VALUES ($1, $2, $3)"

func (*groupDaoImpl) InsertRow(req *http.Request, name string, desc string, systemOwned bool) error {
	sqlh := getSqlExecutorFromRequest(req)
	_, err := sqlh.Exec(sqlInsertGroupsRow, name, desc, systemOwned)
	if err != nil {
		logger.Error(nil, "Got error on GroupDao/InsertRow query: ", err)
		return err
	}
	return nil
}

var sqlSelectGroupsRowAll = "SELECT id, name, description, system_owned, _created, _updated FROM groups ORDER BY name"

func (*groupDaoImpl) SelectRowAll(req *http.Request) ([]GroupRow, error) {
	sqlh := getSqlExecutorFromRequest(req)
	rows, err := sqlh.Query(sqlSelectGroupsRowAll)
	if err != nil {
		logger.Error(nil, "Got error on GroupDao/SelectRowAll query: ", err)
		return nil, err
	}
	var result []GroupRow
	for rows.Next() {
		thisRow := GroupRow{}
		err = rows.Scan(&thisRow.Id, &thisRow.Name, &thisRow.Description, &thisRow.SystemOwned, &thisRow.Created, &thisRow.Updated)
		if err != nil {
			return nil, err
		}
		result = append(result, thisRow)
	}
	return result, nil
}

var sqlSelectGroupsRowByName = "SELECT id, name, description, system_owned FROM groups WHERE name = $1"

func (*groupDaoImpl) SelectRowByName(req *http.Request, name string) (*GroupRow, error) {
	sqlh := getSqlExecutorFromRequest(req)
	result := &GroupRow{}
	row := sqlh.QueryRow(sqlSelectGroupsRowByName, name)
	err := row.Scan(&result.Id, &result.Name, &result.Description, &result.SystemOwned)
	if err != nil {
		logger.Error(nil, "Got error on GroupDao/SelectRowByName query: ", err)
		return nil, err
	}
	return result, nil
}
