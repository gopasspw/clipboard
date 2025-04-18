// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clipboard_test

import (
	"os"
	"runtime"
	"testing"

	. "github.com/gopasspw/clipboard"
)

func TestCopyAndPaste(t *testing.T) {
	if runtime.GOOS == "linux" && os.Getenv("WAYLAND_DISPLAY") == "" && os.Getenv("DISPLAY") == "" {
		t.Skip("Skipping test on Linux without Wayland or X11")
	}

	expected := "日本語"

	err := WriteAllString(expected)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := ReadAllString()
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Errorf("want %s, got %s", expected, actual)
	}
}

func TestMultiCopyAndPaste(t *testing.T) {
	if runtime.GOOS == "linux" && os.Getenv("WAYLAND_DISPLAY") == "" && os.Getenv("DISPLAY") == "" {
		t.Skip("Skipping test on Linux without Wayland or X11")
	}

	expected1 := "French: éèêëàùœç"
	expected2 := "Weird UTF-8: 💩☃"

	err := WriteAllString(expected1)
	if err != nil {
		t.Fatal(err)
	}

	actual1, err := ReadAllString()
	if err != nil {
		t.Fatal(err)
	}
	if actual1 != expected1 {
		t.Errorf("want %s, got %s", expected1, actual1)
	}

	err = WriteAllString(expected2)
	if err != nil {
		t.Fatal(err)
	}

	actual2, err := ReadAllString()
	if err != nil {
		t.Fatal(err)
	}
	if actual2 != expected2 {
		t.Errorf("want %s, got %s", expected2, actual2)
	}
}

func BenchmarkReadAll(b *testing.B) {
	for b.Loop() {
		ReadAll() //nolint:errcheck
	}
}

func BenchmarkWriteAll(b *testing.B) {
	text := "いろはにほへと"
	for b.Loop() {
		WriteAllString(text) //nolint:errcheck
	}
}
