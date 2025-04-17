package main

import (
	"flag"
	"io"
	"os"
	"time"

	"github.com/gopasspw/clipboard"
)

func main() {
	timeout := flag.Duration("t", 0, "Erase clipboard after timeout.  Durations are specified like \"20s\" or \"2h45m\".  0 (default) means never erase.")
	flag.Parse()

	out, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	if err := clipboard.WriteAll(string(out)); err != nil {
		panic(err)
	}

	if timeout != nil && *timeout > 0 {
		<-time.After(*timeout)
		var text string
		text, err = clipboard.ReadAll()
		if err != nil {
			os.Exit(1)
		}
		if text == string(out) {
			err = clipboard.WriteAll("")
		}
	}
	if err != nil {
		os.Exit(1)
	}
}
