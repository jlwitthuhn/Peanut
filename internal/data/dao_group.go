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

type GroupDao interface {
	CreateDBObjects(tx *sql.Tx) error
	InsertRow(tx *sql.Tx, name string, desc string, systemOwned bool) error
}

var groupDaoInstance GroupDao
var groupDaoInstanceOnce sync.Once

func GroupDaoInst() GroupDao {
	groupDaoInstanceOnce.Do(func() {
		groupDaoInstance = &groupDaoImpl{}
	})
	return groupDaoInstance
}

type groupDaoImpl struct{}

func (*groupDaoImpl) SelectRowByName(name string, tx *sql.Tx) {
	//TODO implement me
	panic("implement me")
}

var sqlCreateTableGroups = `
	CREATE TABLE groups (
		id BIGSERIAL PRIMARY KEY,
		name VARCHAR(127) UNIQUE NOT NULL,
		description VARCHAR(255) NOT NULL,
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

func (*groupDaoImpl) CreateDBObjects(tx *sql.Tx) error {
	sqlh := selectExecutor(datasource.PostgresHandle(), tx)
	_, err := sqlh.Exec(sqlCreateTableGroups)
	if err != nil {
		logger.Error(nil, "Got error on GroupDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlInsertGroupsRow = "INSERT INTO groups(name, description, system_owned) VALUES ($1, $2, $3)"

func (*groupDaoImpl) InsertRow(tx *sql.Tx, name string, desc string, systemOwned bool) error {
	sqlh := selectExecutor(datasource.PostgresHandle(), tx)
	_, err := sqlh.Exec(sqlInsertGroupsRow, name, desc, systemOwned)
	if err != nil {
		logger.Error(nil, "Got error on GroupDao/InsertRow query: ", err)
		return err
	}
	return nil
}
