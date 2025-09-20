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
	"peanut/internal/template"

	_ "github.com/lib/pq"
)

//go:embed static
var staticFs embed.FS

//go:embed template
var templateFs embed.FS

func main() {
	logger.Info("++ Starting Peanut ++")

	rawMux := http.NewServeMux()
	middlewareMux := http.NewServeMux()
	rootMux := http.NewServeMux()
	rootMux.Handle("/favicon.ico", rawMux)
	rootMux.Handle("/static/", rawMux)

	logger.Info("Preparing static files...")
	rawMux.Handle("/static/", http.FileServer(http.FS(staticFs)))

	logger.Info("Preparing templates...")
	justTemplates, err := fs.Sub(templateFs, "template")
	if err != nil {
		log.Fatal("Failed to find embedded template directory: ", err)
	}
	template.LoadTemplates(justTemplates)

	logger.Info("Registering routes...")
	pages.RegisterIndexHandlers(middlewareMux)
	pages.RegisterSetupHandlers(middlewareMux)
	wrappedMiddlewareMux := middleware.WrapHandler(middlewareMux, middleware.RequestLog(), middleware.RequestTimer(), middleware.DatabaseInitCheck())
	rootMux.Handle("/", wrappedMiddlewareMux)

	logger.Info("Connecting to postgres...")
	datasource.PostgresConnect()

	logger.Info("Startup complete, listening on :8080")
	logger.Fatal(http.ListenAndServe(":8080", rootMux))
}
