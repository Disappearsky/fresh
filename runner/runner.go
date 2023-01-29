package runner

import (
	"io"
	"os"
	"os/exec"
)

func run() bool {
	runnerLog("Running...")

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

	go io.Copy(os.Stderr, stderr)
	go io.Copy(os.Stdout, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
	}()

	return true
}
