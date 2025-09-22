// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package passhash

func GenerateDefaultPhcString(plaintext string) string {
	return EncodeDefaultArgon2IdPhcString(plaintext)
}

func ValidatePassword(password string, phcString string) bool {
	if IsPlaintextPhcString(phcString) {
		return ValidatePlaintextPassword(password, phcString)
	} else if IsArgon2IdPhcString(phcString) {
		return ValidateArgon2IdPhcString(password, phcString)
	}
	return false
}
