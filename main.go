// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

//go:embed static
var staticFs embed.FS

func connectDb() {
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
		log.Fatal("Error connecting to database: ", err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal("Error pinging database: ", pingErr)
	}
	defer db.Close()
}

func main() {
	log.Println("Preparing static files...")
	http.Handle("/static/", http.FileServer(http.FS(staticFs)))

	log.Println("Connecting to database...")
	connectDb()

	http.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, "Welcome to Peanut")
		if err != nil {
			log.Println("Error:", err)
		}
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
