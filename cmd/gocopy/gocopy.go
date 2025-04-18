package main

import (
	"context"
	"flag"
	"io"
	"os"
	"time"

	"github.com/gopasspw/clipboard"
)

func main() {
	ctx := context.Background()
	timeout := flag.Duration("t", 0, "Erase clipboard after timeout.  Durations are specified like \"20s\" or \"2h45m\".  0 (default) means never erase.")
	flag.Parse()

	out, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	if err := clipboard.WriteAll(ctx, out); err != nil {
		panic(err)
	}

	if timeout != nil && *timeout > 0 {
		<-time.After(*timeout)
		var text string
		text, err = clipboard.ReadAllString(ctx)
		if err != nil {
			os.Exit(1)
		}
		if text == string(out) {
			err = clipboard.WriteAllString(ctx, "")
		}
	}
	if err != nil {
		os.Exit(1)
	}
}
