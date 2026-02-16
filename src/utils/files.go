package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func ExtractZip(data []byte, dest string) error {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("zip open failed: %w", err)
	}

	for _, f := range r.File {
		target := filepath.Join(dest, f.Name)

		if !isValidPath(target, dest) {
			return fmt.Errorf("invalid path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(target, 0755)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		src, err := f.Open()
		if err != nil {
			return err
		}

		dst, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			src.Close()
			return err
		}

		_, err = io.Copy(dst, src)
		src.Close()
		dst.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func isValidPath(path, base string) bool {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return false
	}
	return !filepath.IsAbs(rel) && rel != ".." && !filepath.HasPrefix(rel, ".."+string(filepath.Separator))
}

func CreateShortcut(shortcut, target, workDir, desc string) error {
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	defer ole.CoUninitialize()

	shell, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer shell.Release()

	wshell, err := shell.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshell.Release()

	v, err := oleutil.CallMethod(wshell, "CreateShortcut", shortcut)
	if err != nil {
		return err
	}

	sc := v.ToIDispatch()
	defer sc.Release()

	oleutil.PutProperty(sc, "TargetPath", target)
	oleutil.PutProperty(sc, "WorkingDirectory", workDir)
	oleutil.PutProperty(sc, "Description", desc)
	oleutil.PutProperty(sc, "IconLocation", target+",0")

	_, err = oleutil.CallMethod(sc, "Save")
	return err
}