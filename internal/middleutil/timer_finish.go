// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleutil

import (
	"net/http"
	"time"
)

var RequestTimerBeginKey string = "requestTimerBegin"

func RequestTimerFinish(req *http.Request) float64 {
	requestBegin := req.Context().Value(RequestTimerBeginKey).(time.Time)
	requestDurationUs := time.Now().Sub(requestBegin).Microseconds()
	requestDurationUs -= requestDurationUs % 10 // Two decimal places in milliseconds
	return float64(requestDurationUs) / 1000.0
}
