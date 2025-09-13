package logger

import (
	"log"
	"os"
)

var traceLogger = log.New(os.Stdout, "[TRAC] ", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)
var infoLogger = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)
var warnLogger = log.New(os.Stdout, "[WARN] ", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)
var errorLogger = log.New(os.Stdout, "[EROR] ", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)
var fatalLogger = log.New(os.Stdout, "[!!!!] ", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)

func Trace(args ...any) {
	traceLogger.Println(args...)
}

func Info(args ...any) {
	infoLogger.Println(args...)
}

func Warn(args ...any) {
	warnLogger.Println(args...)
}

func Error(args ...any) {
	errorLogger.Println(args...)
}

func Fatal(args ...any) {
	fatalLogger.Fatal(args...)
}
