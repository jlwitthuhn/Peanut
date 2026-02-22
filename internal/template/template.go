// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package template

import (
	"html/template"
	"io/fs"
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

	adminGroupsFiles := []string{"base.html", "css/common.css", "admin/base.html", "admin/groups.html"}
	loadTemplateOrDie(fs, "_admin/groups", adminGroupsFiles...)

	adminGroupsListFiles := []string{"base.html", "css/common.css", "admin/base.html", "admin/groups_list.html"}
	loadTemplateOrDie(fs, "_admin/groups_list", adminGroupsListFiles...)

	adminScheduledJobFiles := []string{"base.html", "css/common.css", "admin/base.html", "admin/scheduled_jobs.html"}
	loadTemplateOrDie(fs, "_admin/scheduled_jobs", adminScheduledJobFiles...)

	adminSimpleMessageFiles := []string{"base.html", "css/common.css", "admin/base.html", "admin/simple_message.html"}
	loadTemplateOrDie(fs, "_admin/simple_message", adminSimpleMessageFiles...)

	adminUsersFiles := []string{"base.html", "css/common.css", "admin/base.html", "admin/users.html"}
	loadTemplateOrDie(fs, "_admin/users", adminUsersFiles...)

	adminUsersListFiles := []string{"base.html", "css/common.css", "admin/base.html", "admin/users_list.html"}
	loadTemplateOrDie(fs, "_admin/users_list", adminUsersListFiles...)

	indexFiles := []string{"base.html", "css/common.css", "index.html"}
	loadTemplateOrDie(fs, "_index", indexFiles...)

	loginFiles := []string{"base.html", "css/common.css", "login.html"}
	loadTemplateOrDie(fs, "_login", loginFiles...)

	profileFiles := []string{"base.html", "css/common.css", "profile.html"}
	loadTemplateOrDie(fs, "_profile", profileFiles...)

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
		logger.Fatal(nil, "Error parsing template: ", err)
	}
	templatesByName[name] = theTemplate
	logger.Trace(nil, "Template loaded: "+name)
}
