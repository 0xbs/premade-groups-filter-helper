package logger

import (
	"fmt"
	"os"
	"time"
)

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var purple = "\033[35m"
var cyan = "\033[36m"
var gray = "\033[37m"
var white = "\033[97m"

func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05.000000")
}

func Infof(format string, v ...any) {
	msg := timestamp() + green + " INFO" + reset + " " + format + "\n"
	fmt.Printf(msg, v...)
}

func Warnf(format string, v ...any) {
	msg := timestamp() + yellow + " WARN" + reset + " " + format + "\n"
	fmt.Printf(msg, v...)
}

func Errorf(format string, v ...any) {
	msg := timestamp() + red + " ERROR" + reset + " " + format + "\n"
	fmt.Printf(msg, v...)
}

func Fatalf(format string, v ...any) {
	msg := timestamp() + purple + " FATAL" + reset + " " + format + "\n"
	fmt.Printf(msg, v...)
	os.Exit(1)
}
