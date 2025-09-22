// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package security

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSessionId() string {
	idBytes := make([]byte, 24)
	_, err := rand.Read(idBytes)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(idBytes)
}
