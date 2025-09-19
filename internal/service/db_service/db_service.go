// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package db_service

import "peanut/internal/database/data/data_meta"

func DoesTableExist(tableName string) (bool, error) {
	return data_meta.DoesTableExist(tableName)
}
