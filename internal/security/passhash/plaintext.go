// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package passhash

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func EncodePlaintextPhcString(plaintext string) string {
	base64String := base64.StdEncoding.EncodeToString([]byte(plaintext))
	return fmt.Sprintf("$plaintext$%s", base64String)
}

func IsPlaintextPhcString(maybePhc string) bool {
	segments := strings.Split(maybePhc, "$")
	if len(segments) != 3 {
		return false
	}
	if segments[0] != "" {
		return false
	}
	if segments[1] != "plaintext" {
		return false
	}
	if len(segments[2]) == 0 {
		return false
	}
	return true
}

func ValidatePlaintextPassword(plaintext string, phc string) bool {
	if IsPlaintextPhcString(phc) == false {
		return false
	}
	segments := strings.Split(phc, "$")
	existingBase64 := segments[2]
	newBase64 := base64.StdEncoding.EncodeToString([]byte(plaintext))
	return existingBase64 == newBase64
}
