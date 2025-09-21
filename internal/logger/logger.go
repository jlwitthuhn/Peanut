// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"peanut/internal/middleutil"
)

var theLogger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)

func Trace(r *http.Request, args ...any) {
	fullArgs := append([]any{formatRequestId(r), "[TRAC]"}, args...)
	theLogger.Println(fullArgs...)
}

func Info(r *http.Request, args ...any) {
	fullArgs := append([]any{formatRequestId(r), "[INFO]"}, args...)
	theLogger.Println(fullArgs...)
}

func Warn(r *http.Request, args ...any) {
	fullArgs := append([]any{formatRequestId(r), "[WARN]"}, args...)
	theLogger.Println(fullArgs...)
}

func Error(r *http.Request, args ...any) {
	fullArgs := append([]any{formatRequestId(r), "[EROR]"}, args...)
	theLogger.Println(fullArgs...)
}

func Fatal(r *http.Request, args ...any) {
	fullArgs := append([]any{formatRequestId(r), "[!!!!]"}, args...)
	theLogger.Fatal(fullArgs...)
}

func formatRequestId(r *http.Request) string {
	return fmt.Sprintf("(%s)", middleutil.RetrieveRequestId(r))
}
