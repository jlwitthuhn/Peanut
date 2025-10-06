// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package requtil

import (
	"database/sql"
	"errors"
	"net/http"
	"peanut/internal/logger"
)

func CommitTransactionForRequest(req *http.Request) error {
	tx, ok := req.Context().Value("tx").(*sql.Tx)
	if !ok || tx == nil {
		logger.Error(req, "Attempted to commit transaction with no transaction.")
		return errors.New("No transaction to commit.")
	}

	err := tx.Commit()
	if err != nil {
		logger.Error(req, "Error committing transaction:", err.Error())
		return err
	}

	return nil
}
