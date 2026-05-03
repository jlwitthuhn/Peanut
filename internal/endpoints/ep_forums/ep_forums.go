// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_forums

import (
	"net/http"
	"peanut/internal/data"
	"peanut/internal/service"
)

func RegisterForumHandlers(mux *http.ServeMux, forumHomeDvao data.ForumHomeDvao, forumsService service.ForumsService) {
	registerForumHomeHandlers(mux, forumHomeDvao)
	registerForumIndexHandlers(mux, forumsService)
}
