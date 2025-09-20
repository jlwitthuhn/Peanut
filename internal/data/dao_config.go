// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"database/sql"
	"sync"
)

type ConfigDao interface {
	CreateDBObjects(*sql.Tx) error
}

var configDaoInstance ConfigDao
var configDaoInstanceOnce sync.Once

func ConfigDaoInst() ConfigDao {
	configDaoInstanceOnce.Do(func() {
		configDaoInstance = &configDaoImpl{}
	})
	return configDaoInstance
}

type configDaoImpl struct{}

func (*configDaoImpl) CreateDBObjects(tx *sql.Tx) error {
	_, intErr := tx.Exec(sqlCreateTableInt)
	if intErr != nil {
		return intErr
	}
	return nil
}

var sqlCreateTableInt = `
	CREATE TABLE config_int (
		name VARCHAR(255) PRIMARY KEY,
		value BIGINT NOT NULL,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE TRIGGER
		config_int_trigger_created_updated_before_insert
	BEFORE INSERT ON
		config_int
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		config_int_trigger_created_updated_before_update
	BEFORE UPDATE ON
		config_int
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`
