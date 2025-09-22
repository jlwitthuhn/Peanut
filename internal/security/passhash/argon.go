// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package passhash

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func EncodeArgon2IdPhcString(plaintext string, salt string, memory uint32, time uint32, parallel uint8) string {
	saltB64 := base64.StdEncoding.EncodeToString([]byte(salt))
	bytes := argon2.IDKey([]byte(plaintext), []byte(salt), time, memory, parallel, 32)
	b64 := base64.StdEncoding.EncodeToString(bytes)
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", memory, time, parallel, saltB64, b64)
}

func EncodeDefaultArgon2IdPhcString(plaintext string) string {
	salt := "salt"
	var memory uint32 = 24 * 1024
	var time uint32 = 2
	var parallel uint8 = 1
	return EncodeArgon2IdPhcString(plaintext, salt, memory, time, parallel)
}
