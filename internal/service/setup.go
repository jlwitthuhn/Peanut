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
	"peanut/internal/security/perms/permgroups"
	"time"
)

type SetupService interface {
	InitializeDatabase(r *http.Request, adminName string, adminEmail string, adminPlainPassword string) error
}

func NewSetupService(
	configDao data.ConfigDao,
	groupDao data.GroupDao,
	groupMembershipDao data.GroupMembershipDao,
	metaDao data.MetaDao,
	sessionDao data.SessionDao,
	userDao data.UserDao,
	configService ConfigService,
	groupService GroupService,
	userService UserService,
) SetupService {
	return &setupServiceImpl{
		configDao:          configDao,
		groupDao:           groupDao,
		groupMembershipDao: groupMembershipDao,
		metaDao:            metaDao,
		sessionDao:         sessionDao,
		userDao:            userDao,
		configService:      configService,
		groupService:       groupService,
		userService:        userService,
	}
}

type setupServiceImpl struct {
	configDao          data.ConfigDao
	groupDao           data.GroupDao
	groupMembershipDao data.GroupMembershipDao
	metaDao            data.MetaDao
	sessionDao         data.SessionDao
	userDao            data.UserDao
	configService      ConfigService
	groupService       GroupService
	userService        UserService
}

func (this *setupServiceImpl) InitializeDatabase(r *http.Request, adminName string, adminEmail string, adminPlainPassword string) error {
	logger.Debug(r, "Preparing transaction...")
	ctx := context.Background()
	tx, txErr := datasource.PostgresHandle().BeginTx(ctx, nil)
	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	logger.Debug(r, "Creating tables...")
	{
		metaErr := this.metaDao.CreateDBObjects(tx)
		if metaErr != nil {
			return metaErr
		}
		configErr := this.configDao.CreateDBObjects(tx)
		if configErr != nil {
			return configErr
		}
		groupErr := this.groupDao.CreateDBObjects(tx)
		if groupErr != nil {
			return groupErr
		}
		userErr := this.userDao.CreateDBObjects(tx)
		if userErr != nil {
			return userErr
		}
		groupMembershipErr := this.groupMembershipDao.CreateDBObjects(tx)
		if groupMembershipErr != nil {
			return groupMembershipErr
		}
		sessionErr := this.sessionDao.CreateDBObjects(tx)
		if sessionErr != nil {
			return sessionErr
		}
	}

	logger.Debug(r, "Populating data...")
	configTimeErr := this.configService.SetInt(configkey.IntInitializedTime, time.Now().Unix(), tx)
	if configTimeErr != nil {
		return configTimeErr
	}
	{
		groupTurboErr := this.groupService.CreateGroup(tx, permgroups.TurboAdmin, "Full control over everything.", true)
		if groupTurboErr != nil {
			return groupTurboErr
		}
		groupAdminErr := this.groupService.CreateGroup(tx, permgroups.Admin, "Full control over everything except mass database updates and exports.", true)
		if groupAdminErr != nil {
			return groupAdminErr
		}
		groupUserErr := this.groupService.CreateGroup(tx, permgroups.User, "Ordinary registered user.", true)
		if groupUserErr != nil {
			return groupUserErr
		}
	}
	userId, userErr := this.userService.CreateUser(tx, adminName, adminEmail, adminPlainPassword)
	if userErr != nil {
		return userErr
	}
	{
		memberTurboErr := this.groupService.EnrollUserInGroup(r, tx, userId, permgroups.TurboAdmin)
		if memberTurboErr != nil {
			return memberTurboErr
		}
		memberAdminErr := this.groupService.EnrollUserInGroup(r, tx, userId, permgroups.Admin)
		if memberAdminErr != nil {
			return memberAdminErr
		}
		memberUserErr := this.groupService.EnrollUserInGroup(r, tx, userId, permgroups.User)
		if memberUserErr != nil {
			return memberUserErr
		}
	}

	logger.Debug(r, "Commiting transaction...")
	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}

	return nil
}
