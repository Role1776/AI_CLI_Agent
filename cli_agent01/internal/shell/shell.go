package shell

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

func ExecuteCommand(command string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-NoProfile", "-Command", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}

	output, err := cmd.CombinedOutput()

	return strings.ToValidUTF8(string(output), ""), err
}

func GetSystemInfo() string {
	osName := runtime.GOOS

	shell := "bash"
	if osName == "windows" {
		shell = "PowerShell"
	}

	currentUser, err := user.Current()
	username := "unknown"
	if err == nil {
		username = currentUser.Username
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "unknown"
	}

	return fmt.Sprintf(
		"\n- OS: %s\n- Shell: %s\n- User: %s\n- CWD: %s\n",
		osName, shell, username, cwd,
	)
}
