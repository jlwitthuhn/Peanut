// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package perms

import (
	"peanut/internal/logger"
	"peanut/internal/perms/permgroups"
)

const Admin_Gui_View = "Admin/Gui/View"

func GetPermissionsForGroup(group string) map[string]struct{} {
	result := make(map[string]struct{})
	switch group {
	case permgroups.TurboAdmin:
		fallthrough
	case permgroups.Admin:
		result[Admin_Gui_View] = struct{}{}
		fallthrough
	case permgroups.User:
	case permgroups.Guest:
	default:
		logger.Error(nil, "Attempted to get perms for illegal group: ", group)
	}
	return result
}

func GetPermissionsUnionForGroups(groups ...string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, group := range groups {
		groupPerms := GetPermissionsForGroup(group)
		for perm := range groupPerms {
			result[perm] = struct{}{}
		}
	}
	return result
}
