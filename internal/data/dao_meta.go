// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"database/sql"
	"net/http"
	"peanut/internal/data/datasource"
	"peanut/internal/logger"
)

type MetaDao interface {
	CreateDBObjects(req *http.Request) error
	DoesTableExist(tableName string) (bool, error)
	SelectVersion() (string, error)
}

type metaDaoImpl struct{}

func (*metaDaoImpl) SelectRowByName(name string, tx *sql.Tx) {
	//TODO implement me
	panic("implement me")
}

func NewMetaDao() MetaDao {
	return &metaDaoImpl{}
}

func (*metaDaoImpl) CreateDBObjects(req *http.Request) error {
	sqlh := getSqlExecutorFromRequest(req)
	_, errInsert := sqlh.Exec(sqlCreatedUpdatedBeforeInsert)
	if errInsert != nil {
		return errInsert
	}
	_, errUpdate := sqlh.Exec(sqlCreatedUpdatedBeforeUpdate)
	if errUpdate != nil {
		return errUpdate
	}
	return nil
}

func (*metaDaoImpl) DoesTableExist(tableName string) (bool, error) {
	sqlh := selectExecutor(datasource.PostgresHandle(), nil)
	rows, err := sqlh.Query("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1", tableName)
	if err != nil {
		logger.Warn(nil, "Error querying in data_meta.DoesTableExist:", tableName, err)
		return false, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Warn(nil, "Error closing rows in data_meta.DoesTableExist:", err)
		}
	}(rows)

	theCount := 0
	for rows.Next() {
		rowErr := rows.Scan(&theCount)
		if rowErr != nil {
			logger.Warn(nil, "Error reading query result in data_meta.DoesTableExist:", tableName, err)
			return false, rowErr
		}
		break
	}
	return theCount > 0, nil
}

var sqlShowServerVersion = "SHOW server_version;"

func (*metaDaoImpl) SelectVersion() (string, error) {
	sqlh := selectExecutor(datasource.PostgresHandle(), nil)

	var version string
	row := sqlh.QueryRow(sqlShowServerVersion)
	err := row.Scan(&version)
	if err != nil {
		return "", err
	}
	return version, nil
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
