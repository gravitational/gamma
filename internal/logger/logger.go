package logger

import (
	"fmt"
	"os"

	"github.com/gravitational/gamma/internal/color"
)

var (
	infoColor    = color.Magenta("ℹ")
	successColor = color.Green("✔")
	warningColor = color.Yellow("⚠")
	errorColor   = color.Red("✖")
)

func Info(message any) {
	fmt.Printf("%s %s\n", infoColor, message)
}

func Infof(format string, a ...any) {
	fmt.Printf("%s ", infoColor)
	fmt.Printf(format, a...)
	fmt.Print("\n")
}

func Success(message any) {
	fmt.Printf("%s %s\n", successColor, message)
}

func Successf(format string, a ...any) {
	fmt.Printf("%s ", successColor)
	fmt.Printf(format, a...)
	fmt.Print("\n")
}

func Warning(message any) {
	fmt.Printf("%s %s\n", warningColor, message)
}

func Warningf(format string, a ...any) {
	fmt.Printf("%s ", warningColor)
	fmt.Printf(format, a...)
	fmt.Print("\n")
}

func Error(message any) {
	fmt.Printf("%s %s\n", errorColor, message)
}

func Errorf(format string, a ...any) {
	fmt.Printf("%s ", errorColor)
	fmt.Printf(format, a...)
	fmt.Print("\n")
}

func Fatal(message any) {
	fmt.Printf("%s %s\n", errorColor, message)
	os.Exit(1)
}

func Fatalf(format string, a ...any) {
	Errorf(format, a...)
	os.Exit(1)
}
