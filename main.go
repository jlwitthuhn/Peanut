// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"peanut/internal/logger"
	"peanut/internal/template"
)

//go:embed static
var staticFs embed.FS

//go:embed template
var templateFs embed.FS

func connectDb() *sql.DB {
	var host string
	if os.Getenv("PEANUT_DB_HOST") != "" {
		host = os.Getenv("PEANUT_DB_HOST")
	} else {
		host = "localhost"
	}
	var user string
	if os.Getenv("PEANUT_DB_USER") != "" {
		user = os.Getenv("PEANUT_DB_USER")
	} else {
		user = "peanut"
	}
	var password string
	if os.Getenv("PEANUT_DB_PASSWORD") != "" {
		password = os.Getenv("PEANUT_DB_PASSWORD")
	} else {
		password = "password"
	}
	var dbname string
	if os.Getenv("PEANUT_DB_NAME") != "" {
		dbname = os.Getenv("PEANUT_DB_NAME")
	} else {
		dbname = "peanut"
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal("Error connecting to database:", err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		logger.Warn("Error pinging database:", pingErr)
	}

	return db
}

func main() {
	logger.Info("++ Starting Peanut ++")

	logger.Info("Preparing static files...")
	http.Handle("/static/", http.FileServer(http.FS(staticFs)))

	logger.Info("Preparing templates...")
	justTemplates, err := fs.Sub(templateFs, "template")
	if err != nil {
		log.Fatal("Failed to find embedded template directory: ", err)
	}
	template.LoadTemplates(justTemplates)

	logger.Info("Registering routes...")
	http.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		theTemplate := template.GetTemplate("_index")
		err = theTemplate.Execute(w, nil)
		if err != nil {
			logger.Error("Error executing template:", err)
		}
	})

	logger.Info("Connecting to database...")
	connectDb()

	logger.Info("Startup complete, listening on :8080")
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
