// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"errors"
	"net/http"
	"peanut/internal/data"
	"peanut/internal/middleutil"
	"peanut/internal/sched_job"
	"peanut/internal/security/perms"
	"time"
)

type ScheduledJobService interface {
	AddJobDefinition(req *http.Request, jobName string, runInterval time.Duration) error
	GetAllJobSummaries(req *http.Request) ([]data.ScheduledJobSummary, error)
	GetJobNameById(req *http.Request, id string) (string, error)
	RunJob(req *http.Request, jobName string) error
}

func NewScheduledJobService(multiTableDao data.MultiTableDao, scheduledJobDao data.ScheduledJobDao, scheduledJobRunDao data.ScheduledJobRunDao) ScheduledJobService {
	return &scheduledJobServiceImpl{multiTableDao: multiTableDao, scheduledJobDao: scheduledJobDao, scheduledJobRunDao: scheduledJobRunDao}
}

type scheduledJobServiceImpl struct {
	multiTableDao      data.MultiTableDao
	scheduledJobDao    data.ScheduledJobDao
	scheduledJobRunDao data.ScheduledJobRunDao
}

func (this *scheduledJobServiceImpl) AddJobDefinition(req *http.Request, jobName string, runInterval time.Duration) error {
	return this.scheduledJobDao.InsertRow(req, jobName, runInterval)
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
		scheduled_job.RunExpiredSessionsJob()
	} else {
		return errors.New("not implemented")
	}

	err = this.scheduledJobRunDao.InsertRow(req, jobDetails.Id, true)
	if err != nil {
		return err
	}

	return nil
}
