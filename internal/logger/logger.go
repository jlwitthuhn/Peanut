// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"peanut/internal/middleutil"
	"runtime"
)

var theLogger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)

func Trace(r *http.Request, args ...any) {
	fullArgs := append([]any{formatRequestId(r), "[TRAC]", getCallerLocation()}, args...)
	theLogger.Println(fullArgs...)
}

func Debug(r *http.Request, args ...any) {
	fullArgs := append([]any{formatRequestId(r), "[DBUG]", getCallerLocation()}, args...)
	theLogger.Println(fullArgs...)
}

func Info(r *http.Request, args ...any) {
	fullArgs := append([]any{formatRequestId(r), "[INFO]", getCallerLocation()}, args...)
	theLogger.Println(fullArgs...)
}

func Warn(r *http.Request, args ...any) {
	fullArgs := append([]any{formatRequestId(r), "[WARN]", getCallerLocation()}, args...)
	theLogger.Println(fullArgs...)
}

func Error(r *http.Request, args ...any) {
	fullArgs := append([]any{formatRequestId(r), "[EROR]", getCallerLocation()}, args...)
	theLogger.Println(fullArgs...)
}

func Fatal(r *http.Request, args ...any) {
	fullArgs := append([]any{formatRequestId(r), "[!!!!]", getCallerLocation()}, args...)
	theLogger.Fatal(fullArgs...)
}

func formatRequestId(r *http.Request) string {
	return fmt.Sprintf("(%s)", middleutil.RetrieveRequestId(r))
}

func getCallerLocation() string {
	// 0 is getCallerLocation
	// 1 is Info, Warn, Error, etc
	// 2 is the code that called into this package
	_, file, line, _ := runtime.Caller(2)
	fileName := filepath.Base(file)
	return fmt.Sprintf("%s:%d:", fileName, line)
}
