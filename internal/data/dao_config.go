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

type ConfigDao interface {
	CreateDBObjects(*sql.Tx) error
	UpsertIntByName(name string, value int64, tx *sql.Tx) error
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

func (*configDaoImpl) CreateDBObjects(tx *sql.Tx) error {
	_, err := tx.Exec(sqlCreateTableInt)
	if err != nil {
		logger.Error(nil, "Got error on ConfigDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlUpsertIntByName = `
	INSERT INTO
		config_int (name, value)
	VALUES
		($1, $2)
	ON CONFLICT (name)
		DO UPDATE SET value = EXCLUDED.value;
`

func (*configDaoImpl) UpsertIntByName(name string, value int64, tx *sql.Tx) error {
	sqlh := selectExecutor(datasource.PostgresHandle(), tx)
	_, err := sqlh.Exec(sqlUpsertIntByName, name, value)
	if err != nil {
		logger.Error(nil, "Got error on ConfigDao/UpsertIntByName query: ", err)
		return err
	}
	return nil
}
