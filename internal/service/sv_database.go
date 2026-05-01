// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"peanut/internal/data"
)

type DatabaseService interface {
	DoesTableExist(ctx context.Context, tableName string) (bool, error)
	GetPostgresVersion(ctx context.Context) (string, error)
}

func NewDatabaseService(metaDao data.MetaDao) DatabaseService {
	return &databaseServiceImpl{metaDao: metaDao}
}

type databaseServiceImpl struct {
	metaDao data.MetaDao
}

func (this *databaseServiceImpl) DoesTableExist(ctx context.Context, tableName string) (bool, error) {
	return this.metaDao.DoesTableExist(ctx, tableName)
}

func (this *databaseServiceImpl) GetPostgresVersion(ctx context.Context) (string, error) {
	return this.metaDao.SelectVersion(ctx)
}
