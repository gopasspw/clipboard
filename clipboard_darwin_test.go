// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin
// +build darwin

package clipboard

import (
	"testing"
)

// TestCopyViaOsascriptSuccess verifies that copyViaOsascript writes the
// expected value to the system clipboard and returns no error.
func TestCopyViaOsascriptSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping osascript integration test in short mode")
	}

	ctx := t.Context()
	const want = "s3cr3t-pass"

	if err := copyViaOsascript(ctx, want); err != nil {
		t.Fatalf("copyViaOsascript returned error: %v", err)
	}

	got, err := readAll(ctx)
	if err != nil {
		t.Fatalf("readAll failed: %v", err)
	}
	if string(got) != want {
		t.Errorf("clipboard = %q, want %q", got, want)
	}
}

// TestCopyViaOsascriptSpecialChars verifies that passwords containing
// characters that are special in AppleScript string literals (quotes,
// ampersands, angle brackets) survive the base64 round-trip correctly.
func TestCopyViaOsascriptSpecialChars(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping osascript integration test in short mode")
	}

	ctx := t.Context()
	const want = `p@$$w0rd "quoted" & <tags>`

	if err := copyViaOsascript(ctx, want); err != nil {
		t.Fatalf("copyViaOsascript returned error: %v", err)
	}

	got, err := readAll(ctx)
	if err != nil {
		t.Fatalf("readAll failed: %v", err)
	}
	if string(got) != want {
		t.Errorf("clipboard = %q, want %q", got, want)
	}
}

// TestWriteAllSecretFallback verifies that when osascript is unavailable,
// writeAll quietly falls back to pbcopy and still places the text on the
// clipboard without returning an error.
func TestWriteAllSecretFallback(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping osascript integration test in short mode")
	}

	// Point osascriptPath at a nonexistent binary so that copyViaOsascript
	// always fails and the pbcopy fallback path inside writeAll is exercised.
	orig := osascriptPath
	osascriptPath = "/nonexistent-osascript-binary"
	t.Cleanup(func() { osascriptPath = orig })

	ctx := t.Context()
	const want = "fallback-test-pass"

	if err := writeAll(ctx, []byte(want), true); err != nil {
		t.Fatalf("writeAll (secret, fallback) returned error: %v", err)
	}

	got, err := readAll(ctx)
	if err != nil {
		t.Fatalf("readAll failed: %v", err)
	}
	if string(got) != want {
		t.Errorf("clipboard = %q, want %q", got, want)
	}
}
