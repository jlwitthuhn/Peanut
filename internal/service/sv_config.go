// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"fmt"
	"peanut/internal/data"
	"peanut/internal/keynames/contextkeys"
	"sync"
	"time"
)

const cacheTTL = 5 * time.Second

type cachedInt struct {
	value    int64
	readTime time.Time
}

type cachedString struct {
	value    string
	readTime time.Time
}

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
	return &configServiceImpl{
		configDao:    configDao,
		systemLogDao: systemLogDao,
		intCache:     make(map[string]cachedInt),
		stringCache:  make(map[string]cachedString),
	}
}

func NewSetupConfigService(configDao data.ConfigDao, systemLogDao data.SystemLogDao) SetupConfigService {
	return &configServiceImpl{
		configDao:    configDao,
		systemLogDao: systemLogDao,
		intCache:     make(map[string]cachedInt),
		stringCache:  make(map[string]cachedString),
	}
}

type configServiceImpl struct {
	configDao        data.ConfigDao
	systemLogDao     data.SystemLogDao
	intCacheMutex    sync.Mutex
	intCache         map[string]cachedInt
	stringCacheMutex sync.Mutex
	stringCache      map[string]cachedString
}

func (this *configServiceImpl) GetInt(ctx context.Context, key string) (int64, error) {
	this.intCacheMutex.Lock()
	defer this.intCacheMutex.Unlock()
	if cached, ok := this.intCache[key]; ok && time.Since(cached.readTime) < cacheTTL {
		return cached.value, nil
	}
	row, err := this.configDao.SelectIntRowByName(ctx, key)
	if err != nil {
		return 0, err
	}
	this.intCache[key] = cachedInt{value: row.Value, readTime: time.Now()}
	return row.Value, nil
}

func (this *configServiceImpl) GetString(ctx context.Context, key string) (string, error) {
	this.stringCacheMutex.Lock()
	defer this.stringCacheMutex.Unlock()
	if cached, ok := this.stringCache[key]; ok && time.Since(cached.readTime) < cacheTTL {
		return cached.value, nil
	}
	row, err := this.configDao.SelectStringRowByName(ctx, key)
	if err != nil {
		return "", err
	}
	this.stringCache[key] = cachedString{value: row.Value, readTime: time.Now()}
	return row.Value, nil
}

func (this *configServiceImpl) clearIntCache(name string) {
	this.intCacheMutex.Lock()
	defer this.intCacheMutex.Unlock()
	delete(this.intCache, name)
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
	this.clearIntCache(name)
	message := fmt.Sprintf("Set config int '%s' to: %d", name, value)
	return this.systemLogDao.InsertRow(ctx, userId, "", message)
}

func (this *configServiceImpl) clearStringCache(name string) {
	this.stringCacheMutex.Lock()
	defer this.stringCacheMutex.Unlock()
	delete(this.stringCache, name)
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
	this.clearStringCache(name)
	message := fmt.Sprintf("Set config string '%s' to: %s", name, value)
	return this.systemLogDao.InsertRow(ctx, userId, "", message)
}

func (this *configServiceImpl) SetIntSetup(ctx context.Context, name string, value int64) error {
	return this.configDao.UpsertIntByName(ctx, name, value)
}

func (this *configServiceImpl) SetStringSetup(ctx context.Context, name string, value string) error {
	return this.configDao.UpsertStringByName(ctx, name, value)
}
