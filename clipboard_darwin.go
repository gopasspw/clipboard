// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin
// +build darwin

package clipboard

import (
	"os/exec"
)

var (
	pasteCmdArgs = "pbpaste"
	copyCmdArgs  = "pbcopy"
)

func getPasteCommand() *exec.Cmd {
	return exec.Command(pasteCmdArgs)
}

func getCopyCommand() *exec.Cmd {
	return exec.Command(copyCmdArgs)
}

func readAll() ([]byte, error) {
	pasteCmd := getPasteCommand()
	out, err := pasteCmd.Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func writeAll(text []byte, _ bool) error {
	copyCmd := getCopyCommand()
	in, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}

	return copyCmd.Wait()
}
