// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"context"
	"database/sql"
	"errors"
	"peanut/internal/logger"
	"time"
)

type ConfigRow struct {
	Name        string
	ValueInt    sql.NullInt64
	ValueString sql.NullString
	ValueUuid   sql.NullString
	Created     time.Time
	Updated     time.Time
}

type ConfigDao interface {
	CreateDBObjects(ctx context.Context) error
	SelectRowByName(ctx context.Context, name string) (*ConfigRow, error)
	UpsertIntByName(ctx context.Context, name string, value int64) error
	UpsertStringByName(ctx context.Context, name string, value string) error
	UpsertUuidByName(ctx context.Context, name string, value string) error
}

func NewConfigDao() ConfigDao {
	return &configDaoImpl{}
}

type configDaoImpl struct{}

var sqlCreateTableConfig = `
	CREATE TABLE config (
		name VARCHAR(100) PRIMARY KEY,
		value_int BIGINT,
		value_string VARCHAR(2000),
		value_uuid UUID,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL,
		CONSTRAINT config_exactly_one_value CHECK (
			num_nonnulls(value_int, value_string, value_uuid) = 1
		)
	);

	CREATE TRIGGER
		config_trigger_created_updated_before_insert
	BEFORE INSERT ON
		config
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		config_trigger_created_updated_before_update
	BEFORE UPDATE ON
		config
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func (*configDaoImpl) CreateDBObjects(ctx context.Context) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlCreateTableConfig)
	if err != nil {
		logger.Error(ctx, "Got database error on ConfigDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlSelectConfigRowByName = "SELECT name, value_int, value_string, value_uuid, _created, _updated FROM config WHERE name = $1"

func (*configDaoImpl) SelectRowByName(ctx context.Context, name string) (*ConfigRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	result := &ConfigRow{}
	row := sqlh.QueryRow(sqlSelectConfigRowByName, name)
	err := row.Scan(&result.Name, &result.ValueInt, &result.ValueString, &result.ValueUuid, &result.Created, &result.Updated)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		logger.Error(ctx, "Got database error on ConfigDao/SelectRowByName query: ", err)
		return nil, err
	}
	return result, nil
}

var sqlUpsertConfigIntByName = `
	INSERT INTO
		config (name, value_int)
	VALUES
		($1, $2)
	ON CONFLICT (name)
		DO UPDATE SET
			value_int = EXCLUDED.value_int,
			value_string = NULL,
			value_uuid = NULL;
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
		config (name, value_string)
	VALUES
		($1, $2)
	ON CONFLICT (name)
		DO UPDATE SET
			value_string = EXCLUDED.value_string,
			value_int = NULL,
			value_uuid = NULL;
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

var sqlUpsertConfigUuidByName = `
	INSERT INTO
		config (name, value_uuid)
	VALUES
		($1, $2::uuid)
	ON CONFLICT (name)
		DO UPDATE SET
			value_uuid = EXCLUDED.value_uuid,
			value_int = NULL,
			value_string = NULL;
`

func (*configDaoImpl) UpsertUuidByName(ctx context.Context, name string, value string) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlUpsertConfigUuidByName, name, value)
	if err != nil {
		logger.Error(ctx, "Got database error on ConfigDao/UpsertUuidByName query: ", err)
		return err
	}
	return nil
}
