// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"net/http"
	"peanut/internal/data"
)

type DatabaseService interface {
	DoesTableExist(req *http.Request, tableName string) (bool, error)
	GetPostgresVersion(req *http.Request) (string, error)
}

func NewDatabaseService(metaDao data.MetaDao) DatabaseService {
	return &databaseServiceImpl{metaDao: metaDao}
}

type databaseServiceImpl struct {
	metaDao data.MetaDao
}

func (this *databaseServiceImpl) DoesTableExist(req *http.Request, tableName string) (bool, error) {
	return this.metaDao.DoesTableExist(req, tableName)
}

func (this *databaseServiceImpl) GetPostgresVersion(req *http.Request) (string, error) {
	return this.metaDao.SelectVersion(req)
}
