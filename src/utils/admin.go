package utils

import (
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

func IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func RunAsAdmin() {
	exe, err := os.Executable()
	if err != nil {
		return
	}

	cwd, _ := os.Getwd()
	verb, _ := syscall.UTF16PtrFromString("runas")
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)

	windows.ShellExecute(0, verb, exePtr, nil, cwdPtr, windows.SW_NORMAL)
}
