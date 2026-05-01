// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleutil

import (
	"context"
	"peanut/internal/keynames/contextkeys"
)

func RetrieveRequestId(ctx context.Context) string {
	if ctx == nil {
		return "--------"
	}
	id, ok := ctx.Value(contextkeys.RequestId).(string)
	if !ok {
		return "--------"
	}
	return id
}
