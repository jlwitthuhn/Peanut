// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"net/http"
	"peanut/internal/logger"
)

type ScheduledJobSummary struct {
	Name        string
	RunInterval string
}

type MultiTableDao interface {
	SelectAllScheduledJobSummaries(req *http.Request) ([]ScheduledJobSummary, error)
	SelectGroupNamesByUserId(r *http.Request, userId string) ([]string, error)
	SelectUserRowsByGroupName(r *http.Request, groupName string) ([]UserRow, error)
}

func NewMultiTableDao() MultiTableDao {
	return &multiTableDaoImpl{}
}

type multiTableDaoImpl struct{}

var sqlSelectAllScheduledJobSummaries = `
	SELECT
		name, run_interval
	FROM
	    scheduled_jobs
	ORDER BY
	    name
`

func (*multiTableDaoImpl) SelectAllScheduledJobSummaries(req *http.Request) ([]ScheduledJobSummary, error) {
	sqlh := getSqlExecutorFromRequest(req)
	rows, err := sqlh.Query(sqlSelectAllScheduledJobSummaries)
	if err != nil {
		logger.Error(nil, "Got error on SelectAllScheduledJobSummaries query:", err)
		return nil, err
	}
	defer rows.Close()

	var result []ScheduledJobSummary
	for rows.Next() {
		thisSummary := ScheduledJobSummary{}
		scanErr := rows.Scan(&thisSummary.Name, &thisSummary.RunInterval)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, thisSummary)
	}
	return result, nil
}

var sqlSelectGroupNamesByUserId = `
	SELECT
		groups.name
	FROM
	    groups INNER JOIN group_membership ON groups.id = group_membership.group_id
	WHERE
	    group_membership.user_id = $1
`

func (*multiTableDaoImpl) SelectGroupNamesByUserId(req *http.Request, userId string) ([]string, error) {
	sqlh := getSqlExecutorFromRequest(req)
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

var sqlSelectUsersByGroupName = `
WITH
	q_group_ids AS (
		SELECT
			id
		FROM
			groups
		WHERE
			name = $1
	),

	q_user_ids AS (
		SELECT
			user_id
		FROM
			group_membership
		WHERE
			group_id IN (SELECT id FROM q_group_ids)
	)

	SELECT
		id, display_name, email, password, _created, _updated
	FROM
		users
	WHERE
	    id IN (SELECT user_id FROM q_user_ids)
	ORDER BY
	    _created
`

func (*multiTableDaoImpl) SelectUserRowsByGroupName(req *http.Request, groupName string) ([]UserRow, error) {
	sqlh := getSqlExecutorFromRequest(req)
	rows, err := sqlh.Query(sqlSelectUsersByGroupName, groupName)
	if err != nil {
		logger.Error(nil, "Got error on MultiTableDao/SelectRowsLikeName query:", err)
		return nil, err
	}
	var result []UserRow
	for rows.Next() {
		thisRow := UserRow{}
		err = rows.Scan(&thisRow.Id, &thisRow.DisplayName, &thisRow.Email, &thisRow.Password, &thisRow.Created, &thisRow.Updated)
		if err != nil {
			return nil, err
		}
		result = append(result, thisRow)
	}
	return result, nil
}
