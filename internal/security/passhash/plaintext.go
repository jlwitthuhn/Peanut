// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package passhash

import (
	"encoding/base64"
	"fmt"
)

func GeneratePlaintextPhcString(plaintext string) string {
	base64String := base64.StdEncoding.EncodeToString([]byte(plaintext))
	return fmt.Sprintf("$plaintext$%s", base64String)
}
