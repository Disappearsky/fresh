package runner

import (
	"io"
	"os/exec"
)

func run() bool {
	runnerLog("Running...")
	appLogFile := newAppLog(buildErrorsFilePath())

	cmd := exec.Command(buildPath(), args())

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	go io.Copy(appLogFile, stderr)
	go io.Copy(appLogFile, stdout)

	go func() {
		<-stopChannel
		appLogFile.Close()
		runnerLog("Close log %s", buildErrorsFilePath())
		err = removeBuildErrorsLog()
		if err != nil {
			runnerLog(err.Error())
		}
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
		fileSig <- struct{}{}
	}()

	return true
}
