// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package database

import (
	"database/sql"
	"peanut/internal/logger"
	"sync"
)

type MetaDao struct{}

var metaDaoInstance *MetaDao
var metaDaoInstanceOnce sync.Once

func MetaDaoInst() *MetaDao {
	metaDaoInstanceOnce.Do(func() {
		metaDaoInstance = &MetaDao{}
	})
	return metaDaoInstance
}

func (*MetaDao) CreateDBObjects(tx *sql.Tx) error {
	_, errInsert := tx.Exec(sqlCreatedUpdatedBeforeInsert)
	if errInsert != nil {
		return errInsert
	}
	_, errUpdate := tx.Exec(sqlCreatedUpdatedBeforeUpdate)
	if errUpdate != nil {
		return errUpdate
	}
	return nil
}

func (dao *MetaDao) DoesTableExist(tableName string) (bool, error) {
	db := PostgresHandle()
	rows, err := db.Query("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1", tableName)
	if err != nil {
		logger.Warn("Error querying in data_meta.DoesTableExist:", tableName, err)
		return false, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Warn("Error closing rows in data_meta.DoesTableExist:", err)
		}
	}(rows)

	theCount := 0
	for rows.Next() {
		rowErr := rows.Scan(&theCount)
		if rowErr != nil {
			logger.Warn("Error reading query result in data_meta.DoesTableExist:", tableName, err)
			return false, rowErr
		}
		break
	}
	return theCount > 0, nil
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
