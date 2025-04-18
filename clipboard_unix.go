// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows && !darwin && !plan9
// +build !windows,!darwin,!plan9

package clipboard

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	Primary bool
	trimDOS bool

	pasteCmdArgs   []string
	copyCmdArgs    []string
	copySecretArgs []string

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

	missingCommands = errors.New("No clipboard utilities available. Please install xsel, xclip, wl-clipboard or Termux:API add-on for termux-clipboard-get/set.")
)

func init() {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		pasteCmdArgs = wlpasteArgs
		copyCmdArgs = wlcopyArgs
		copySecretArgs = append(wlcopyArgs, "-o", "--type", "x-kde-passwordManagerHint/secret")

		if _, err := exec.LookPath(wlcopy); err == nil {
			if _, err := exec.LookPath(wlpaste); err == nil {
				return
			}
		}
	}

	pasteCmdArgs = xclipPasteArgs
	copyCmdArgs = xclipCopyArgs

	if _, err := exec.LookPath(xclip); err == nil {
		return
	}

	pasteCmdArgs = xselPasteArgs
	copyCmdArgs = xselCopyArgs

	if _, err := exec.LookPath(xsel); err == nil {
		return
	}

	pasteCmdArgs = termuxPasteArgs
	copyCmdArgs = termuxCopyArgs

	if _, err := exec.LookPath(termuxClipboardSet); err == nil {
		if _, err := exec.LookPath(termuxClipboardGet); err == nil {
			return
		}
	}

	pasteCmdArgs = powershellExePasteArgs
	copyCmdArgs = clipExeCopyArgs
	trimDOS = true

	if _, err := exec.LookPath(clipExe); err == nil {
		if _, err := exec.LookPath(powershellExe); err == nil {
			return
		}
	}

	Unsupported = true
}

func getPasteCommand() *exec.Cmd {
	if Primary {
		pasteCmdArgs = pasteCmdArgs[:1]
	}
	return exec.Command(pasteCmdArgs[0], pasteCmdArgs[1:]...)
}

func getCopyCommand() *exec.Cmd {
	if Primary {
		copyCmdArgs = copyCmdArgs[:1]
	}
	return exec.Command(copyCmdArgs[0], copyCmdArgs[1:]...)
}

func getCopySecretCommand() *exec.Cmd {
	if len(copySecretArgs) < 1 {
		copySecretArgs = copyCmdArgs
	}
	if Primary {
		copySecretArgs = copySecretArgs[:1]
	}
	return exec.Command(copySecretArgs[0], copySecretArgs[1:]...)
}

func readAll() ([]byte, error) {
	if Unsupported {
		return nil, missingCommands
	}

	pasteCmd := getPasteCommand()
	// capture errors
	eOut := &bytes.Buffer{}
	pasteCmd.Stderr = eOut

	out, err := pasteCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w. Output: %s", err, eOut.String())
	}
	result := out
	if trimDOS && len(result) > 1 {
		result = result[:len(result)-2]
	}
	return result, nil
}

func writeAll(text []byte, secret bool) error {
	if Unsupported {
		return missingCommands
	}
	var copyCmd *exec.Cmd
	if secret {
		copyCmd = getCopySecretCommand()
	} else {
		copyCmd = getCopyCommand()
	}
	// capture errors
	eOut := &bytes.Buffer{}
	copyCmd.Stderr = eOut

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
		return fmt.Errorf("failed to wait for command: %w. Output: %s", err, eOut.String())
	}
	return nil
}
