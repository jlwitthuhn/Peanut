// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"database/sql"
	"peanut/internal/data/datasource"
	"sync"
)

type MultiTableDao interface {
	SelectGroupNamesByUserId(tx *sql.Tx, userId string) ([]string, error)
}

var multiTableDaoInstance MultiTableDao
var multiTableDaoInstanceOnce sync.Once

func MultiTableDaoInst() MultiTableDao {
	multiTableDaoInstanceOnce.Do(func() {
		multiTableDaoInstance = &multiTableDaoImpl{}
	})
	return multiTableDaoInstance
}

type multiTableDaoImpl struct{}

var sqlSelectGroupNamesByUserId = `
	SELECT
		groups.name
	FROM
	    groups INNER JOIN group_membership ON groups.id = group_membership.group_id
	WHERE
	    group_membership.user_id = $1
`

func (*multiTableDaoImpl) SelectGroupNamesByUserId(tx *sql.Tx, userId string) ([]string, error) {
	sqlh := selectExecutor(datasource.PostgresHandle(), tx)

	rows, err := sqlh.Query(sqlSelectGroupNamesByUserId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []string{}
	for rows.Next() {
		var thisGroup string
		scanErr := rows.Scan(&thisGroup)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, thisGroup)
	}
	return result, nil
}
