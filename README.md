# Clipboard for Go

[![GoDoc](https://godoc.org/github.com/gopasspw/clipboard?status.svg)](http://godoc.org/github.com/gopasspw/clipboard)

Provide copying and pasting to the Clipboard for Go.

Build:

    go get github.com/gopasspw/clipboard

Platforms:

* OSX
* Windows 7 (probably work on other Windows)
* Linux, Unix (requires 'xclip', 'xsel' or 'wl-clipboard' commands to be installed)

Document:

* http://godoc.org/github.com/gopasspw/clipboard

Notes:

* Text string only
* UTF-8 text encoding only (no conversion)

## Commands

paste shell command:

    go get github.com/gopasspw/clipboard/cmd/gopaste
    # example:
    gopaste > document.txt

copy shell command:

    go get github.com/gopasspw/clipboard/cmd/gocopy
    # example:
    cat document.txt | gocopy
