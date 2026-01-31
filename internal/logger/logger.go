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

func getLogLevel() int {
	envLogLevel := os.Getenv("PEANUT_LOG_LEVEL")
	switch envLogLevel {
	case "fatal":
		return 0
	case "error":
		return 1
	case "warn":
		return 2
	case "info":
		return 3
	case "debug":
		return 4
	case "trace":
		return 5
	default:
		return 5
	}
}

var theLogger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)
var logLevel = getLogLevel()

func Trace(r *http.Request, args ...any) {
	if logLevel >= 5 {
		fullArgs := append([]any{formatRequestId(r), "[TRAC]", getCallerLocation()}, args...)
		theLogger.Println(fullArgs...)
	}
}

func Debug(r *http.Request, args ...any) {
	if logLevel >= 4 {
		fullArgs := append([]any{formatRequestId(r), "[DBUG]", getCallerLocation()}, args...)
		theLogger.Println(fullArgs...)
	}
}

func Info(r *http.Request, args ...any) {
	if logLevel >= 3 {
		fullArgs := append([]any{formatRequestId(r), "[INFO]", getCallerLocation()}, args...)
		theLogger.Println(fullArgs...)
	}
}

func Warn(r *http.Request, args ...any) {
	if logLevel >= 2 {
		fullArgs := append([]any{formatRequestId(r), "[WARN]", getCallerLocation()}, args...)
		theLogger.Println(fullArgs...)
	}
}

func Error(r *http.Request, args ...any) {
	if logLevel >= 1 {
		fullArgs := append([]any{formatRequestId(r), "[EROR]", getCallerLocation()}, args...)
		theLogger.Println(fullArgs...)
	}
}

func Fatal(r *http.Request, args ...any) {
	if logLevel >= 0 {
		fullArgs := append([]any{formatRequestId(r), "[!!!!]", getCallerLocation()}, args...)
		theLogger.Fatal(fullArgs...)
	} else {
		// End process if logging is disabled
		os.Exit(1)
	}
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
