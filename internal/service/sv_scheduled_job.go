// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"peanut/internal/data"
	"peanut/internal/data/datasource"
	"peanut/internal/keynames/contextkeys"
	"peanut/internal/logger"
	"peanut/internal/middleutil"
	"peanut/internal/security/perms"
	"time"
)

type ScheduledJobService interface {
	AddJobDefinition(req *http.Request, jobName string, runInterval time.Duration) error
	BackgroundThreadFunc()
	GetAllJobSummaries(req *http.Request) ([]data.ScheduledJobSummary, error)
	GetJobNameById(req *http.Request, id string) (string, error)
	RunJob(req *http.Request, jobName string) error
}

func NewScheduledJobService(multiTableDao data.MultiTableDao, scheduledJobDao data.ScheduledJobDao, scheduledJobRunDao data.ScheduledJobRunDao, sessionDao data.SessionDao, dbService DatabaseService) ScheduledJobService {
	return &scheduledJobServiceImpl{multiTableDao: multiTableDao, scheduledJobDao: scheduledJobDao, scheduledJobRunDao: scheduledJobRunDao, sessionDao: sessionDao, dbService: dbService}
}

type scheduledJobServiceImpl struct {
	multiTableDao      data.MultiTableDao
	scheduledJobDao    data.ScheduledJobDao
	scheduledJobRunDao data.ScheduledJobRunDao
	sessionDao         data.SessionDao
	dbService          DatabaseService
}

func (this *scheduledJobServiceImpl) AddJobDefinition(req *http.Request, jobName string, runInterval time.Duration) error {
	return this.scheduledJobDao.InsertRow(req, jobName, runInterval)
}

// Application infrastructure expects each db access to be associated with a request
// We make a fake request here for all background db operations that are not associated with a real request
// The request here is never sent
func createBackgroundHttpRequest() (*http.Request, error) {
	result := httptest.NewRequest(http.MethodGet, "http://127.0.0.1", nil)

	result = result.WithContext(context.WithValue(result.Context(), contextkeys.RequestId, "BGTHREAD"))

	// Set up db transaction
	tx, err := datasource.PostgresHandle().BeginTx(result.Context(), nil)
	if err != nil {
		logger.Error(result, "Failed to create db transaction for scheduled jobs")
		return nil, err
	}
	ctx := context.WithValue(result.Context(), contextkeys.PostgresTx, tx)
	result = result.WithContext(ctx)

	// Add permissions
	permissions := []string{perms.Admin_ScheduledJob_Run}
	ctx = context.WithValue(result.Context(), contextkeys.UserPerms, permissions)
	result = result.WithContext(ctx)

	return result, nil
}

func (this *scheduledJobServiceImpl) backgroundThreadIter() {
	req, err := createBackgroundHttpRequest()
	if err != nil {
		return
	}
	tx, txOk := req.Context().Value(contextkeys.PostgresTx).(*sql.Tx)
	if !txOk {
		logger.Debug(req, "Database not yet initialized, waiting another cycle...")
		return
	}
	defer tx.Rollback()

	// Check if the database exists
	exists, err := this.dbService.DoesTableExist(req, "config_int")
	if err != nil {
		logger.Error(req, "Failed to check if database exists")
		return
	}
	if !exists {
		logger.Debug(req, "Database not yet initialized, waiting another cycle...")
		return
	}

	row, err := this.multiTableDao.SelectScheduledJobByNextPending(req)
	if err != nil {
		logger.Error(req, "Failed to find next pending scheduled job, aborting")
		return
	}

	logger.Debug(req, "Running job: "+row.Name)
	err = this.RunJob(req, row.Name)
	if err != nil {
		logger.Error(req, "Failed to run job, aborting")
		return
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(req, "Failed to commit transaction")
	}
}

func (this *scheduledJobServiceImpl) BackgroundThreadFunc() {
	for {
		this.backgroundThreadIter()
		time.Sleep(1 * time.Minute)
	}
}

func (this *scheduledJobServiceImpl) GetAllJobSummaries(req *http.Request) ([]data.ScheduledJobSummary, error) {
	return this.multiTableDao.SelectAllScheduledJobSummaries(req)
}

func (this *scheduledJobServiceImpl) GetJobNameById(req *http.Request, id string) (string, error) {
	row, err := this.scheduledJobDao.SelectRowById(req, id)
	if err != nil {
		return "", err
	}
	return row.Name, nil
}

func (this *scheduledJobServiceImpl) RunJob(req *http.Request, jobName string) error {
	if middleutil.RequestHasPermission(req, perms.Admin_ScheduledJob_Run) == false {
		return errors.New("permission denied")
	}

	jobDetails, err := this.scheduledJobDao.SelectRowByName(req, jobName)
	if err != nil {
		return err
	}

	if jobName == "DeleteExpiredSessions" {
		this.runExpiredSessionsJob(req)
	} else {
		return errors.New("not implemented")
	}

	err = this.scheduledJobRunDao.InsertRow(req, jobDetails.Id, true)
	if err != nil {
		return err
	}

	return nil
}

func (this *scheduledJobServiceImpl) runExpiredSessionsJob(req *http.Request) error {
	logger.Info(req, "Expired sessions job beginning")
	err := this.sessionDao.DeleteRowsByExpired(req)
	if err != nil {
		logger.Error(req, "Failed to delete expired sessions")
		return errors.New("failed to delete expired sessions")
	}
	logger.Debug(req, "Expired sessions job complete")
	return nil
}
