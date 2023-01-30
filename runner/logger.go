package runner

import (
	"fmt"
	logPkg "log"
	"os"
	"sync"
	"time"

	"github.com/mattn/go-colorable"
)

var once sync.Once

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

type appLogWriter struct {
	logFile *os.File
}

func newAppLog(path string) *appLogWriter {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		appLog(err.Error())
	}
	return &appLogWriter{
		logFile: file,
	}
}

func (a appLogWriter) Write(p []byte) (n int, err error) {
	a.logFile.Write(p)

	return len(p), nil
}

func (a *appLogWriter) Close() error {
	return a.logFile.Close()
}
