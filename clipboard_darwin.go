// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin
// +build darwin

package clipboard

import (
	"context"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
)

var (
	pasteCmdArgs  = "pbpaste"
	copyCmdArgs   = "pbcopy"
	osascriptPath = "osascript"
)

func readAll(ctx context.Context) ([]byte, error) {
	pasteCmd := exec.CommandContext(ctx, pasteCmdArgs)
	out, err := pasteCmd.Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func writeAll(ctx context.Context, text []byte, secret bool) error {
	if secret {
		// Use osascript to set the clipboard contents, if that fails
		// just use pbcopy as a fallback.
		if err := copyViaOsascript(ctx, string(text)); err == nil {
			return nil
		}
	}

	// If osascript fails or we're not dealing with a secret, use pbcopy.
	copyCmd := exec.CommandContext(ctx, copyCmdArgs)
	in, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write(text); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}

	return copyCmd.Wait()
}

func copyViaOsascript(ctx context.Context, password string) error {
	// Base64-encode the password so it can be safely embedded in the AppleScript
	// source without any special escaping — the base64 alphabet [A-Za-z0-9+/=]
	// contains no characters that are special inside AppleScript string literals.
	b64 := base64.StdEncoding.EncodeToString([]byte(password))

	// The entire script is piped to "osascript -" via stdin, which avoids the
	// fd#3 inheritance issue introduced in macOS 15 (Sequoia) where subprocesses
	// spawned by "do shell script" no longer reliably inherit extra file
	// descriptors from the osascript process.
	script := fmt.Sprintf(`use framework "Foundation"
use scripting additions
set pw to do shell script "printf '%%s' '%s' | base64 -D"
set type to current application's NSPasteboardTypeString
set pb to current application's NSPasteboard's generalPasteboard()
pb's clearContents()
pb's setString:"" forType:"org.nspasteboard.ConcealedType"
pb's setString:pw forType:type
`, b64)

	cmd := exec.CommandContext(ctx, osascriptPath, "-")
	cmd.Stdin = strings.NewReader(script)

	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("osascript: %w", err)
	}

	// osascript prints the return value of the last expression; setString:forType:
	// returns a BOOL which AppleScript renders as "true" on success.
	if string(out) != "true\n" {
		return fmt.Errorf("osascript failed to set password: %s", string(out))
	}

	return nil
}

func unsupported() bool {
	return false
}
