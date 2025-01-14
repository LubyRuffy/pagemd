package aitag

import (
	"context"
	"fmt"
	"github.com/LubyRuffy/pagemd/pkg/pagecontent"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestNew(t *testing.T) {
	testF := func(f string, tags []string) {
		h, err := os.ReadFile(filepath.Join("../../testdata", f))
		assert.NoError(t, err)

		ci, err := pagecontent.NewAnalysis(pagecontent.WithURL("http://127.0.0.1/a/b/c.html"),
			pagecontent.WithHTML(string(h))).ExtractMainContent()

		llmTags, err := New().Tag(context.Background(),
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

	testF("a.html", []string{"cybersecurity", "fuzzing", "technical", "vulnerability exploitation"})
	testF("b.html", []string{"cybersecurity", "security tools", "technical"})
	testF("c.html", []string{"programming", "technical", "software engineering"})
	testF("d.html", []string{"AI", "javascript", "python", "ollama", "AI tools", "technical", "programming"})
	testF("e.html", []string{"golang", "technical", "programming"})
}
