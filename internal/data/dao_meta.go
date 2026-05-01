// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"context"
	"database/sql"
	"peanut/internal/logger"
)

type MetaDao interface {
	CreateDBObjects(ctx context.Context) error
	DoesTableExist(ctx context.Context, tableName string) (bool, error)
	SelectVersion(ctx context.Context) (string, error)
	Vacuum(ctx context.Context, dbh *sql.DB) error
}

type metaDaoImpl struct{}

func (*metaDaoImpl) SelectRowByName(name string, tx *sql.Tx) {
	//TODO implement me
	panic("implement me")
}

func NewMetaDao() MetaDao {
	return &metaDaoImpl{}
}

func (*metaDaoImpl) CreateDBObjects(ctx context.Context) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, errInsert := sqlh.Exec(sqlCreatedUpdatedBeforeInsert)
	if errInsert != nil {
		logger.Error(ctx, "Got database error on MetaDao/CreateDBObjects query: ", errInsert)
		return errInsert
	}
	_, errUpdate := sqlh.Exec(sqlCreatedUpdatedBeforeUpdate)
	if errUpdate != nil {
		logger.Error(ctx, "Got database error on MetaDao/CreateDBObjects query: ", errUpdate)
		return errUpdate
	}
	_, errVis := sqlh.Exec(sqlVisibilityEnum)
	if errVis != nil {
		logger.Error(ctx, "Got database error on MetaDao/CreateDBObjects query: ", errVis)
		return errVis
	}
	return nil
}

func (*metaDaoImpl) DoesTableExist(ctx context.Context, tableName string) (bool, error) {
	sqlh := getSqlExecutorFromContext(ctx)

	rows, err := sqlh.Query("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1", tableName)
	if err != nil {
		logger.Error(ctx, "Got database error on MetaDao/DoesTableExist query: ", err)
		return false, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error(ctx, "Got database error on MetaDao/DoesTableExist query: ", err)
		}
	}(rows)

	theCount := 0
	for rows.Next() {
		rowErr := rows.Scan(&theCount)
		if rowErr != nil {
			logger.Error(ctx, "Got database error on MetaDao/DoesTableExist query: ", rowErr)
			return false, rowErr
		}
		break
	}
	return theCount > 0, nil
}

var sqlShowServerVersion = "SHOW server_version;"

func (*metaDaoImpl) SelectVersion(ctx context.Context) (string, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	var version string
	row := sqlh.QueryRow(sqlShowServerVersion)
	err := row.Scan(&version)
	if err != nil {
		logger.Error(ctx, "Got database error on MetaDao/SelectVersion query: ", err)
		return "", err
	}
	return version, nil
}

var sqlVacuumDb = "VACUUM;"

func (*metaDaoImpl) Vacuum(ctx context.Context, dbh *sql.DB) error {
	_, err := dbh.Exec(sqlVacuumDb)
	if err != nil {
		logger.Error(ctx, "Got database error on MetaDao/Vacuum query: ", err)
		return err
	}
	return nil
}

var sqlCreatedUpdatedBeforeInsert = `
	CREATE FUNCTION fn_created_updated_before_insert()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW._created := now();
		NEW._updated := NEW._created;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
`

var sqlCreatedUpdatedBeforeUpdate = `
	CREATE FUNCTION fn_created_updated_before_update()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW._created := OLD._created;
		NEW._updated := now();
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
`

var sqlVisibilityEnum = `
	CREATE TYPE visibility_enum
	AS ENUM('Public', 'Private', 'Deleted');
`
