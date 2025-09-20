// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"database/sql"
	"peanut/internal/data"
)

type ConfigService interface {
	SetInt(key string, value int64, tx *sql.Tx) error
}

func NewConfigService() ConfigService {
	return &configServiceImpl{}
}

type configServiceImpl struct{}

func (*configServiceImpl) SetInt(name string, value int64, tx *sql.Tx) error {
	err := data.ConfigDaoInst().UpsertIntByName(name, value, tx)
	return err
}
