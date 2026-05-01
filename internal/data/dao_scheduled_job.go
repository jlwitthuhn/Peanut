// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"context"
	"peanut/internal/data/dataformat"
	"peanut/internal/logger"
	"time"
)

type ScheduledJobRow struct {
	Id          string
	Name        string
	RunInterval string
	Created     time.Time
	Updated     time.Time
}

type ScheduledJobDao interface {
	CreateDBObjects(ctx context.Context) error
	InsertRow(ctx context.Context, name string, runInterval time.Duration) error
	SelectRowByName(ctx context.Context, name string) (*ScheduledJobRow, error)
	SelectRowById(ctx context.Context, id string) (*ScheduledJobRow, error)
}

func NewScheduledJobDao() ScheduledJobDao {
	return &scheduledJobDaoImpl{}
}

type scheduledJobDaoImpl struct{}

var sqlCreateTableScheduledJobs = `
	CREATE TABLE scheduled_jobs (
		id UUID PRIMARY KEY DEFAULT uuidv7(),
	    name VARCHAR(100) NOT NULL UNIQUE,
	    run_interval INTERVAL NOT NULL,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE TRIGGER
		scheduled_jobs_trigger_created_updated_before_insert
	BEFORE INSERT ON
		scheduled_jobs
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		scheduled_jobs_trigger_created_updated_before_update
	BEFORE UPDATE ON
		scheduled_jobs
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func (this *scheduledJobDaoImpl) CreateDBObjects(ctx context.Context) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlCreateTableScheduledJobs)
	if err != nil {
		logger.Error(ctx, "Got error on CreateDBObjects query:", err)
		return err
	}
	return nil
}

var sqlInsertScheduledJobsRow = "INSERT INTO scheduled_jobs(name, run_interval) VALUES ($1, $2)"

func (this *scheduledJobDaoImpl) InsertRow(ctx context.Context, name string, runInterval time.Duration) error {
	sqlh := getSqlExecutorFromContext(ctx)
	formattedInterval := dataformat.FormatDurationAsPostgresInterval(runInterval)
	_, err := sqlh.Exec(sqlInsertScheduledJobsRow, name, formattedInterval)
	if err != nil {
		logger.Error(ctx, "Got error on InsertRow query:", err)
		return err
	}
	return nil
}

var sqlSelectScheduledJobsRowByName = `
	SELECT
		id, name, run_interval, _created, _updated
	FROM scheduled_jobs
		WHERE name = $1
`

func (this *scheduledJobDaoImpl) SelectRowByName(ctx context.Context, name string) (*ScheduledJobRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	result := &ScheduledJobRow{}
	row := sqlh.QueryRow(sqlSelectScheduledJobsRowByName, name)
	err := row.Scan(&result.Id, &result.Name, &result.RunInterval, &result.Created, &result.Updated)
	if err != nil {
		logger.Error(ctx, "Got error on SelectRowByName query:", err)
		return nil, err
	}
	return result, nil
}

var sqlSelectScheduledJobsRowById = `
	SELECT
		id, name, run_interval, _created, _updated
	FROM scheduled_jobs
		WHERE id = $1
`

func (this *scheduledJobDaoImpl) SelectRowById(ctx context.Context, id string) (*ScheduledJobRow, error) {
	sqlh := getSqlExecutorFromContext(ctx)
	result := &ScheduledJobRow{}
	row := sqlh.QueryRow(sqlSelectScheduledJobsRowById, id)
	err := row.Scan(&result.Id, &result.Name, &result.RunInterval, &result.Created, &result.Updated)
	if err != nil {
		logger.Error(ctx, "Got error on SelectRowById query:", err)
		return nil, err
	}
	return result, nil
}
