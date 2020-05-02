package vcmd

import (
	"bufio"
	"context"
	"io/ioutil"
	"os/exec"
	"syscall"
	"time"
)

// Command 运行windows命令，兼容所有系统
func Command(name string, args ...string) (stdout []byte, stderr []byte, err error) {
	var (
		cmd        *exec.Cmd
		argsHandle []string
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	argsHandle = append(args, "/C", name)
	argsHandle = append(argsHandle, args...)
	cmd = exec.CommandContext(ctx, "cmd", argsHandle...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

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
