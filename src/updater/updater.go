package updater

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/hindsightchat/installer/src/colours"
	apptheme "github.com/hindsightchat/installer/src/theme"
	"github.com/hindsightchat/installer/src/utils"
	"github.com/hindsightchat/installer/src/widgets"
)

type Window struct {
	app     fyne.App
	window  fyne.Window
	status  *widget.Label
	spinner *widgets.Spinner
	dir     string
	zipData []byte
	logo    []byte

	exeName string
	appName string

	fileLog *os.File
}

func New(appName string, exeName string, zipData, logo []byte) *Window {
	logFile, err := os.OpenFile("updater.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		os.Stdout = logFile
		os.Stderr = logFile
	}
	return &Window{
		dir:     utils.DetectInstallDir(appName, exeName),
		zipData: zipData,
		logo:    logo,
		appName: appName,
		exeName: exeName,
		fileLog: logFile,
	}
}

func (w *Window) Run() {

	w.app = app.New()
	w.app.Settings().SetTheme(&apptheme.Dark{})

	drv := fyne.CurrentApp().Driver()
	if d, ok := drv.(desktop.Driver); ok {
		w.window = d.CreateSplashWindow()
	} else {
		w.window = w.app.NewWindow("")
	}

	w.window.SetFixedSize(true)
	w.window.Resize(fyne.NewSize(300, 300))
	w.window.CenterOnScreen()
	w.window.SetPadded(false)

	w.window.SetContent(w.buildUI())

	go func() {
		time.Sleep(500 * time.Millisecond)
		w.window.RequestFocus()
		w.runUpdate()
	}()

	w.window.ShowAndRun()

}

func (w *Window) buildUI() fyne.CanvasObject {

	logoRes := fyne.NewStaticResource("logo.png", w.logo)
	logo := canvas.NewImageFromResource(logoRes)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(140, 50))

	w.spinner = widgets.NewSpinner(60, colours.Accent)
	w.spinner.Start()

	w.status = widget.NewLabel("preparing update...")
	w.status.Alignment = fyne.TextAlignCenter

	content := container.NewVBox(
		container.NewStack(
			container.NewPadded(layout.NewSpacer()),
		),
		layout.NewSpacer(),
		container.NewCenter(logo),
		layout.NewSpacer(),
		container.NewCenter(w.spinner),
		layout.NewSpacer(),
		container.NewCenter(w.status),
		layout.NewSpacer(),
	)

	return content

}

func (w *Window) setStatus(msg string) {
	w.status.SetText(msg)

	if w.fileLog != nil {
		timeStamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(w.fileLog, "[%s] %s\n", timeStamp, msg)
	}
}

func (w *Window) runUpdate() {
	exe := filepath.Join(w.dir, w.exeName)

	w.setStatus("stopping running instances...")
	utils.KillInstances(w.exeName, w.dir)
	utils.KillInstances("rpc.exe", w.dir)

	w.setStatus("preparing...")
	time.Sleep(1 * time.Second)

	if _, err := os.Stat(w.dir); err == nil {
		w.setStatus("removing old files...")
		if err := os.RemoveAll(w.dir); err != nil {
			w.showError(fmt.Sprintf("failed: %v", err))
			return
		}
	}

	if err := os.MkdirAll(w.dir, 0755); err != nil {
		w.showError(fmt.Sprintf("failed: %v", err))
		return
	}

	w.setStatus("extracting files...")
	if err := utils.ExtractZip(w.zipData, w.dir); err != nil {
		w.showError(fmt.Sprintf("failed: %v", err))
		return
	}

	if _, err := os.Stat(exe); os.IsNotExist(err) {
		w.showError("app not found in package")
		return
	}

	w.setStatus("setting registry key...")
	if err := utils.WriteRegistry(w.appName, w.dir); err != nil {
		w.showError(fmt.Sprintf("failed: %v", err))
		return
	}

	w.setStatus("creating shortcuts...")
	startMenu := filepath.Join(os.Getenv("ProgramData"), "Microsoft", "Windows", "Start Menu", "Programs")
	utils.CreateShortcut(filepath.Join(startMenu, w.appName+".lnk"), exe, w.dir, w.appName)

	desktop := filepath.Join(os.Getenv("PUBLIC"), "Desktop")
	utils.CreateShortcut(filepath.Join(desktop, w.appName+".lnk"), exe, w.dir, w.appName)

	w.setStatus("launching...")
	time.Sleep(300 * time.Millisecond)

	cmd := exec.Command(exe)
	cmd.Dir = w.dir
	cmd.Start()

	time.Sleep(500 * time.Millisecond)
	w.spinner.Stop()
	w.window.Close()
}

func (w *Window) showError(msg string) {
	w.setStatus("error: " + msg)
	time.Sleep(3 * time.Second)
	w.spinner.Stop()
	w.window.Close()
}
