// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"context"
	"peanut/internal/logger"
)

type ConfigIntRow struct {
	Name  string
	Value int64
}

type ConfigStringRow struct {
	Name  string
	Value string
}

type ConfigDao interface {
	CreateDBObjects(ctx context.Context) error
	SelectIntRowByName(ctx context.Context, name string) (*ConfigIntRow, error)
	SelectStringRowByName(ctx context.Context, name string) (*ConfigStringRow, error)
	UpsertIntByName(ctx context.Context, name string, value int64) error
	UpsertStringByName(ctx context.Context, name string, value string) error
}

func NewConfigDao() ConfigDao {
	return &configDaoImpl{}
}

type configDaoImpl struct{}

var sqlCreateTableConfigInt = `
	CREATE TABLE config_int (
		name VARCHAR(100) PRIMARY KEY,
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
		name VARCHAR(100) PRIMARY KEY,
		value VARCHAR(2000) NOT NULL,
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

func (*configDaoImpl) CreateDBObjects(ctx context.Context) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlCreateTableConfigInt)
	if err != nil {
		logger.Error(ctx, "Got database error on ConfigDao/CreateDBObjects query: ", err)
		return err
	}
	_, err = sqlh.Exec(sqlCreateTableConfigString)
	if err != nil {
		logger.Error(ctx, "Got database error on ConfigDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlSelectConfigIntRowByName = "SELECT name, value FROM config_int WHERE name = $1"

func (*configDaoImpl) SelectIntRowByName(ctx context.Context, name string) (*ConfigIntRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	result := &ConfigIntRow{}
	row := sqlh.QueryRow(sqlSelectConfigIntRowByName, name)
	err := row.Scan(&result.Name, &result.Value)
	if err != nil {
		logger.Error(ctx, "Got database error on ConfigDao/SelectIntRowByName query: ", err)
		return nil, err
	}
	return result, nil
}

var sqlSelectConfigStringRowByName = "SELECT name, value FROM config_string WHERE name = $1"

func (*configDaoImpl) SelectStringRowByName(ctx context.Context, name string) (*ConfigStringRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	result := &ConfigStringRow{}
	row := sqlh.QueryRow(sqlSelectConfigStringRowByName, name)
	err := row.Scan(&result.Name, &result.Value)
	if err != nil {
		logger.Error(ctx, "Got database error on ConfigDao/SelectStringRowByName query: ", err)
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

func (*configDaoImpl) UpsertIntByName(ctx context.Context, name string, value int64) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlUpsertConfigIntByName, name, value)
	if err != nil {
		logger.Error(ctx, "Got database error on ConfigDao/UpsertIntByName query: ", err)
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

func (*configDaoImpl) UpsertStringByName(ctx context.Context, name string, value string) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlUpsertConfigStringByName, name, value)
	if err != nil {
		logger.Error(ctx, "Got database error on ConfigDao/UpsertStringByName query: ", err)
		return err
	}
	return nil
}
