# Clipboard access for Go

[![GoDoc](https://godoc.org/github.com/gopasspw/clipboard?status.svg)](http://godoc.org/github.com/gopasspw/clipboard)

Provide copying and pasting to the Clipboard for Go.

## Installation

    go get github.com/gopasspw/clipboard

## Platforms

We aim to support all platforms that Go supports, as long as they provide a clipboard.

The following platforms are know to work:

* Darwin (macOS)
* Windows
* Linux, Unix (requires 'xclip', 'xsel' or 'wl-clipboard' commands to be installed)

## Commands

paste shell command:

    go get github.com/gopasspw/clipboard/cmd/gopaste
    # example:
    gopaste > document.txt

copy shell command:

    go get github.com/gopasspw/clipboard/cmd/gocopy
    # example:
    cat document.txt | gocopy

## License and Credit

This package is licensed under the [BSD 3-Clause License](https://opensource.org/licenses/BSD-3-Clause).
It is a detached fork of github.com/atotto/clipboard. Since the origional repository has been inactive
for over 10 years, we decided to fork it and add some features and bug fixes.

This repository is maintained to support the needs of [gopass](https://github.com/gopasspw/gopass)
password manager. It is not an exact drop-in replacement, but
we encourage everyone to use it as a replacement for the original package.
