// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"net/http"
	"peanut/internal/data"
	"peanut/internal/data/configkey"
	"peanut/internal/data/datasource"
	"peanut/internal/logger"
	"time"
)

type SetupService interface {
	InitializeDatabase(r *http.Request, adminName string, adminEmail string, adminPlainPassword string) error
}

func NewSetupService(configService ConfigService, groupService GroupService, userService UserService) SetupService {
	return &setupServiceImpl{configService: configService, groupService: groupService, userService: userService}
}

type setupServiceImpl struct {
	configService ConfigService
	groupService  GroupService
	userService   UserService
}

func (this *setupServiceImpl) InitializeDatabase(r *http.Request, adminName string, adminEmail string, adminPlainPassword string) error {
	logger.Trace(r, "Preparing transaction...")
	ctx := context.Background()
	tx, txErr := datasource.PostgresHandle().BeginTx(ctx, nil)
	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	logger.Trace(r, "Creating tables...")
	{
		metaErr := data.MetaDaoInst().CreateDBObjects(tx)
		if metaErr != nil {
			return metaErr
		}
		configErr := data.ConfigDaoInst().CreateDBObjects(tx)
		if configErr != nil {
			return configErr
		}
		groupErr := data.GroupDaoInst().CreateDBObjects(tx)
		if groupErr != nil {
			return groupErr
		}
		userErr := data.UserDaoInst().CreateDBObjects(tx)
		if userErr != nil {
			return userErr
		}
		sessionErr := data.SessionDaoInst().CreateDBObjects(tx)
		if sessionErr != nil {
			return sessionErr
		}
	}

	logger.Trace(r, "Populating data...")
	configTimeErr := this.configService.SetInt(configkey.IntInitializedTime, time.Now().Unix(), tx)
	if configTimeErr != nil {
		return configTimeErr
	}
	{
		groupTurboErr := this.groupService.CreateGroup(tx, "TurboAdmin", "Full control over everything.", true)
		if groupTurboErr != nil {
			return groupTurboErr
		}
		groupAdminErr := this.groupService.CreateGroup(tx, "Admin", "Full control over everything except mass database updates and exports.", true)
		if groupAdminErr != nil {
			return groupAdminErr
		}
		groupUserErr := this.groupService.CreateGroup(tx, "User", "Ordinary registered user.", true)
		if groupUserErr != nil {
			return groupUserErr
		}
	}
	userErr := this.userService.CreateUser(tx, adminName, adminEmail, adminPlainPassword)
	if userErr != nil {
		return userErr
	}

	logger.Trace(r, "Commiting transaction...")
	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}

	return nil
}
