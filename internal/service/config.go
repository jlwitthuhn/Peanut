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
	SetInt(tx *sql.Tx, key string, value int64) error
	SetString(tx *sql.Tx, key string, value string) error
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

func (this *configServiceImpl) SetInt(tx *sql.Tx, name string, value int64) error {
	err := this.configDao.UpsertIntByName(tx, name, value)
	return err
}

func (this *configServiceImpl) SetString(tx *sql.Tx, name string, value string) error {
	err := this.configDao.UpsertStringByName(tx, name, value)
	return err
}
