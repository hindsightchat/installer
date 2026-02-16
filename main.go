package main

import (
	_ "embed"

	"github.com/hindsightchat/installer/src/updater"
	"github.com/hindsightchat/installer/src/utils"
)

//go:embed app.zip
var AppZipData []byte

//go:embed src/hindsight.png
var LogoPng []byte

const (
	AppName     = "Hindsight Chat"
	ExeName     = "hindsightchat.exe"
	RegistryKey = `Software\HindsightChat`
)

func main() {
	if !utils.IsAdmin() {
		utils.RunAsAdmin()
		return
	}

	u := updater.New(AppName, ExeName, AppZipData, LogoPng)
	u.Run()
}
