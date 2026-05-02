// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package passhash

import "testing"

const (
	testArg2Memory   uint32 = 8
	testArg2Time     uint32 = 1
	testArg2Parallel uint8  = 1
)

var testArg2Salt = []byte{1, 2, 3, 4, 5, 6, 7, 8}

func TestValidatePassword_Plaintext(t *testing.T) {
	phc := EncodePlaintextPhcString("hunter2")
	if !ValidatePassword("hunter2", phc) {
		t.Fatalf("ValidatePassword: correct plaintext password rejected")
	}
	if ValidatePassword("wrong", phc) {
		t.Fatalf("ValidatePassword: wrong plaintext password accepted")
	}
}

func TestValidatePassword_Argon2Id(t *testing.T) {
	phc := EncodeArgon2IdPhcString("hunter2", testArg2Salt, testArg2Memory, testArg2Time, testArg2Parallel)
	if !ValidatePassword("hunter2", phc) {
		t.Fatalf("ValidatePassword: correct argon2id password rejected")
	}
	if ValidatePassword("wrong", phc) {
		t.Fatalf("ValidatePassword: wrong argon2id password accepted")
	}
}

func TestValidatePassword_BadFormat(t *testing.T) {
	if ValidatePassword("hunter2", "not-a-phc-string") {
		t.Fatalf("ValidatePassword: bad PHC format accepted")
	}
	if ValidatePassword("hunter2", "") {
		t.Fatalf("ValidatePassword: empty PHC accepted")
	}
}
