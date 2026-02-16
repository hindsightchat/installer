package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func KillInstances(exeName, installDir string) {
	// kill by process name
	cmd := exec.Command("taskkill", "/F", "/IM", exeName)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()

	// kill any from install dir except this process
	ps := fmt.Sprintf(
		"Get-Process | Where-Object { $_.Path -eq '%s' -and $_.Id -ne %d } | Stop-Process -Force",
		filepath.Join(installDir, exeName),
		os.Getpid(),
	)
	cmd = exec.Command("powershell", "-Command", ps)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
}