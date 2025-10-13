// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package dataformat

import (
	"fmt"
	"time"
)

func FormatDurationAsPostgresInterval(duration time.Duration) string {
	output := fmt.Sprintf("%f minutes", duration.Minutes())
	return output
}
