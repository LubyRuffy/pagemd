package pagecontent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Contains(t, ci.Markdown, "http://127.0.0.1/a/b/a.jpg")

	ci, err = NewAnalysis(WithURL("http://127.0.0.1/a/b/c.html"),
		WithHTML(`<div><img src="/a.jpg"/>this is a text. this is a text. this is a text. `+
			`this is a text. this is a text. this is a text. this is a text. this is a text. this is a text. `+
			`this is a text. this is a text. this is a text. this is a text. this is a text. this is a text. `+
			`this is a text. this is a text. this is a text. this is a text. this is a text. this is a text. `+
			`this is a text. this is a text. this is a text. this is a text. this is a text. this is a text. <div>`)).ExtractMainContent()
	assert.NoError(t, err)
	assert.Contains(t, ci.Markdown, "http://127.0.0.1/a.jpg")

	//ci, err = NewAnalysis(WithURL("https://www.cnblogs.com/fanghan/p/13075290.html")).ExtractMainContent()
	//assert.NoError(t, err)
	//assert.Contains(t, ci.Markdown, "http://127.0.0.1/a.jpg")
}

func testFile(t *testing.T, file string, f func(ci *ContentInfo)) {
	// 创建一个文件服务器
	fileServer := http.FileServer(http.Dir("../../testdata"))
	// 创建一个测试服务器
	server := httptest.NewServer(fileServer)
	defer server.Close()
	url := server.URL + "/" + file
	ci, err := NewAnalysis(WithURL(url), WithDebug(true), WithHeadless(true)).ExtractMainContent()
	assert.NoError(t, err)
	f(ci)
}

func TestExtractMainContent_a(t *testing.T) {
	testFile(t, "a.html", func(ci *ContentInfo) {
		assert.Contains(t, ci.Markdown, "This blog post presents our fuzzer for")
		assert.Contains(t, ci.Markdown, "Blog post is published.")
	})
}

func TestExtractMainContent_c(t *testing.T) {
	testFile(t, "c.html", func(ci *ContentInfo) {
		assert.Contains(t, ci.Markdown, "I used to picture my dream job as this: I work for an early-stage startup that recently raised a round.")
		assert.Contains(t, ci.Markdown, "pulling you towards the messy codebase.")
	})
}

func TestExtractMainContent_d(t *testing.T) {
	testFile(t, "d.html", func(ci *ContentInfo) {
		assert.Contains(t, ci.Markdown, "Ollama now supports structured outputs making it possible to constrain a model")
		assert.Contains(t, ci.Markdown, "Additional format support beyond JSON schema")
	})
}

func TestExtractMainContent_e(t *testing.T) {
	testFile(t, "e.html", func(ci *ContentInfo) {
		assert.Contains(t, ci.Markdown, "Optimising and Visualising Go Tests Parallelism")
		assert.Contains(t, ci.Markdown, "There are more ways to make your tests more useful")
	})
}

func TestExtractMainContent_f(t *testing.T) {
	testFile(t, "f.html", func(ci *ContentInfo) {
		assert.Contains(t, ci.Markdown, "annual_report_2024_en.pdf")
	})
}

func TestExtractMainContent_g(t *testing.T) {
	testFile(t, "g.html", func(ci *ContentInfo) {
		assert.Contains(t, ci.Markdown, "2001年")
	})
}
