// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"database/sql"
	"errors"
	"net/http"
	"peanut/internal/logger"
)

type ScheduledJobSummary struct {
	JobId            string
	Name             string
	RunInterval      string
	LastRunTime      string
	LastRunResult    bool
	LastRunTimeRaw   sql.NullTime
	LastRunResultRaw sql.NullBool
}

func (this *ScheduledJobSummary) populateStrings() {
	if this.LastRunTimeRaw.Valid {
		this.LastRunTime = this.LastRunTimeRaw.Time.Format("2006-01-02 15:04:05 MST")
	} else {
		this.LastRunTime = "Never"
	}
	if this.LastRunResultRaw.Valid {
		this.LastRunResult = this.LastRunResultRaw.Bool
	} else {
		this.LastRunResult = false
	}
}

type MultiTableDao interface {
	SelectAllScheduledJobSummaries(req *http.Request) ([]ScheduledJobSummary, error)
	SelectGroupNamesByUserId(r *http.Request, userId string) ([]string, error)
	SelectScheduledJobByNextPending(r *http.Request) (*ScheduledJobRow, error)
	SelectUserRowsByGroupName(r *http.Request, groupName string) ([]UserRow, error)
}

func NewMultiTableDao() MultiTableDao {
	return &multiTableDaoImpl{}
}

type multiTableDaoImpl struct{}

var sqlSelectAllScheduledJobSummaries = `
	WITH last_run AS (
		SELECT
			DISTINCT ON (job_id)
			job_id, success, _created
		FROM
			scheduled_job_runs
		ORDER BY
			job_id, _created DESC
	)
	
	SELECT
		scheduled_jobs.id, name, run_interval, last_run._created, last_run.success
	FROM
	    scheduled_jobs LEFT JOIN last_run ON scheduled_jobs.id = last_run.job_id
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
		scanErr := rows.Scan(&thisSummary.JobId, &thisSummary.Name, &thisSummary.RunInterval, &thisSummary.LastRunTimeRaw, &thisSummary.LastRunResultRaw)
		if scanErr != nil {
			return nil, scanErr
		}
		thisSummary.populateStrings()
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

var sqlSelectScheduledJobByNextPending = `
	WITH
	all_job_ids AS (
		SELECT
			id,
			name,
			run_interval,
			_created,
			_updated
		FROM
			scheduled_jobs
	),

	latest_success AS (
		SELECT DISTINCT ON (job_id)
			job_id,
			_created as run_time
		FROM
			scheduled_job_runs
		ORDER BY
			job_id, _created DESC
	),

	full_successes AS (
		SELECT
			all_job_ids.id AS job_id,
			all_job_ids.name AS job_name,
			all_job_ids.run_interval AS run_interval,
			all_job_ids._created AS created,
			all_job_ids._updated AS updated,
			COALESCE(NOW() - latest_success.run_time, INTERVAL '1 day') AS since_last_run
		FROM
			all_job_ids LEFT JOIN latest_success ON all_job_ids.id = latest_success.job_id
	),

	full_successes_ordered AS (
		SELECT
			job_id, job_name, run_interval, created, updated
		FROM
			full_successes
		WHERE
			since_last_run >= run_interval
		ORDER BY
			since_last_run DESC
	)

	SELECT
		job_id, job_name, run_interval, created, updated
	FROM
		full_successes_ordered
	FETCH FIRST 1 ROWS ONLY
`

func (*multiTableDaoImpl) SelectScheduledJobByNextPending(req *http.Request) (*ScheduledJobRow, error) {
	sqlh := getSqlExecutorFromRequest(req)
	result := &ScheduledJobRow{}
	row := sqlh.QueryRow(sqlSelectScheduledJobByNextPending)
	err := row.Scan(&result.Id, &result.Name, &result.RunInterval, &result.Created, &result.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Error(nil, "Got error on MultiTableDao/SelectScheduledJobByNextPending query:", err)
		return nil, err
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
	defer rows.Close()

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
