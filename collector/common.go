package collector

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/astaxie/beego/logs"
)

var DebugMode bool

func RunCmd(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	outStr := stdoutBuf.String()
	errStr := stderrBuf.String()

	if DebugMode {
		logs.Debug("Command: %s %s", command, strings.Join(args, " "))
		logs.Debug("Stdout: %s", strings.TrimSpace(outStr))
		logs.Debug("Stderr: %s", strings.TrimSpace(errStr))
		logs.Debug("Status: %v", err)
	}

	if err != nil {
		return "", fmt.Errorf("command '%s %s' failed: %v\nstderr: %s", command, strings.Join(args, " "), err, errStr)
	}

	return outStr, nil
}

func BytesToMb(bytes uint64) uint64 {
	const megabyte = uint64(1024 * 1024)
	return bytes / megabyte
}
