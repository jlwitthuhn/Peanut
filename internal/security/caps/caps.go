// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package caps

import (
	"peanut/internal/logger"
	"peanut/internal/security/groupnames"
)

const Admin_Forums_Structure_Edit = "Admin/Forums/Structure/Edit"
const Admin_FrontPage_Edit = "Admin/FrontPage/Edit"
const Admin_Gui_View = "Admin/Gui/View"
const Admin_ScheduledJob_Run = "Admin/ScheduledJob/Run"
const Forum_Moderate = "Forum/Moderate"
const Forum_Thread_Post = "Forum/Thread/Post"
const Forum_Thread_Reply = "Forum/Thread/Reply"

func GetCapabilitiesForGroup(group string) map[string]struct{} {
	result := make(map[string]struct{})
	switch group {
	case groupnames.TurboAdmin:
		fallthrough
	case groupnames.Admin:
		result[Admin_Forums_Structure_Edit] = struct{}{}
		result[Admin_FrontPage_Edit] = struct{}{}
		result[Admin_Gui_View] = struct{}{}
		result[Admin_ScheduledJob_Run] = struct{}{}
		result[Forum_Moderate] = struct{}{}
		fallthrough
	case groupnames.User:
		result[Forum_Thread_Post] = struct{}{}
		result[Forum_Thread_Reply] = struct{}{}
		fallthrough
	case groupnames.Guest:
	default:
		logger.Error(nil, "Attempted to get capabilities for illegal group: ", group)
	}
	return result
}

func GetCapabilitiesForGroups(groups ...string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, group := range groups {
		for capability := range GetCapabilitiesForGroup(group) {
			result[capability] = struct{}{}
		}
	}
	return result
}
