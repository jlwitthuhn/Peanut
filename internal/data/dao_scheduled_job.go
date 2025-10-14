// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package data

import (
	"net/http"
	"peanut/internal/data/dataformat"
	"peanut/internal/logger"
	"time"
)

type ScheduledJobDao interface {
	CreateDBObjects(req *http.Request) error
	InsertRow(req *http.Request, name string, runInterval time.Duration) error
}

func NewScheduledJobDao() ScheduledJobDao {
	return &scheduledJobDaoImpl{}
}

type scheduledJobDaoImpl struct{}

var sqlCreateTableScheduledJobs = `
	CREATE TABLE scheduled_jobs (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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

func (this *scheduledJobDaoImpl) CreateDBObjects(req *http.Request) error {
	sqlh := getSqlExecutorFromRequest(req)
	_, err := sqlh.Exec(sqlCreateTableScheduledJobs)
	if err != nil {
		logger.Error(nil, "Got error on ScheduledJobDao/CreateDBObjects query: ", err)
		return err
	}
	return nil
}

var sqlInsertScheduledJobsRow = "INSERT INTO scheduled_jobs(name, run_interval) VALUES ($1, $2)"

func (this *scheduledJobDaoImpl) InsertRow(req *http.Request, name string, runInterval time.Duration) error {
	sqlh := getSqlExecutorFromRequest(req)
	formattedInterval := dataformat.FormatDurationAsPostgresInterval(runInterval)
	_, err := sqlh.Exec(sqlInsertScheduledJobsRow, name, formattedInterval)
	if err != nil {
		logger.Error(nil, "Got error on ScheduledJobDao/InsertRow query: ", err)
		return err
	}
	return nil
}
