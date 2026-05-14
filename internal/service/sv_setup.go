// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"peanut/internal/data"
	"peanut/internal/keynames/configkey"
	"peanut/internal/logger"
	"peanut/internal/security/perms/permgroups"
	"strconv"
	"strings"
	"time"
)

type SetupService interface {
	InitializeDatabase(ctx context.Context, adminName string, adminEmail string, adminPlainPassword string) error
}

func NewSetupService(
	configDao data.ConfigDao,
	forumsDao data.ForumsDao,
	forumPostsDao data.ForumPostsDao,
	forumSectionsDao data.ForumSectionsDao,
	forumThreadsDao data.ForumThreadsDao,
	groupDao data.GroupDao,
	groupMembershipDao data.GroupMembershipDao,
	metaDao data.MetaDao,
	scheduledJobDao data.ScheduledJobDao,
	scheduledJobRunDao data.ScheduledJobRunDao,
	sessionDao data.SessionDao,
	sessionStringDao data.SessionStringDao,
	systemLogDao data.SystemLogDao,
	userDao data.UserDao,
	forumHomeDvao data.ForumHomeDvao,
	forumPostListingDvao data.ForumPostListingDvao,
	forumThreadSummaryDvao data.ForumThreadSummaryDvao,
	setupConfigService SetupConfigService,
	databaseService DatabaseService,
	groupService GroupService,
	scheduledJobService ScheduledJobService,
	userService UserService,
) SetupService {
	return &setupServiceImpl{
		configDao:              configDao,
		forumsDao:              forumsDao,
		forumPostsDao:          forumPostsDao,
		forumSectionsDao:       forumSectionsDao,
		forumThreadsDao:        forumThreadsDao,
		databaseService:        databaseService,
		groupDao:               groupDao,
		groupMembershipDao:     groupMembershipDao,
		metaDao:                metaDao,
		scheduledJobDao:        scheduledJobDao,
		scheduledJobRunDao:     scheduledJobRunDao,
		sessionDao:             sessionDao,
		sessionStringDao:       sessionStringDao,
		systemLogDao:           systemLogDao,
		userDao:                userDao,
		forumHomeDvao:          forumHomeDvao,
		forumPostListingDvao:   forumPostListingDvao,
		forumThreadSummaryDvao: forumThreadSummaryDvao,
		setupConfigService:     setupConfigService,
		groupService:           groupService,
		scheduledJobService:    scheduledJobService,
		userService:            userService,
	}
}

type setupServiceImpl struct {
	configDao              data.ConfigDao
	forumsDao              data.ForumsDao
	forumPostsDao          data.ForumPostsDao
	forumSectionsDao       data.ForumSectionsDao
	forumThreadsDao        data.ForumThreadsDao
	groupDao               data.GroupDao
	groupMembershipDao     data.GroupMembershipDao
	metaDao                data.MetaDao
	scheduledJobDao        data.ScheduledJobDao
	scheduledJobRunDao     data.ScheduledJobRunDao
	sessionDao             data.SessionDao
	sessionStringDao       data.SessionStringDao
	systemLogDao           data.SystemLogDao
	userDao                data.UserDao
	forumHomeDvao          data.ForumHomeDvao
	forumPostListingDvao   data.ForumPostListingDvao
	forumThreadSummaryDvao data.ForumThreadSummaryDvao
	setupConfigService     SetupConfigService
	databaseService        DatabaseService
	groupService           GroupService
	scheduledJobService    ScheduledJobService
	userService            UserService
}

func (this *setupServiceImpl) InitializeDatabase(ctx context.Context, adminName string, adminEmail string, adminPlainPassword string) error {
	logger.Info(ctx, "Database initialization starting")

	logger.Debug(ctx, "Checking postgres version...")
	{
		pgVersion, err := this.databaseService.GetPostgresVersion(ctx)
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

	logger.Debug(ctx, "Creating tables...")

	// Core functionality
	err := this.metaDao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}
	err = this.configDao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}
	err = this.groupDao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}
	err = this.scheduledJobDao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}
	err = this.scheduledJobRunDao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}
	err = this.userDao.CreateDBObjects(ctx)
	// Below depend on `users` existing
	if err != nil {
		return err
	}
	err = this.groupMembershipDao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}
	err = this.sessionDao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}
	err = this.sessionStringDao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}
	err = this.systemLogDao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}

	// Forums
	err = this.forumSectionsDao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}
	err = this.forumsDao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}
	err = this.forumThreadsDao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}
	err = this.forumPostsDao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}
	err = this.forumHomeDvao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}
	err = this.forumPostListingDvao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}
	err = this.forumThreadSummaryDvao.CreateDBObjects(ctx)
	if err != nil {
		return err
	}

	logger.Debug(ctx, "Populating data...")

	err = this.scheduledJobService.AddJobDefinition(ctx, "DeleteExpiredSessions", time.Hour)
	if err != nil {
		return err
	}
	err = this.scheduledJobService.AddJobDefinition(ctx, "VacuumDatabase", 6*time.Hour)
	if err != nil {
		return err
	}

	err = this.setupConfigService.SetIntSetup(ctx, configkey.IntInitializedTime, time.Now().Unix())
	if err != nil {
		return err
	}
	err = this.setupConfigService.SetIntSetup(ctx, configkey.IntSessionLengthMinutes, 120)
	if err != nil {
		return err
	}
	err = this.setupConfigService.SetStringSetup(ctx, configkey.StringWelcomeMessage, "Haldo.")
	if err != nil {
		return err
	}

	err = this.groupService.CreateGroup(ctx, permgroups.TurboAdmin, "Full control over everything.", true)
	if err != nil {
		return err
	}
	err = this.groupService.CreateGroup(ctx, permgroups.Admin, "Much control over most things.", true)
	if err != nil {
		return err
	}
	err = this.groupService.CreateGroup(ctx, permgroups.User, "Ordinary registered user.", true)
	if err != nil {
		return err
	}

	logger.Debug(ctx, "Creating admin user...")

	userId, err := this.userService.CreateUser(ctx, adminName, adminEmail, adminPlainPassword)
	if err != nil {
		return err
	}
	err = this.groupService.EnrollUserInGroup(ctx, userId, permgroups.TurboAdmin)
	if err != nil {
		return err
	}
	err = this.groupService.EnrollUserInGroup(ctx, userId, permgroups.Admin)
	if err != nil {
		return err
	}
	err = this.groupService.EnrollUserInGroup(ctx, userId, permgroups.User)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Database initialization succeeded")

	return nil
}
