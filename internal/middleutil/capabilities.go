// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleutil

import (
	"context"
	"peanut/internal/keynames/contextkeys"
)

func ContextHasCapability(ctx context.Context, capability string) bool {
	capabilities, ok := ctx.Value(contextkeys.UserCaps).(map[string]struct{})
	if !ok {
		return false
	}
	_, has := capabilities[capability]
	return has
}
