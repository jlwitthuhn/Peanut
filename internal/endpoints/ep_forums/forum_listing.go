// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_forums

import (
	"net/http"
	"peanut/internal/data"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/service"
)

type forumHomeSection struct {
	Id     string
	Name   string
	Forums []forumHomeForum
}

type forumHomeForum struct {
	Name         string
	Id           string
	Description  string
	ThreadCount  int32
	PostCount    int32
	LastPostDate string
}

func groupForumsBySection(rows []data.ForumHomeViewRow) []forumHomeSection {
	var sections []forumHomeSection
	var currentSection *forumHomeSection

	for _, row := range rows {
		if currentSection == nil || currentSection.Id != row.SectionId {
			sections = append(sections, forumHomeSection{Id: row.SectionId, Name: row.SectionName})
			currentSection = &sections[len(sections)-1]
		}
		lastPostDate := "Never"
		if row.LastTimestamp != nil {
			lastPostDate = row.LastTimestamp.Format("2006-01-02 15:04 MST")
		}
		currentSection.Forums = append(currentSection.Forums, forumHomeForum{Name: row.ForumName, Id: row.ForumId, Description: row.ForumDescription, ThreadCount: row.ThreadCount, PostCount: row.PostCount, LastPostDate: lastPostDate})
	}

	return sections
}

func registerForumListingHandlers(mux *http.ServeMux, forumsService service.ForumService) {
	getForumHomeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rows, err := forumsService.GetHomeViewRowsPublic(r.Context())
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forums.", w, r)
			return
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["Sections"] = groupForumsBySection(rows)
		ep_util.RenderTemplate("view_forum/forum_listing", templateCtx, w, r)
	})
	mux.Handle("GET /forum", getForumHomeHandler)

	getSectionIndexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sectionId := r.PathValue("sectionId")

		readable, err := forumsService.CanReadSection(r.Context(), sectionId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load section.", w, r)
			return
		}
		if !readable {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return
		}

		rows, err := forumsService.GetHomeViewRowsBySectionPublic(r.Context(), sectionId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forums.", w, r)
			return
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["Sections"] = groupForumsBySection(rows)
		ep_util.RenderTemplate("view_forum/forum_listing", templateCtx, w, r)
	})
	mux.Handle("GET /forum/section/{sectionId}", getSectionIndexHandler)
}
