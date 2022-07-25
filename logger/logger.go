package logger

import (
	"fmt"
	"io"
	"os"
)

//logger包提供一些列的日志处理函数
var writer io.Writer = fileWriter("log/log")

func fileWriter(path string) *os.File {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		return nil
	}
	return file
}

func Fatalf(format string, args ...any) {
	fmt.Fprintf(writer, "Fatal error"+format+"\n", args)
}
func Fatal(format string) {
	fmt.Fprintln(writer, "Fatal error"+format)
}

func Infof(format string, args ...any) {
	fmt.Fprintf(writer, "Info "+format+"\n", args)
}

func Info(format string) {
	fmt.Fprintln(writer, "Info "+format)
}

func Warnf(format string, args ...any) {
	fmt.Fprintf(writer, "Warn "+format+"\n", args)
}

func Warn(format string) {
	fmt.Fprintln(writer, "Warn"+format)
}

func Errorf(format string, args ...any) {
	fmt.Fprintf(writer, "Error "+format+"\n", args)
}

func Error(format string) {
	fmt.Fprintln(writer, "Error "+format)
}
