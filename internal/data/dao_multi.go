// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"net/http"
)

type MultiTableDao interface {
	SelectGroupNamesByUserId(r *http.Request, userId string) ([]string, error)
}

func NewMultiTableDao() MultiTableDao {
	return &multiTableDaoImpl{}
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

func (*multiTableDaoImpl) SelectGroupNamesByUserId(r *http.Request, userId string) ([]string, error) {
	sqlh := getSqlExecutorFromRequest(r)

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
