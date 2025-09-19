// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package database

import (
	"database/sql"
	"fmt"
	"os"

	"peanut/internal/logger"
)

var pgDb *sql.DB = nil

func PostgresConnect() {
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
		logger.Warn("Database did not respond to ping:", pingErr)
	}

	pgDb = db
}

func PostgresHandle() *sql.DB {
	return pgDb
}
