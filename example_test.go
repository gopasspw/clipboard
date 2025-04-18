package clipboard_test

import (
	"fmt"

	"github.com/gopasspw/clipboard"
)

func Example() {
	clipboard.WriteAllString("日本語") //nolint:errcheck
	text, _ := clipboard.ReadAll()  //nolint:errcheck
	fmt.Println(text)

	// Output:
	// 日本語
}
