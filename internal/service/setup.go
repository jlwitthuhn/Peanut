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
	sessionStringDao data.SessionStringDao,
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
		sessionStringDao:   sessionStringDao,
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
	sessionStringDao   data.SessionStringDao
	userDao            data.UserDao
	configService      ConfigService
	groupService       GroupService
	userService        UserService
}

func (this *setupServiceImpl) InitializeDatabase(r *http.Request, adminName string, adminEmail string, adminPlainPassword string) error {
	logger.Debug(r, "Preparing transaction...")
	ctx := context.Background()
	tx, err := datasource.PostgresHandle().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	logger.Debug(r, "Creating tables...")
	{
		err := this.metaDao.CreateDBObjects(tx)
		if err != nil {
			return err
		}
		err = this.configDao.CreateDBObjects(tx)
		if err != nil {
			return err
		}
		err = this.groupDao.CreateDBObjects(tx)
		if err != nil {
			return err
		}
		err = this.userDao.CreateDBObjects(tx)
		if err != nil {
			return err
		}
		err = this.groupMembershipDao.CreateDBObjects(tx)
		if err != nil {
			return err
		}
		err = this.sessionDao.CreateDBObjects(tx)
		if err != nil {
			return err
		}
		err = this.sessionStringDao.CreateDBObjects(tx)
		if err != nil {
			return err
		}
	}

	logger.Debug(r, "Populating data...")
	err = this.configService.SetInt(tx, configkey.IntInitializedTime, time.Now().Unix())
	if err != nil {
		return err
	}
	err = this.configService.SetString(tx, configkey.StringWelcomeMessage, "Haldo.")
	if err != nil {
		return err
	}
	{
		err := this.groupService.CreateGroup(tx, permgroups.TurboAdmin, "Full control over everything.", true)
		if err != nil {
			return err
		}
		err = this.groupService.CreateGroup(tx, permgroups.Admin, "Full control over everything except mass database updates and exports.", true)
		if err != nil {
			return err
		}
		err = this.groupService.CreateGroup(tx, permgroups.User, "Ordinary registered user.", true)
		if err != nil {
			return err
		}
	}
	userId, err := this.userService.CreateUser(tx, adminName, adminEmail, adminPlainPassword)
	if err != nil {
		return err
	}
	{
		err := this.groupService.EnrollUserInGroup(r, tx, userId, permgroups.TurboAdmin)
		if err != nil {
			return err
		}
		err = this.groupService.EnrollUserInGroup(r, tx, userId, permgroups.Admin)
		if err != nil {
			return err
		}
		err = this.groupService.EnrollUserInGroup(r, tx, userId, permgroups.User)
		if err != nil {
			return err
		}
	}

	logger.Debug(r, "Commiting transaction...")
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
