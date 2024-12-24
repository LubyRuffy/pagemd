package aitrans

import (
	"context"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	New().TranslateToChinese(context.Background(),
		"hello\n```python\na='this is a test'```",
		func(s string) {
			fmt.Printf("%s", s)
		})
}
