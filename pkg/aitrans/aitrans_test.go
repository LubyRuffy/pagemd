package aitrans

import (
	"context"
	"fmt"
	"github.com/LubyRuffy/pagemd/pkg/llm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	cfg, err := llm.Load("../../config.yaml")
	assert.NoError(t, err)
	New(cfg).TranslateToChinese(context.Background(),
		"hello\n```python\na='this is a test'```",
		func(s string) {
			fmt.Printf("%s", s)
		})
}
