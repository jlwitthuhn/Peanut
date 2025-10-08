// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package security

func GenerateCsrfSmallToken() string {
	return GenerateSecureBase64Token(8)
}

func GenerateCsrfToken() string {
	return GenerateSecureBase64Token(16)
}
