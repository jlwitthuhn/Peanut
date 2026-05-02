// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package validator

import "testing"

func TestValidateEmail(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty", "", true},
		{"one char", "a", true},
		{"normal", "user@example.com", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEmail(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ValidateEmail(%q) error = %v, wantErr = %v", tc.input, err, tc.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty", "", true},
		{"one char", "x", true},
		{"two chars", "xy", false},
		{"long", "This password is a sentence!", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ValidatePassword(%q) error = %v, wantErr = %v", tc.input, err, tc.wantErr)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty", "", true},
		{"one char", "a", true},
		{"two chars", "ab", false},
		{"normal", "pingas", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateUsername(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ValidateUsername(%q) error = %v, wantErr = %v", tc.input, err, tc.wantErr)
			}
		})
	}
}
