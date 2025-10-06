// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"net/http"
	"peanut/internal/data"
)

type ConfigService interface {
	GetInt(req *http.Request, key string) (int64, error)
	GetString(req *http.Request, key string) (string, error)
	SetInt(req *http.Request, key string, value int64) error
	SetString(req *http.Request, key string, value string) error
}

func NewConfigService(configDao data.ConfigDao) ConfigService {
	return &configServiceImpl{configDao: configDao}
}

type configServiceImpl struct {
	configDao data.ConfigDao
}

func (this *configServiceImpl) GetInt(req *http.Request, key string) (int64, error) {
	row, err := this.configDao.SelectIntRowByName(req, key)
	if err != nil {
		return 0, err
	}
	return row.Value, nil
}

func (this *configServiceImpl) GetString(req *http.Request, key string) (string, error) {
	row, err := this.configDao.SelectStringRowByName(req, key)
	if err != nil {
		return "", err
	}
	return row.Value, nil
}

func (this *configServiceImpl) SetInt(req *http.Request, name string, value int64) error {
	err := this.configDao.UpsertIntByName(req, name, value)
	return err
}

func (this *configServiceImpl) SetString(req *http.Request, name string, value string) error {
	err := this.configDao.UpsertStringByName(req, name, value)
	return err
}
