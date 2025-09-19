// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data_meta

import (
	"database/sql"

	"peanut/internal/database"
	"peanut/internal/logger"
)

func DoesTableExist(tableName string) (bool, error) {
	db := database.PostgresHandle()
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
