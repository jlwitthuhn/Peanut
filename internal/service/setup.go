// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"peanut/internal/data"
	"peanut/internal/data/configkey"
	"peanut/internal/data/datasource"
	"peanut/internal/logger"
	"time"
)

type SetupService interface {
	InitializeDatabase(context.Context) error
}

func NewSetupService(configService ConfigService) SetupService {
	return &setupServiceImpl{configService: configService}
}

type setupServiceImpl struct {
	configService ConfigService
}

func (this *setupServiceImpl) InitializeDatabase(ctx context.Context) error {
	logger.Trace("Preparing transaction...")
	tx, txErr := datasource.PostgresHandle().BeginTx(ctx, nil)
	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	logger.Trace("Creating tables...")
	metaErr := data.MetaDaoInst().CreateDBObjects(tx)
	if metaErr != nil {
		return metaErr
	}
	configErr := data.ConfigDaoInst().CreateDBObjects(tx)
	if configErr != nil {
		return configErr
	}

	logger.Trace("Populating data...")
	configTimeErr := this.configService.SetInt(configkey.IntInitializedTime, time.Now().Unix(), tx)
	if configTimeErr != nil {
		return configTimeErr
	}

	logger.Trace("Commiting transaction...")
	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}

	return nil
}
