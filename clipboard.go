// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package clipboard read/write on clipboard
package clipboard

// ReadAll reads from the clipboard.
func ReadAll() ([]byte, error) {
	return readAll()
}

// ReadAllString reads a string from the clipboard.
func ReadAllString() (string, error) {
	text, err := readAll()
	if err != nil {
		return "", err
	}
	return string(text), nil
}

// WriteAll write string to clipboard.
func WriteAll(text []byte) error {
	return writeAll(text, false)
}

// WriteAllString writes a string to the clipboard.
func WriteAllString(text string) error {
	return writeAll([]byte(text), false)
}

// WritePassword writes a password to the clipboard.
func WritePassword(text []byte) error {
	return writeAll(text, true)
}

// Unsupported might be set true during clipboard init, to help callers decide
// whether or not to offer clipboard options.
var Unsupported bool
