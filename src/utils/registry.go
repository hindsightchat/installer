package utils

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func DetectInstallDir(appName, exeName string) string {
	// check registry first
	if dir, ok := readRegistry(appName); ok {
		return dir
	}

	// check exe location
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		if _, err := os.Stat(filepath.Join(exeDir, exeName)); err == nil {
			return exeDir
		}
		parent := filepath.Dir(exeDir)
		if _, err := os.Stat(filepath.Join(parent, exeName)); err == nil {
			return parent
		}
	}

	return defaultInstallDir(appName)
}

func defaultInstallDir(appName string) string {
	pf := os.Getenv("ProgramFiles")
	if pf == "" {
		pf = `C:\Program Files`
	}
	return filepath.Join(pf, appName)
}

func readRegistry(appName string) (string, bool) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\`+appName, registry.QUERY_VALUE)
	if err != nil {
		return "", false
	}
	defer k.Close()

	path, _, err := k.GetStringValue("path")
	if err != nil {
		return "", false
	}
	return path, true
}