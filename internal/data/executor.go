// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"database/sql"
	"net/http"
	"peanut/internal/data/datasource"
)

type sqlExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

func getSqlExecutorFromRequest(r *http.Request) sqlExecutor {
	tx, ok := r.Context().Value("tx").(*sql.Tx)
	if !ok || tx == nil {
		return datasource.PostgresHandle()
	}
	return tx
}

func selectExecutor(db *sql.DB, tx *sql.Tx) sqlExecutor {
	if tx == nil {
		return db
	} else {
		return tx
	}
}
