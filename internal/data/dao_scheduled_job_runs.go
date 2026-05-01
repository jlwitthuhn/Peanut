// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"context"
	"peanut/internal/logger"
)

type ScheduledJobRunDao interface {
	CreateDBObjects(ctx context.Context) error
	InsertRow(ctx context.Context, jobId string, success bool) error
}

func NewScheduledJobRunDao() ScheduledJobRunDao {
	return &scheduledJobRunDaoImpl{}
}

type scheduledJobRunDaoImpl struct{}

var sqlCreateTableScheduledJobRuns = `
	CREATE TABLE scheduled_job_runs (
		id UUID PRIMARY KEY DEFAULT uuidv7(),
		job_id UUID NOT NULL REFERENCES scheduled_jobs(id) ON DELETE CASCADE,
	    success BOOLEAN NOT NULL,
		_created TIMESTAMP WITH TIME ZONE NOT NULL,
		_updated TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE INDEX ON scheduled_job_runs (job_id);

	CREATE TRIGGER
		scheduled_job_runs_trigger_created_updated_before_insert
	BEFORE INSERT ON
		scheduled_job_runs
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_insert();

	CREATE TRIGGER
		scheduled_job_runs_trigger_created_updated_before_update
	BEFORE UPDATE ON
		scheduled_job_runs
	FOR EACH ROW EXECUTE FUNCTION
		fn_created_updated_before_update();
`

func (*scheduledJobRunDaoImpl) CreateDBObjects(ctx context.Context) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlCreateTableScheduledJobRuns)
	if err != nil {
		logger.Error(ctx, "Got error on ScheduledJobRunDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlInsertScheduledJobRunRow = "INSERT INTO scheduled_job_runs(job_id, success) VALUES ($1, $2)"

func (*scheduledJobRunDaoImpl) InsertRow(ctx context.Context, jobId string, success bool) error {
	sqlh := getSqlExecutorFromContext(ctx)
	_, err := sqlh.Exec(sqlInsertScheduledJobRunRow, jobId, success)
	if err != nil {
		logger.Error(ctx, "Got error on InsertRow query:", err)
		return err
	}
	return nil
}
