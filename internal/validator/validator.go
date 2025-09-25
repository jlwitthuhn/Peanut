// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package validator

import "errors"

func ValidateEmail(email string) error {
	if len(email) < 2 {
		return errors.New("Email must contain at least two characters")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 2 {
		return errors.New("Password must contain at least two characters")
	}
	return nil
}

func ValidateUsername(username string) error {
	if len(username) < 2 {
		return errors.New("Username must contain at least two characters")
	}
	return nil
}
