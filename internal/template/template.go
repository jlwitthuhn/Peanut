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

	adminIndexFiles := []string{"base.html", "admin/base.html", "admin/index.html"}
	loadTemplateOrDie(fs, "view_admin/index", adminIndexFiles...)

	adminDebugFiles := []string{"base.html", "admin/base.html", "admin/debug.html"}
	loadTemplateOrDie(fs, "view_admin/debug", adminDebugFiles...)

	adminForumForumsFiles := []string{"base.html", "admin/base.html", "admin/forum/forums.html"}
	loadTemplateOrDie(fs, "view_admin/forum/forums", adminForumForumsFiles...)

	adminForumSectionsFiles := []string{"base.html", "admin/base.html", "admin/forum/sections.html"}
	loadTemplateOrDie(fs, "view_admin/forum/sections", adminForumSectionsFiles...)

	adminForumSectionsAddEditFiles := []string{"base.html", "admin/base.html", "admin/forum/sections_add_edit.html"}
	loadTemplateOrDie(fs, "view_admin/forum/sections/add_edit", adminForumSectionsAddEditFiles...)

	adminForumForumsAddEditFiles := []string{"base.html", "admin/base.html", "admin/forum/forums_add_edit.html"}
	loadTemplateOrDie(fs, "view_admin/forum/forums/add_edit", adminForumForumsAddEditFiles...)

	adminFrontPageFiles := []string{"base.html", "admin/base.html", "admin/front_page.html"}
	loadTemplateOrDie(fs, "view_admin/front_page", adminFrontPageFiles...)

	adminGroupsFiles := []string{"base.html", "admin/base.html", "admin/groups.html"}
	loadTemplateOrDie(fs, "view_admin/groups", adminGroupsFiles...)

	adminGroupsListFiles := []string{"base.html", "admin/base.html", "admin/groups_list.html"}
	loadTemplateOrDie(fs, "view_admin/groups_list", adminGroupsListFiles...)

	adminScheduledJobFiles := []string{"base.html", "admin/base.html", "admin/scheduled_jobs.html"}
	loadTemplateOrDie(fs, "view_admin/scheduled_jobs", adminScheduledJobFiles...)

	adminSimpleMessageFiles := []string{"base.html", "admin/base.html", "admin/simple_message.html"}
	loadTemplateOrDie(fs, "view_admin/simple_message", adminSimpleMessageFiles...)

	adminUsersFiles := []string{"base.html", "admin/base.html", "admin/users.html"}
	loadTemplateOrDie(fs, "view_admin/users", adminUsersFiles...)

	adminUsersListFiles := []string{"base.html", "admin/base.html", "admin/users_list.html"}
	loadTemplateOrDie(fs, "view_admin/users_list", adminUsersListFiles...)

	indexFiles := []string{"base.html", "index.html"}
	loadTemplateOrDie(fs, "view_index", indexFiles...)

	loginFiles := []string{"base.html", "login.html"}
	loadTemplateOrDie(fs, "view_login", loginFiles...)

	profileFiles := []string{"base.html", "profile.html"}
	loadTemplateOrDie(fs, "view_profile", profileFiles...)

	registerFiles := []string{"base.html", "register.html"}
	loadTemplateOrDie(fs, "view_register", registerFiles...)

	setupFiles := []string{"base.html", "setup.html"}
	loadTemplateOrDie(fs, "view_setup", setupFiles...)

	forumListingFiles := []string{"base.html", "forum/forum_listing.html", "widget/forum_section.html"}
	loadTemplateOrDie(fs, "view_forum/forum_listing", forumListingFiles...)

	forumModerateFiles := []string{"base.html", "forum/moderate.html"}
	loadTemplateOrDie(fs, "view_forum/moderate", forumModerateFiles...)

	forumNewThreadFiles := []string{"base.html", "forum/new_thread.html"}
	loadTemplateOrDie(fs, "view_forum/new_thread", forumNewThreadFiles...)

	forumPostListingFiles := []string{"base.html", "forum/post_listing.html"}
	loadTemplateOrDie(fs, "view_forum/post_listing", forumPostListingFiles...)

	forumReplyFiles := []string{"base.html", "forum/reply.html"}
	loadTemplateOrDie(fs, "view_forum/reply", forumReplyFiles...)

	forumThreadFiles := []string{"base.html", "forum/thread_listing.html"}
	loadTemplateOrDie(fs, "view_forum/thread_listing", forumThreadFiles...)

	simpleMessageFiles := []string{"base.html", "simple_message.html"}
	loadTemplateOrDie(fs, "view_simple_message", simpleMessageFiles...)
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
