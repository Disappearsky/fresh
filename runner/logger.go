package runner

import (
	"fmt"
	logPkg "log"
	"os"
	"time"

	"github.com/mattn/go-colorable"
)

var file *os.File

func init() {
	var err error
	file, err = os.OpenFile(buildErrorsFilePath(), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		appLog(err.Error())
	}
}

type logFunc func(string, ...interface{})

var logger = logPkg.New(colorable.NewColorableStderr(), "", 0)

func newLogFunc(prefix string) func(string, ...interface{}) {
	color, clear := "", ""
	if settings["colors"] == "1" {
		color = fmt.Sprintf("\033[%sm", logColor(prefix))
		clear = fmt.Sprintf("\033[%sm", colors["reset"])
	}
	prefix = fmt.Sprintf("%-11s", prefix)

	return func(format string, v ...interface{}) {
		now := time.Now()
		timeString := fmt.Sprintf("%d:%d:%02d", now.Hour(), now.Minute(), now.Second())
		format = fmt.Sprintf("%s%s %s |%s %s", color, timeString, prefix, clear, format)
		logger.Printf(format, v...)
	}
}

func fatal(err error) {
	logger.Fatal(err)
}

type appLogWriter struct{}

func (a appLogWriter) Write(p []byte) (n int, err error) {
	// write log to local text
	file.Write(p)

	return len(p), nil
}
