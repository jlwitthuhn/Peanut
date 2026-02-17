// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"errors"
	"net/http"
	"peanut/internal/data"
	"peanut/internal/data/configkey"
	"peanut/internal/logger"
	"peanut/internal/security/perms/permgroups"
	"strconv"
	"strings"
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
	scheduledJobDao data.ScheduledJobDao,
	scheduledJobRunDao data.ScheduledJobRunDao,
	sessionDao data.SessionDao,
	sessionStringDao data.SessionStringDao,
	userDao data.UserDao,
	configService ConfigService,
	databaseService DatabaseService,
	groupService GroupService,
	scheduledJobService ScheduledJobService,
	userService UserService,
) SetupService {
	return &setupServiceImpl{
		configDao:           configDao,
		databaseService:     databaseService,
		groupDao:            groupDao,
		groupMembershipDao:  groupMembershipDao,
		metaDao:             metaDao,
		scheduledJobDao:     scheduledJobDao,
		scheduledJobRunDao:  scheduledJobRunDao,
		sessionDao:          sessionDao,
		sessionStringDao:    sessionStringDao,
		userDao:             userDao,
		configService:       configService,
		groupService:        groupService,
		scheduledJobService: scheduledJobService,
		userService:         userService,
	}
}

type setupServiceImpl struct {
	configDao           data.ConfigDao
	groupDao            data.GroupDao
	groupMembershipDao  data.GroupMembershipDao
	metaDao             data.MetaDao
	scheduledJobDao     data.ScheduledJobDao
	scheduledJobRunDao  data.ScheduledJobRunDao
	sessionDao          data.SessionDao
	sessionStringDao    data.SessionStringDao
	userDao             data.UserDao
	configService       ConfigService
	databaseService     DatabaseService
	groupService        GroupService
	scheduledJobService ScheduledJobService
	userService         UserService
}

func (this *setupServiceImpl) InitializeDatabase(r *http.Request, adminName string, adminEmail string, adminPlainPassword string) error {
	logger.Info(r, "Database initialization starting")

	logger.Debug(r, "Checking postgres version...")
	{
		pgVersion, err := this.databaseService.GetPostgresVersion(r)
		if err != nil {
			return err
		}
		dotIndex := strings.IndexRune(pgVersion, '.')
		if dotIndex == -1 {
			return errors.New("failed to parse postgres version string")
		}
		versionInt, err := strconv.Atoi(pgVersion[:dotIndex])
		if err != nil {
			return errors.New("failed to parse postgres version number")
		}
		if versionInt < 18 {
			return errors.New("only postgres version 18.0 and higher is supported")
		}
	}

	logger.Debug(r, "Creating tables...")

	err := this.metaDao.CreateDBObjects(r)
	if err != nil {
		return err
	}
	err = this.configDao.CreateDBObjects(r)
	if err != nil {
		return err
	}
	err = this.groupDao.CreateDBObjects(r)
	if err != nil {
		return err
	}
	err = this.scheduledJobDao.CreateDBObjects(r)
	if err != nil {
		return err
	}
	err = this.scheduledJobRunDao.CreateDBObjects(r)
	if err != nil {
		return err
	}
	err = this.userDao.CreateDBObjects(r)
	if err != nil {
		return err
	}
	err = this.groupMembershipDao.CreateDBObjects(r)
	if err != nil {
		return err
	}
	err = this.sessionDao.CreateDBObjects(r)
	if err != nil {
		return err
	}
	err = this.sessionStringDao.CreateDBObjects(r)
	if err != nil {
		return err
	}

	logger.Debug(r, "Populating data...")

	err = this.scheduledJobService.AddJobDefinition(r, "DeleteExpiredSessions", time.Hour)
	if err != nil {
		return err
	}
	err = this.scheduledJobService.AddJobDefinition(r, "VacuumDatabase", 6*time.Hour)
	if err != nil {
		return err
	}

	err = this.configService.SetInt(r, configkey.IntInitializedTime, time.Now().Unix())
	if err != nil {
		return err
	}
	err = this.configService.SetString(r, configkey.StringWelcomeMessage, "Haldo.")
	if err != nil {
		return err
	}

	err = this.groupService.CreateGroup(r, permgroups.TurboAdmin, "Full control over everything.", true)
	if err != nil {
		return err
	}
	err = this.groupService.CreateGroup(r, permgroups.Admin, "Much control over most things.", true)
	if err != nil {
		return err
	}
	err = this.groupService.CreateGroup(r, permgroups.User, "Ordinary registered user.", true)
	if err != nil {
		return err
	}

	logger.Debug(r, "Creating admin user...")

	userId, err := this.userService.CreateUser(r, adminName, adminEmail, adminPlainPassword)
	if err != nil {
		return err
	}
	err = this.groupService.EnrollUserInGroup(r, userId, permgroups.TurboAdmin)
	if err != nil {
		return err
	}
	err = this.groupService.EnrollUserInGroup(r, userId, permgroups.Admin)
	if err != nil {
		return err
	}
	err = this.groupService.EnrollUserInGroup(r, userId, permgroups.User)
	if err != nil {
		return err
	}

	logger.Info(r, "Database initialization succeeded")

	return nil
}
