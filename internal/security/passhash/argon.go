// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package passhash

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"unicode"

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

func IsArgon2IdPhcString(maybePhc string) bool {
	segments := strings.Split(maybePhc, "$")
	if len(segments) != 6 {
		return false
	}
	if segments[0] != "" {
		return false
	}
	if segments[1] != "argon2id" {
		return false
	}
	if segments[2] != "v=19" {
		return false
	}
	return true
}

func ValidateArgon2IdPhcString(plaintext string, phc string) bool {
	if IsArgon2IdPhcString(phc) == false {
		return false
	}
	segments := strings.Split(phc, "$")
	params := segments[3]
	saltB64 := segments[4]
	hashB64 := segments[5]
	m, t, p := parsePhcMtp(params)
	salt, saltErr := base64.StdEncoding.DecodeString(saltB64)
	if saltErr != nil {
		return false
	}

	newHash := argon2.IDKey([]byte(plaintext), salt, t, m, p, 32)

	hash, hashErr := base64.StdEncoding.DecodeString(hashB64)
	if hashErr != nil {
		return false
	}

	return bytes.Equal(newHash, hash)
}

func genSalt() []byte {
	saltBytes := make([]byte, 8)
	_, err := rand.Read(saltBytes)
	if err != nil {
		panic(err)
	}
	return saltBytes
}

func parsePhcMtp(input string) (uint32, uint32, uint8) {
	m := parsePhcParam(input, "m=")
	t := parsePhcParam(input, "t=")
	p := parsePhcParam(input, "p=")
	return m, t, uint8(p)
}

func parsePhcParam(input string, prefix string) uint32 {
	prefixIndex := strings.Index(input, prefix)
	numberBeginIndex := prefixIndex + len(prefix)
	if numberBeginIndex == len(input) {
		// Nothing exists after '='
		return 0
	}
	numberEndIndex := numberBeginIndex + 1
	for numberEndIndex < len(input) && unicode.IsDigit(rune(input[numberEndIndex])) {
		numberEndIndex++
	}
	substr := input[numberBeginIndex:numberEndIndex]
	result, err := strconv.Atoi(substr)
	if err != nil {
		return 0
	}
	return uint32(result)
}
