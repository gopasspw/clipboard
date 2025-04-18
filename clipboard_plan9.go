// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9
// +build plan9

package clipboard

import (
	"context"
	"io"
	"os"
)

func readAll(_ context.Context) ([]byte, error) {
	f, err := os.Open("/dev/snarf")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	str, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return str, nil
}

func writeAll(_ context.Context, text []byte, _ bool) error {
	f, err := os.OpenFile("/dev/snarf", os.O_WRONLY, 0o666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(text)
	if err != nil {
		return err
	}

	return nil
}

func unsupported() bool {
	return false
}
