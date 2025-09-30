// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"database/sql"
	"peanut/internal/data/datasource"
	"peanut/internal/logger"
)

type ConfigIntRow struct {
	Name  string
	Value int64
}

type ConfigDao interface {
	CreateDBObjects(*sql.Tx) error
	SelectIntRowByName(tx *sql.Tx, name string) (*ConfigIntRow, error)
	UpsertIntByName(tx *sql.Tx, name string, value int64) error
	UpsertStringByName(tx *sql.Tx, name string, value string) error
}

func NewConfigDao() ConfigDao {
	return &configDaoImpl{}
}

type configDaoImpl struct{}

func (*configDaoImpl) SelectRowByName(name string, tx *sql.Tx) {
	//TODO implement me
	panic("implement me")
}

var sqlCreateTableConfigInt = `
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

var sqlCreateTableConfigString = `
	CREATE TABLE config_string (
		name VARCHAR(255) PRIMARY KEY,
		value VARCHAR(4096) NOT NULL,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE TRIGGER
		config_string_trigger_created_updated_before_insert
	BEFORE INSERT ON
		config_string
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		config_string_trigger_created_updated_before_update
	BEFORE UPDATE ON
		config_string
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func (*configDaoImpl) CreateDBObjects(tx *sql.Tx) error {
	_, err := tx.Exec(sqlCreateTableConfigInt)
	if err != nil {
		logger.Error(nil, "Got error on ConfigDao/CreateDBObjects query: ", err)
		return err
	}
	_, err = tx.Exec(sqlCreateTableConfigString)
	if err != nil {
		logger.Error(nil, "Got error on ConfigDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlSelectConfigIntRowByName = "SELECT name, value FROM config_int WHERE name = $1"

func (*configDaoImpl) SelectIntRowByName(tx *sql.Tx, name string) (*ConfigIntRow, error) {
	sqlh := selectExecutor(datasource.PostgresHandle(), tx)

	result := &ConfigIntRow{}
	row := sqlh.QueryRow(sqlSelectConfigIntRowByName, name)
	err := row.Scan(&result.Name, &result.Value)
	if err != nil {
		return nil, err
	}
	return result, nil
}

var sqlUpsertConfigIntByName = `
	INSERT INTO
		config_int (name, value)
	VALUES
		($1, $2)
	ON CONFLICT (name)
		DO UPDATE SET value = EXCLUDED.value;
`

func (*configDaoImpl) UpsertIntByName(tx *sql.Tx, name string, value int64) error {
	sqlh := selectExecutor(datasource.PostgresHandle(), tx)
	_, err := sqlh.Exec(sqlUpsertConfigIntByName, name, value)
	if err != nil {
		logger.Error(nil, "Got error on ConfigDao/UpsertIntByName query: ", err)
		return err
	}
	return nil
}

var sqlUpsertConfigStringByName = `
	INSERT INTO
		config_string (name, value)
	VALUES
		($1, $2)
	ON CONFLICT (name)
		DO UPDATE SET value = EXCLUDED.value;
`

func (*configDaoImpl) UpsertStringByName(tx *sql.Tx, name string, value string) error {
	sqlh := selectExecutor(datasource.PostgresHandle(), tx)
	_, err := sqlh.Exec(sqlUpsertConfigStringByName, name, value)
	if err != nil {
		logger.Error(nil, "Got error on ConfigDao/UpsertStringByName query: ", err)
		return err
	}
	return nil
}
