// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"database/sql"
	"errors"
	"peanut/internal/data"
	"peanut/internal/data/datasource"
	"peanut/internal/keynames/contextkeys"
	"peanut/internal/logger"
	"peanut/internal/middleutil"
	"peanut/internal/security/perms"
	"time"
)

type ScheduledJobService interface {
	AddJobDefinition(ctx context.Context, jobName string, runInterval time.Duration) error
	BackgroundThreadFunc()
	GetAllJobSummaries(ctx context.Context) ([]data.ScheduledJobSummary, error)
	GetJobNameById(ctx context.Context, id string) (string, error)
	RunJob(ctx context.Context, jobName string) error
}

func NewScheduledJobService(
	metaDao data.MetaDao,
	multiTableDao data.MultiTableDao,
	scheduledJobDao data.ScheduledJobDao,
	scheduledJobRunDao data.ScheduledJobRunDao,
	sessionDao data.SessionDao,
	dbService DatabaseService,
) ScheduledJobService {
	return &scheduledJobServiceImpl{
		metaDao:            metaDao,
		multiTableDao:      multiTableDao,
		scheduledJobDao:    scheduledJobDao,
		scheduledJobRunDao: scheduledJobRunDao,
		sessionDao:         sessionDao,
		dbService:          dbService,
	}
}

type scheduledJobServiceImpl struct {
	metaDao            data.MetaDao
	multiTableDao      data.MultiTableDao
	scheduledJobDao    data.ScheduledJobDao
	scheduledJobRunDao data.ScheduledJobRunDao
	sessionDao         data.SessionDao
	dbService          DatabaseService
}

func (this *scheduledJobServiceImpl) AddJobDefinition(ctx context.Context, jobName string, runInterval time.Duration) error {
	return this.scheduledJobDao.InsertRow(ctx, jobName, runInterval)
}

// Application infrastructure expects each db access to be associated with a request.
// We create a synthetic context here for all background db operations that are not
// associated with a real request.
func createBackgroundContext() (context.Context, error) {
	ctx := context.WithValue(context.Background(), contextkeys.RequestId, "BGTHREAD")

	// Set up db transaction
	tx, err := datasource.PostgresHandle().BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "Failed to create db transaction for scheduled jobs")
		return nil, err
	}
	ctx = context.WithValue(ctx, contextkeys.PostgresTx, tx)

	// Add permissions
	permissions := map[string]struct{}{perms.Admin_ScheduledJob_Run: {}}
	ctx = context.WithValue(ctx, contextkeys.UserPerms, permissions)

	return ctx, nil
}

func (this *scheduledJobServiceImpl) backgroundThreadIter() {
	ctx, err := createBackgroundContext()
	if err != nil {
		return
	}

	tx, txOk := ctx.Value(contextkeys.PostgresTx).(*sql.Tx)
	if !txOk {
		logger.Debug(ctx, "Database not yet initialized, waiting another cycle...")
		return
	}
	defer tx.Rollback()

	// Check if the database exists
	exists, err := this.dbService.DoesTableExist(ctx, "config_int")
	if err != nil {
		logger.Error(ctx, "Failed to check if database exists")
		return
	}
	if !exists {
		logger.Debug(ctx, "Database not yet initialized, waiting another cycle...")
		return
	}

	row, err := this.multiTableDao.SelectScheduledJobByNextPending(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to find next pending scheduled job, aborting")
		return
	}
	if row == nil {
		logger.Trace(ctx, "Nothing to do, waiting another cycle...")
		return
	}

	logger.Debug(ctx, "Running job: "+row.Name)
	err = this.RunJob(ctx, row.Name)
	if err != nil {
		logger.Error(ctx, "Failed to run job, aborting")
		return
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(ctx, "Failed to commit transaction")
	}
}

func (this *scheduledJobServiceImpl) BackgroundThreadFunc() {
	for {
		this.backgroundThreadIter()
		time.Sleep(1 * time.Minute)
	}
}

func (this *scheduledJobServiceImpl) GetAllJobSummaries(ctx context.Context) ([]data.ScheduledJobSummary, error) {
	return this.multiTableDao.SelectAllScheduledJobSummaries(ctx)
}

func (this *scheduledJobServiceImpl) GetJobNameById(ctx context.Context, id string) (string, error) {
	row, err := this.scheduledJobDao.SelectRowById(ctx, id)
	if err != nil {
		return "", err
	}
	return row.Name, nil
}

func (this *scheduledJobServiceImpl) RunJob(ctx context.Context, jobName string) error {
	if middleutil.ContextHasPermission(ctx, perms.Admin_ScheduledJob_Run) == false {
		return errors.New("permission denied")
	}

	jobDetails, err := this.scheduledJobDao.SelectRowByName(ctx, jobName)
	if err != nil {
		return err
	}

	if jobName == "DeleteExpiredSessions" {
		err = this.runExpiredSessionsJob(ctx)
		if err != nil {
			return err
		}
	} else if jobName == "VacuumDatabase" {
		err = this.runVacuumDbJob(ctx)
		if err != nil {
			return err
		}
	} else {
		return errors.New("not implemented")
	}

	err = this.scheduledJobRunDao.InsertRow(ctx, jobDetails.Id, true)
	if err != nil {
		return err
	}

	return nil
}

func (this *scheduledJobServiceImpl) runExpiredSessionsJob(ctx context.Context) error {
	logger.Info(ctx, "Expired sessions job beginning")
	err := this.sessionDao.DeleteRowsByExpired(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to delete expired sessions")
		return errors.New("failed to delete expired sessions")
	}
	logger.Debug(ctx, "Expired sessions job complete")
	return nil
}

func (this *scheduledJobServiceImpl) runVacuumDbJob(ctx context.Context) error {
	logger.Info(ctx, "Vacuum database job beginning")
	dbh := datasource.PostgresHandle()
	if dbh == nil {
		logger.Error(ctx, "No postgres handle, aborting")
		return errors.New("no postgres handle")
	}
	err := this.metaDao.Vacuum(ctx, dbh)
	if err != nil {
		logger.Error(ctx, "Failed to vacuum:", err)
		return errors.New("failed to vacuum")
	}
	logger.Debug(ctx, "Vacuum database job complete")
	return nil
}
