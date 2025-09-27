// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"database/sql"
	"peanut/internal/data"
)

type ConfigService interface {
	GetInt(tx *sql.Tx, key string) (int64, error)
	SetInt(key string, value int64, tx *sql.Tx) error
}

func NewConfigService(configDao data.ConfigDao) ConfigService {
	return &configServiceImpl{configDao: configDao}
}

type configServiceImpl struct {
	configDao data.ConfigDao
}

func (this *configServiceImpl) GetInt(tx *sql.Tx, key string) (int64, error) {
	row, err := this.configDao.SelectIntRowByName(tx, key)
	if err != nil {
		return 0, err
	}
	return row.Value, nil
}

func (this *configServiceImpl) SetInt(name string, value int64, tx *sql.Tx) error {
	err := this.configDao.UpsertIntByName(name, value, tx)
	return err
}
