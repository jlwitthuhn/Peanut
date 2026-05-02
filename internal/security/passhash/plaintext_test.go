// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package passhash

import "testing"

func TestIsPlaintextPhcString(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid", "$plaintext$aGVsbG8=", true},
		{"wrong segment count", "$plaintext$aGVsbG8=$extra", false},
		{"wrong algorithm", "$argon2id$aGVsbG8=", false},
		{"missing leading dollar", "plaintext$aGVsbG8=", false},
		{"empty payload", "$plaintext$", false},
		{"empty string", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsPlaintextPhcString(tc.input); got != tc.want {
				t.Fatalf("IsPlaintextPhcString(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestValidatePlaintextPassword(t *testing.T) {
	phc := EncodePlaintextPhcString("hunter2")

	if !ValidatePlaintextPassword("hunter2", phc) {
		t.Fatalf("ValidatePlaintextPassword: round-trip failed for correct password")
	}
	if ValidatePlaintextPassword("wrong", phc) {
		t.Fatalf("ValidatePlaintextPassword: accepted wrong password")
	}
	if ValidatePlaintextPassword("hunter2", "not-a-phc-string") {
		t.Fatalf("ValidatePlaintextPassword: accepted malformed PHC")
	}
}
