// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleutil

import (
	"context"
	"peanut/internal/keynames/contextkeys"
)

func ContextHasPermission(ctx context.Context, permission string) bool {
	permissions := ctx.Value(contextkeys.UserPerms).(map[string]struct{})
	_, has := permissions[permission]
	return has
}
