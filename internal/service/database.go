// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"peanut/internal/data"
)

type DatabaseService interface {
	DoesTableExist(tableName string) (bool, error)
}

func NewDatabaseService() DatabaseService {
	return &databaseServiceImpl{}
}

type databaseServiceImpl struct{}

func (*databaseServiceImpl) DoesTableExist(tableName string) (bool, error) {
	return data.MetaDaoInst().DoesTableExist(tableName)
}
