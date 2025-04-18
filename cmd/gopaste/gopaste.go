package main

import (
	"context"
	"fmt"

	"github.com/gopasspw/clipboard"
)

func main() {
	text, err := clipboard.ReadAllString(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Print(text)
}
