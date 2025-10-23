// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"errors"
	"net/http"
	"peanut/internal/data"
	"time"
)

type ScheduledJobService interface {
	AddJobDefinition(req *http.Request, jobName string, runInterval time.Duration) error
	GetAllJobSummaries(req *http.Request) ([]data.ScheduledJobSummary, error)
	GetJobNameById(req *http.Request, id string) (string, error)
	RunJob(req *http.Request, jobName string) error
}

func NewScheduledJobService(multiTableDao data.MultiTableDao, scheduledJobDao data.ScheduledJobDao) ScheduledJobService {
	return &scheduledJobServiceImpl{multiTableDao: multiTableDao, scheduledJobDao: scheduledJobDao}
}

type scheduledJobServiceImpl struct {
	multiTableDao   data.MultiTableDao
	scheduledJobDao data.ScheduledJobDao
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
	return errors.New("not implemented")
}
