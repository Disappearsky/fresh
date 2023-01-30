package runner

import (
	"io"
	"os/exec"
)

func run() bool {
	runnerLog("Running...")
	appLog := newAppLog(buildErrorsFilePath())

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

	go io.Copy(appLog, stderr)
	go io.Copy(appLog, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
		appLog.Close()
	}()

	return true
}
