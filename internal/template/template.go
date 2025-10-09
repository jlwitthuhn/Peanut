// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package template

import (
	"html/template"
	"io/fs"
	"log"

	"peanut/internal/logger"
)

var templatesByName map[string]*template.Template = make(map[string]*template.Template)

func GetTemplate(name string) *template.Template {
	return templatesByName[name]
}

// LoadTemplates builds the list of view name to template list mappings.
// Because of the way go manages templates, this needs to be kept separately from template content.
func LoadTemplates(fs fs.FS) {

	adminIndexFiles := []string{"base.html", "css/common.css", "admin/base.html", "admin/index.html"}
	loadTemplateOrDie(fs, "_admin/index", adminIndexFiles...)

	adminFrontPageFiles := []string{"base.html", "css/common.css", "admin/base.html", "admin/front_page.html"}
	loadTemplateOrDie(fs, "_admin/front_page", adminFrontPageFiles...)

	adminSimpleMessageFiles := []string{"base.html", "css/common.css", "admin/base.html", "admin/simple_message.html"}
	loadTemplateOrDie(fs, "_admin/simple_message", adminSimpleMessageFiles...)

	indexFiles := []string{"base.html", "css/common.css", "index.html"}
	loadTemplateOrDie(fs, "_index", indexFiles...)

	loginFiles := []string{"base.html", "css/common.css", "login.html"}
	loadTemplateOrDie(fs, "_login", loginFiles...)

	registerFiles := []string{"base.html", "css/common.css", "register.html"}
	loadTemplateOrDie(fs, "_register", registerFiles...)

	setupFiles := []string{"base.html", "css/common.css", "setup.html"}
	loadTemplateOrDie(fs, "_setup", setupFiles...)

	simpleMessageFiles := []string{"base.html", "css/common.css", "simple_message.html"}
	loadTemplateOrDie(fs, "_simple_message", simpleMessageFiles...)
}

func loadTemplateOrDie(fs fs.FS, name string, files ...string) {
	_, exists := templatesByName[name]
	if exists {
		logger.Fatal(nil, "Template already exists: "+name)
	}
	theTemplate, err := template.ParseFS(fs, files...)
	if err != nil {
		log.Fatal("Error parsing template: ", err)
	}
	templatesByName[name] = theTemplate
	logger.Trace(nil, "Template loaded: "+name)
}
