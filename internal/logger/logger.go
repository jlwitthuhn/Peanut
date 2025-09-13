package logger

import (
	"log"
	"os"
)

var theLogger = log.New(os.Stdout, "[TRAC] ", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)

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
