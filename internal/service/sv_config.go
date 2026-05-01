// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"peanut/internal/data"
)

type ConfigService interface {
	GetInt(ctx context.Context, key string) (int64, error)
	GetString(ctx context.Context, key string) (string, error)
	SetInt(ctx context.Context, key string, value int64) error
	SetString(ctx context.Context, key string, value string) error
}

func NewConfigService(configDao data.ConfigDao) ConfigService {
	return &configServiceImpl{configDao: configDao}
}

type configServiceImpl struct {
	configDao data.ConfigDao
}

func (this *configServiceImpl) GetInt(ctx context.Context, key string) (int64, error) {
	row, err := this.configDao.SelectIntRowByName(ctx, key)
	if err != nil {
		return 0, err
	}
	return row.Value, nil
}

func (this *configServiceImpl) GetString(ctx context.Context, key string) (string, error) {
	row, err := this.configDao.SelectStringRowByName(ctx, key)
	if err != nil {
		return "", err
	}
	return row.Value, nil
}

func (this *configServiceImpl) SetInt(ctx context.Context, name string, value int64) error {
	err := this.configDao.UpsertIntByName(ctx, name, value)
	return err
}

func (this *configServiceImpl) SetString(ctx context.Context, name string, value string) error {
	err := this.configDao.UpsertStringByName(ctx, name, value)
	return err
}
