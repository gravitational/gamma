package logger

import (
	"fmt"
	"os"

	"github.com/gravitational/gamma/internal/color"
)

var (
	info    = color.Magenta("ℹ")
	success = color.Green("✔")
	warning = color.Yellow("⚠")
	error   = color.Red("✖")
)

func Info(message any) {
	fmt.Printf("%s %s\n", info, message)
}

func Infof(format string, a ...any) {
	fmt.Printf("%s ", info)
	fmt.Printf(format, a...)
	fmt.Print("\n")
}

func Success(message any) {
	fmt.Printf("%s %s\n", success, message)
}

func Successf(format string, a ...any) {
	fmt.Printf("%s ", success)
	fmt.Printf(format, a...)
	fmt.Print("\n")
}

func Warning(message any) {
	fmt.Printf("%s %s\n", warning, message)
}

func Warningf(format string, a ...any) {
	fmt.Printf("%s ", warning)
	fmt.Printf(format, a...)
	fmt.Print("\n")
}

func Error(message any) {
	fmt.Printf("%s %s\n", error, message)
}

func Errorf(format string, a ...any) {
	fmt.Printf("%s ", error)
	fmt.Printf(format, a...)
	fmt.Print("\n")
}

func Fatal(message any) {
	fmt.Printf("%s %s\n", error, message)
	os.Exit(1)
}

func Fatalf(format string, a ...any) {
	Errorf(format, a...)
	os.Exit(1)
}
