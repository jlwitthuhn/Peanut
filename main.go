// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"embed"
	"io/fs"
	"net/http"
	"peanut/internal/data"
	"peanut/internal/data/datasource"
	"peanut/internal/endpoints"
	"peanut/internal/endpoints/ep_admin"
	"peanut/internal/logger"
	"peanut/internal/middleware"
	"peanut/internal/service"
	"peanut/internal/template"

	_ "github.com/lib/pq"
)

//go:embed static
var staticFs embed.FS

//go:embed template
var templateFs embed.FS

func main() {
	logger.Info(nil, "++ Starting Peanut ++")

	rawMux := http.NewServeMux()
	middlewareMux := http.NewServeMux()
	rootMux := http.NewServeMux()
	rootMux.Handle("/favicon.ico", rawMux)
	rootMux.Handle("/static/", rawMux)

	logger.Info(nil, "Preparing static files...")
	rawMux.Handle("/static/", http.FileServer(http.FS(staticFs)))

	logger.Info(nil, "Preparing templates...")
	justTemplates, err := fs.Sub(templateFs, "template")
	if err != nil {
		logger.Fatal(nil, "Failed to find embedded template directory: ", err)
	}
	template.LoadTemplates(justTemplates)

	logger.Info(nil, "Initializing services...")

	configDao := data.NewConfigDao()
	groupDao := data.NewGroupDao()
	groupMembershipDao := data.NewGroupMembershipDao()
	metaDao := data.NewMetaDao()
	multiTableDao := data.NewMultiTableDao()
	scheduledJobDao := data.NewScheduledJobDao()
	sessionDao := data.NewSessionDao()
	sessionStringDao := data.NewSessionStringDao()
	userDao := data.NewUserDao()

	configService := service.NewConfigService(configDao)
	dbService := service.NewDatabaseService(metaDao)
	groupService := service.NewGroupService(groupDao, groupMembershipDao, multiTableDao)
	sessionService := service.NewSessionService(sessionDao, sessionStringDao, userDao)
	userService := service.NewUserService(sessionDao, userDao)
	setupService := service.NewSetupService(
		configDao, groupDao, groupMembershipDao, metaDao, scheduledJobDao, sessionDao,
		sessionStringDao, userDao, configService, groupService, userService,
	)

	// Setup mux is separate and is only used from within DatabaseInitCheck
	setupMux := http.NewServeMux()
	endpoints.RegisterSetupHandlers(setupMux, dbService, setupService)

	logger.Info(nil, "Registering routes...")
	ep_admin.RegisterAdminHandlers(middlewareMux, configService, dbService, groupService, sessionService, userService)
	endpoints.RegisterIndexHandlers(middlewareMux, configService)
	endpoints.RegisterLoginHandlers(middlewareMux, sessionService)
	endpoints.RegisterLogoutHandlers(middlewareMux, sessionService)
	endpoints.RegisterRegisterHandlers(middlewareMux, groupService, userService)
	wrappedMiddlewareMux := middleware.WrapHandler(middlewareMux,
		middleware.RequestId(),
		middleware.RequestLog(),
		middleware.RequestTimer(),
		middleware.SecurityHeaders(),
		middleware.PostgresTransaction(),
		middleware.DatabaseInitCheck(dbService, setupMux),
		middleware.Authentication(groupService, sessionService),
		middleware.CsrfProtection(sessionService),
	)
	rootMux.Handle("/", wrappedMiddlewareMux)

	logger.Info(nil, "Connecting to postgres...")
	datasource.PostgresConnect()

	logger.Info(nil, "Startup complete, listening on :8080")
	logger.Fatal(nil, http.ListenAndServe(":8080", rootMux))
}
