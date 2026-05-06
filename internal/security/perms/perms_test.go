// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package perms

import (
	"strings"
	"testing"

	"peanut/internal/security/perms/permgroups"
)

var testAdminPerms = []string{
	Admin_FrontPage_Edit,
	Admin_Gui_View,
}

func TestGetPermissionsForGroup(t *testing.T) {
	t.Run("Admin has expected admin perms", func(t *testing.T) {
		got := GetPermissionsForGroup(permgroups.Admin)
		for _, p := range testAdminPerms {
			if _, ok := got[p]; !ok {
				t.Errorf("Admin: missing perm %q", p)
			}
		}
	})

	t.Run("TurboAdmin has expected admin perms", func(t *testing.T) {
		got := GetPermissionsForGroup(permgroups.TurboAdmin)
		for _, p := range testAdminPerms {
			if _, ok := got[p]; !ok {
				t.Errorf("TurboAdmin: missing perm %q", p)
			}
		}
	})

	t.Run("User has no admin perms", func(t *testing.T) {
		got := GetPermissionsForGroup(permgroups.User)
		for p := range got {
			if strings.HasPrefix(p, "Admin/") {
				t.Errorf("User: unexpected admin perm %q", p)
			}
		}
	})

	t.Run("Guest has no admin perms", func(t *testing.T) {
		got := GetPermissionsForGroup(permgroups.Guest)
		for p := range got {
			if strings.HasPrefix(p, "Admin/") {
				t.Errorf("Guest: unexpected admin perm %q", p)
			}
		}
	})
}

func TestGetGranularPermissionsForGroups_Union(t *testing.T) {
	got := GetGranularPermissionsForGroups(permgroups.User, permgroups.Admin)
	expected := testAdminPerms
	for _, p := range expected {
		if _, ok := got[p]; !ok {
			t.Errorf("missing perm %q", p)
		}
	}
}

func TestGetGranularPermissionsForGroups_Empty(t *testing.T) {
	got := GetGranularPermissionsForGroups()
	if len(got) != 0 {
		t.Fatalf("no groups: got %d perms, want 0", len(got))
	}

	got = GetGranularPermissionsForGroups(permgroups.Guest)
	if len(got) != 0 {
		t.Fatalf("Guest+User: got %d perms, want 0", len(got))
	}
}
