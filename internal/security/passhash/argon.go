// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package passhash

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func EncodeArgon2IdPhcString(plaintext string, salt []byte, memory uint32, time uint32, parallel uint8) string {
	saltB64 := base64.StdEncoding.EncodeToString(salt)
	bytes := argon2.IDKey([]byte(plaintext), salt, time, memory, parallel, 32)
	b64 := base64.StdEncoding.EncodeToString(bytes)
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", memory, time, parallel, saltB64, b64)
}

func EncodeDefaultArgon2IdPhcString(plaintext string) string {
	salt := genSalt()
	var memory uint32 = 24 * 1024
	var time uint32 = 2
	var parallel uint8 = 1
	return EncodeArgon2IdPhcString(plaintext, salt, memory, time, parallel)
}

func genSalt() []byte {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}
