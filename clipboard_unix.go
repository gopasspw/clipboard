// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows && !darwin && !plan9
// +build !windows,!darwin,!plan9

package clipboard

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

const (
	xsel               = "xsel"
	xclip              = "xclip"
	powershellExe      = "powershell.exe"
	clipExe            = "clip.exe"
	wlcopy             = "wl-copy"
	wlpaste            = "wl-paste"
	termuxClipboardGet = "termux-clipboard-get"
	termuxClipboardSet = "termux-clipboard-set"
)

var (
	wrap *wrapper
	once sync.Once

	xselPasteArgs = []string{xsel, "--output", "--clipboard"}
	xselCopyArgs  = []string{xsel, "--input", "--clipboard"}

	xclipPasteArgs = []string{xclip, "-out", "-selection", "clipboard"}
	xclipCopyArgs  = []string{xclip, "-in", "-selection", "clipboard"}

	powershellExePasteArgs = []string{powershellExe, "Get-Clipboard"}
	clipExeCopyArgs        = []string{clipExe}

	wlpasteArgs = []string{wlpaste, "--no-newline"}
	wlcopyArgs  = []string{wlcopy}

	termuxPasteArgs = []string{termuxClipboardGet}
	termuxCopyArgs  = []string{termuxClipboardSet}

	errMissingCommands = errors.New("no clipboard utilities available. Please install xsel, xclip, wl-clipboard or Termux:API add-on for termux-clipboard-get/set")
)

type wrapper struct {
	trimDOS bool

	pasteCmdArgs   []string
	copyCmdArgs    []string
	copySecretArgs []string

	// supported is true if the clipboard is supported
	supported bool
}

func (w *wrapper) unsupported() bool {
	return !w.supported
}

func hasBinary(name string) bool {
	_, err := exec.LookPath(name)

	return err == nil
}

func newWrapper() *wrapper {
	w := &wrapper{}

	switch {
	case os.Getenv("WAYLAND_DISPLAY") != "" && hasBinary(wlcopy) && hasBinary(wlpaste):
		// Wayland
		w.pasteCmdArgs = wlpasteArgs
		w.copyCmdArgs = wlcopyArgs
		w.copySecretArgs = append(wlcopyArgs, "--type", "x-kde-passwordManagerHint/secret")
		w.supported = true
	case hasBinary(xclip):
		// X11 (or Wayland) with xclip
		w.pasteCmdArgs = xclipPasteArgs
		w.copyCmdArgs = xclipCopyArgs
		w.supported = true
	case hasBinary(xsel):
		// X11 (or Wayland) with xsel
		w.pasteCmdArgs = xselPasteArgs
		w.copyCmdArgs = xselCopyArgs
		w.supported = true
	case hasBinary(termuxClipboardSet) && hasBinary(termuxClipboardGet):
		// Termux
		w.pasteCmdArgs = termuxPasteArgs
		w.copyCmdArgs = termuxCopyArgs
		w.supported = true
	case hasBinary(clipExe) && hasBinary(powershellExe):
		// Powershell
		w.pasteCmdArgs = powershellExePasteArgs
		w.copyCmdArgs = clipExeCopyArgs
		w.trimDOS = true
		w.supported = true
	default:
		// Unsupported
		w.supported = false
	}

	return w
}

func getWrapper() *wrapper {
	once.Do(func() {
		wrap = newWrapper()
	})

	return wrap
}

func init() {
}

func readAll(ctx context.Context) ([]byte, error) {
	w := getWrapper()
	if w.unsupported() || ForceUnsupported {
		return nil, errMissingCommands
	}

	pasteCmd := exec.CommandContext(ctx, w.pasteCmdArgs[0], w.pasteCmdArgs[1:]...)
	// capture errors
	eOut := &bytes.Buffer{}
	pasteCmd.Stderr = eOut

	out, err := pasteCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w. Output: %s", err, eOut.String())
	}
	result := out
	if w.trimDOS && len(result) > 1 {
		result = result[:len(result)-2]
	}

	return result, nil
}

func writeAll(ctx context.Context, text []byte, secret bool) error {
	w := getWrapper()
	if w.unsupported() || ForceUnsupported {
		return errMissingCommands
	}

	copyCmd := exec.CommandContext(ctx, w.copyCmdArgs[0], w.copyCmdArgs[1:]...)
	if secret && len(w.copySecretArgs) > 0 {
		copyCmd = exec.CommandContext(ctx, w.copySecretArgs[0], w.copySecretArgs[1:]...)
	}

	// Note: Can not capture stderr from the copyCmd or it will hang.

	in, err := copyCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}
	if err := copyCmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}
	if _, err := in.Write(text); err != nil {
		return fmt.Errorf("failed to write to stdin: %w", err)
	}
	if err := in.Close(); err != nil {
		return fmt.Errorf("failed to close stdin: %w", err)
	}
	if err := copyCmd.Wait(); err != nil {
		return fmt.Errorf("failed to wait for command: %w", err)
	}

	return nil
}

func unsupported() bool {
	w := getWrapper()

	return w.unsupported()
}
