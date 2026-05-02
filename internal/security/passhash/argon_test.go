// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package passhash

import "testing"

func TestIsArgon2IdPhcString(t *testing.T) {
	validPhc := EncodeArgon2IdPhcString("hunter2", testArg2Salt, testArg2Memory, testArg2Time, testArg2Parallel)
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid", validPhc, true},
		{"wrong segment count", "$argon2id$v=19$m=8,t=1,p=1$salt", false},
		{"wrong algorithm", "$argon2i$v=19$m=8,t=1,p=1$c2FsdA==$aGFzaA==", false},
		{"wrong version", "$argon2id$v=18$m=8,t=1,p=1$c2FsdA==$aGFzaA==", false},
		{"missing leading dollar", "argon2id$v=19$m=8,t=1,p=1$c2FsdA==$aGFzaA==", false},
		{"empty string", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsArgon2IdPhcString(tc.input); got != tc.want {
				t.Fatalf("IsArgon2IdPhcString(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestValidateArgon2IdPassword(t *testing.T) {
	phc := EncodeArgon2IdPhcString("hunter2", testArg2Salt, testArg2Memory, testArg2Time, testArg2Parallel)

	if !ValidateArgon2IdPhcString("hunter2", phc) {
		t.Fatalf("ValidateArgon2IdPhcString: round-trip failed for correct password")
	}
	if ValidateArgon2IdPhcString("wrong", phc) {
		t.Fatalf("ValidateArgon2IdPhcString: accepted wrong password")
	}
	if ValidateArgon2IdPhcString("hunter2", "not-a-phc-string") {
		t.Fatalf("ValidateArgon2IdPhcString: accepted malformed PHC")
	}
}

func TestParsePhcParam(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		prefix string
		want   uint32
	}{
		{"test1", "m=5,t=2,p=1", "m=", 5},
		{"test2", "m=512,t=2,p=1", "m=", 512},
		{"test3", "m=8,t=2,p=1", "p=", 1},
		{"test4", "m=24576", "m=", 24576},
		{"test5", "m=", "m=", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parsePhcParam(tc.input, tc.prefix)
			if got != tc.want {
				t.Fatalf("parsePhcParam(%q, %q) = %d, want %d", tc.input, tc.prefix, got, tc.want)
			}
		})
	}
}
