// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"fmt"
	"peanut/internal/data"
	"peanut/internal/keynames/contextkeys"
)

type ConfigService interface {
	GetInt(ctx context.Context, key string) (int64, error)
	GetString(ctx context.Context, key string) (string, error)
	SetInt(ctx context.Context, key string, value int64) error
	SetString(ctx context.Context, key string, value string) error
}

type SetupConfigService interface {
	SetIntSetup(ctx context.Context, key string, value int64) error
	SetStringSetup(ctx context.Context, key string, value string) error
}

func NewConfigService(configDao data.ConfigDao, systemLogDao data.SystemLogDao) ConfigService {
	return &configServiceImpl{configDao: configDao, systemLogDao: systemLogDao}
}

func NewSetupConfigService(configDao data.ConfigDao, systemLogDao data.SystemLogDao) SetupConfigService {
	return &configServiceImpl{configDao: configDao, systemLogDao: systemLogDao}
}

type configServiceImpl struct {
	configDao    data.ConfigDao
	systemLogDao data.SystemLogDao
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
	userId, ok := ctx.Value(contextkeys.UserId).(string)
	if !ok {
		return fmt.Errorf("cannot set config int '%s': no user id in context", name)
	}
	err := this.configDao.UpsertIntByName(ctx, name, value)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("Set config int '%s' to: %d", name, value)
	return this.systemLogDao.InsertRow(ctx, userId, message)
}

func (this *configServiceImpl) SetString(ctx context.Context, name string, value string) error {
	userId, ok := ctx.Value(contextkeys.UserId).(string)
	if !ok {
		return fmt.Errorf("cannot set config string '%s': no user id in context", name)
	}
	err := this.configDao.UpsertStringByName(ctx, name, value)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("Set config string '%s' to: %s", name, value)
	return this.systemLogDao.InsertRow(ctx, userId, message)
}

func (this *configServiceImpl) SetIntSetup(ctx context.Context, name string, value int64) error {
	return this.configDao.UpsertIntByName(ctx, name, value)
}

func (this *configServiceImpl) SetStringSetup(ctx context.Context, name string, value string) error {
	return this.configDao.UpsertStringByName(ctx, name, value)
}
