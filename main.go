// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"peanut/internal/data/datasource"
	"peanut/internal/logger"
	"peanut/internal/middleware"
	"peanut/internal/pages"
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
		log.Fatal("Failed to find embedded template directory: ", err)
	}
	template.LoadTemplates(justTemplates)

	logger.Info(nil, "Initializing services...")
	var configService = service.NewConfigService()
	var dbService = service.NewDatabaseService()
	var groupService = service.NewGroupService()
	var userService = service.NewUserService()
	var setupService = service.NewSetupService(configService, groupService, userService)

	// Setup mux is separate and is only used from within DatabaseInitCheck
	setupMux := http.NewServeMux()
	pages.RegisterSetupHandlers(setupMux, dbService, setupService)

	logger.Info(nil, "Registering routes...")
	pages.RegisterAdminHandlers(middlewareMux)
	pages.RegisterIndexHandlers(middlewareMux)
	pages.RegisterLoginHandlers(middlewareMux, userService)
	pages.RegisterLogoutHandlers(middlewareMux, userService)
	pages.RegisterRegisterHandlers(middlewareMux, groupService, userService)
	wrappedMiddlewareMux := middleware.WrapHandler(middlewareMux,
		middleware.RequestId(),
		middleware.RequestLog(),
		middleware.RequestTimer(),
		middleware.DatabaseInitCheck(dbService, setupMux),
		middleware.Authentication(groupService, userService),
		middleware.SecurityHeaders(),
	)
	rootMux.Handle("/", wrappedMiddlewareMux)

	logger.Info(nil, "Connecting to postgres...")
	datasource.PostgresConnect()

	logger.Info(nil, "Startup complete, listening on :8080")
	logger.Fatal(nil, http.ListenAndServe(":8080", rootMux))
}
