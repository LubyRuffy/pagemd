package aitag

import (
	"context"
	"fmt"
	"github.com/LubyRuffy/pagemd/pkg/llm"
	"github.com/LubyRuffy/pagemd/pkg/pagecontent"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
)

func TestNew(t *testing.T) {
	cfg, err := llm.Load("../../config.yaml")
	assert.NoError(t, err)

	// 创建一个文件服务器
	fileServer := http.FileServer(http.Dir("../../testdata"))

	// 创建一个测试服务器
	server := httptest.NewServer(fileServer)
	defer server.Close()

	testF := func(f string, tags []string) {
		ci, err := pagecontent.NewAnalysis(pagecontent.WithURL(server.URL+f), pagecontent.WithDebug(true)).
			ExtractMainContent()
		assert.NoError(t, err)

		llmTags, err := New(cfg).Tag(context.Background(),
			ci.Markdown,
			func(s string) {
				fmt.Printf("%s", s)
			})
		assert.NoError(t, err)
		sort.Strings(llmTags)
		sort.Strings(tags)

		for i := range tags {
			assert.Contains(t, llmTags, tags[i])
		}
		//assert.Equal(t, tags, llmTags)
	}

	testF("/a.html", []string{"cybersecurity", "fuzzing", "technical", "vulnerability exploitation"})
	testF("/b.html", []string{"cybersecurity", "security tools", "technical"})
	testF("/c.html", []string{"programming", "technical", "software engineering"})
	testF("/d.html", []string{"AI", "javascript", "python", "ollama", "AI tools", "technical", "programming"})
	testF("/e.html", []string{"golang", "technical", "programming"})
}
