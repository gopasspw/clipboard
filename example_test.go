package clipboard_test

import (
	"context"
	"fmt"

	"github.com/gopasspw/clipboard"
)

func Example() {
	clipboard.WriteAllString(context.TODO(), "日本語") //nolint:errcheck
	text, _ := clipboard.ReadAllString(context.TODO()) //nolint:errcheck
	fmt.Println(text)

	// Output:
	// 日本語
}
