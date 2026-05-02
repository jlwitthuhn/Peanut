// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package dataformat

import (
	"testing"
	"time"
)

func TestFormatDurationAsPostgresInterval(t *testing.T) {
	cases := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{"zero", 0, "0.000000 minutes"},
		{"seconds", 90 * time.Second, "1.500000 minutes"},
		{"minutes", time.Minute, "1.000000 minutes"},
		{"hours", time.Hour, "60.000000 minutes"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatDurationAsPostgresInterval(tc.duration)
			if got != tc.want {
				t.Fatalf("FormatDurationAsPostgresInterval(%v) = %q, want %q", tc.duration, got, tc.want)
			}
		})
	}
}
