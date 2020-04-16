package vcmd

import (
	"bufio"
	"context"
	"io/ioutil"
	"os/exec"
	"syscall"
	"time"
)

// Command 运行命令
func Command(name string, args ...string) (stdout []byte, stderr []byte, err error) {
	var (
		cmd *exec.Cmd
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	cmd = exec.CommandContext(ctx, name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	stdoutPipe, err := cmd.StdoutPipe()
	stderrPipe, err := cmd.StderrPipe()

	if err = cmd.Start(); err != nil {
		return
	}

	r := bufio.NewReader(stdoutPipe)
	stdout, err = ioutil.ReadAll(r)

	r = bufio.NewReader(stderrPipe)
	stderr, err = ioutil.ReadAll(r)

	if err = cmd.Wait(); err != nil {
		return
	}

	return
}
