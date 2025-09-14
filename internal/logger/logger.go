// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package logger

import (
	"log"
	"os"
)

var theLogger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)

func Trace(args ...any) {
	fullArgs := append([]any{"[TRAC]"}, args...)
	theLogger.Println(fullArgs...)
}

func Info(args ...any) {
	fullArgs := append([]any{"[INFO]"}, args...)
	theLogger.Println(fullArgs...)
}

func Warn(args ...any) {
	fullArgs := append([]any{"[WARN]"}, args...)
	theLogger.Println(fullArgs...)
}

func Error(args ...any) {
	fullArgs := append([]any{"[EROR]"}, args...)
	theLogger.Println(fullArgs...)
}

func Fatal(args ...any) {
	fullArgs := append([]any{"[!!!!]"}, args...)
	theLogger.Fatal(fullArgs...)
}
