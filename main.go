// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
)

//go:embed static
var staticFs embed.FS

func main() {
	log.Println("Preparing static files...")
	http.Handle("/static/", http.FileServer(http.FS(staticFs)))

	http.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, "Welcome to Peanut")
		if err != nil {
			log.Println("Error:", err)
		}
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
