// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_forums

import (
	"net/http"
	"peanut/internal/data"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
)

type forumHomeSection struct {
	Name   string
	Forums []forumHomeForum
}

type forumHomeForum struct {
	Name string
	Id   string
}

func groupForumsBySection(rows []data.ForumHomeViewRow) []forumHomeSection {
	var sections []forumHomeSection
	var currentSection *forumHomeSection

	for _, row := range rows {
		if currentSection == nil || currentSection.Name != row.SectionName {
			sections = append(sections, forumHomeSection{Name: row.SectionName})
			currentSection = &sections[len(sections)-1]
		}
		currentSection.Forums = append(currentSection.Forums, forumHomeForum{Name: row.ForumName, Id: row.ForumId})
	}

	return sections
}

func registerForumHomeHandlers(mux *http.ServeMux, forumHomeDvao data.ForumHomeDvao) {
	getForumHomeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rows, err := forumHomeDvao.SelectForumHomeViewRowsPublic(r.Context())
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forums.", w, r)
			return
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["Sections"] = groupForumsBySection(rows)
		ep_util.RenderTemplate("_forum/index", templateCtx, w, r)
	})
	mux.Handle("GET /forum", getForumHomeHandler)
}
