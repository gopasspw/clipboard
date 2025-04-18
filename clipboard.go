// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package clipboard read/write on clipboard
package clipboard

import "context"

// ReadAll reads from the clipboard.
func ReadAll(ctx context.Context) ([]byte, error) {
	return readAll(ctx)
}

// ReadAllString reads a string from the clipboard.
func ReadAllString(ctx context.Context) (string, error) {
	text, err := readAll(ctx)
	if err != nil {
		return "", err
	}

	return string(text), nil
}

// WriteAll writes a string to the clipboard.
func WriteAll(ctx context.Context, text []byte) error {
	return writeAll(ctx, text, false)
}

// WriteAllString writes a string to the clipboard.
func WriteAllString(ctx context.Context, text string) error {
	return writeAll(ctx, []byte(text), false)
}

// WritePassword writes a password to the clipboard.
func WritePassword(ctx context.Context, text []byte) error {
	return writeAll(ctx, text, true)
}

// IsUnsupported returns true if the current platform is not supported.
func IsUnsupported() bool {
	return unsupported()
}
