package pagecontent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractMainContent(t *testing.T) {
	// 测试支持图片的相对路径
	ci, err := NewAnalysis(WithURL("http://127.0.0.1/a/b/c.html"),
		WithHTML(`<div><img src="a.jpg"/>this is a text. this is a text. this is a text. `+
			`this is a text. this is a text. this is a text. this is a text. this is a text. this is a text. `+
			`this is a text. this is a text. this is a text. this is a text. this is a text. this is a text. `+
			`this is a text. this is a text. this is a text. this is a text. this is a text. this is a text. `+
			`this is a text. this is a text. this is a text. this is a text. this is a text. this is a text. <div>`)).ExtractMainContent()
	assert.NoError(t, err)
	assert.Contains(t, ci.ContentMarkdown, "http://127.0.0.1/a/b/a.jpg")

	ci, err = NewAnalysis(WithURL("http://127.0.0.1/a/b/c.html"),
		WithHTML(`<div><img src="/a.jpg"/>this is a text. this is a text. this is a text. `+
			`this is a text. this is a text. this is a text. this is a text. this is a text. this is a text. `+
			`this is a text. this is a text. this is a text. this is a text. this is a text. this is a text. `+
			`this is a text. this is a text. this is a text. this is a text. this is a text. this is a text. `+
			`this is a text. this is a text. this is a text. this is a text. this is a text. this is a text. <div>`)).ExtractMainContent()
	assert.NoError(t, err)
	assert.Contains(t, ci.ContentMarkdown, "http://127.0.0.1/a.jpg")
}
