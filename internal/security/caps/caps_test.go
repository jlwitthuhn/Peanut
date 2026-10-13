// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package caps

import (
	"strings"
	"testing"

	"peanut/internal/security/groupnames"
)

var testAdminCaps = []string{
	Admin_FrontPage_Edit,
	Admin_Gui_View,
}

func TestGetCapabilitiesForGroup(t *testing.T) {
	t.Run("Admin has expected admin caps", func(t *testing.T) {
		got := GetCapabilitiesForGroup(groupnames.Admin)
		for _, c := range testAdminCaps {
			if _, ok := got[c]; !ok {
				t.Errorf("Admin: missing capability %q", c)
			}
		}
	})

	t.Run("TurboAdmin has expected admin caps", func(t *testing.T) {
		got := GetCapabilitiesForGroup(groupnames.TurboAdmin)
		for _, c := range testAdminCaps {
			if _, ok := got[c]; !ok {
				t.Errorf("TurboAdmin: missing capability %q", c)
			}
		}
	})

	t.Run("User has no admin caps", func(t *testing.T) {
		got := GetCapabilitiesForGroup(groupnames.User)
		for c := range got {
			if strings.HasPrefix(c, "Admin/") {
				t.Errorf("User: unexpected admin capability %q", c)
			}
		}
	})

	t.Run("Guest has no admin caps", func(t *testing.T) {
		got := GetCapabilitiesForGroup(groupnames.Guest)
		for c := range got {
			if strings.HasPrefix(c, "Admin/") {
				t.Errorf("Guest: unexpected admin capability %q", c)
			}
		}
	})
}

func TestGetCapabilitiesForGroups_Union(t *testing.T) {
	got := GetCapabilitiesForGroups(groupnames.User, groupnames.Admin)
	expected := testAdminCaps
	for _, c := range expected {
		if _, ok := got[c]; !ok {
			t.Errorf("missing capability %q", c)
		}
	}
}

func TestGetCapabilitiesForGroups_Empty(t *testing.T) {
	got := GetCapabilitiesForGroups()
	if len(got) != 0 {
		t.Fatalf("no groups: got %d capabilities, want 0", len(got))
	}

	got = GetCapabilitiesForGroups(groupnames.Guest)
	if len(got) != 0 {
		t.Fatalf("Guest: got %d capabilities, want 0", len(got))
	}
}
