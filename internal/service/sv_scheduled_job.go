// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"net/http"
	"peanut/internal/data"
	"time"
)

type ScheduledJobService interface {
	AddJobDefinition(req *http.Request, jobName string, runInterval time.Duration) error
	GetAllJobSummaries(req *http.Request) ([]data.ScheduledJobSummary, error)
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
