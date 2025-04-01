package pagecontent

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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

func TestExtractMainContent_f(t *testing.T) {
	// 创建一个文件服务器
	fileServer := http.FileServer(http.Dir("../../testdata"))
	// 创建一个测试服务器
	server := httptest.NewServer(fileServer)
	defer server.Close()
	url := server.URL + "/f.html"

	//url := "https://www.huawei.com/en/annual-report"

	ci, err := NewAnalysis(WithURL(url), WithDebug(true), WithHeadless(true)).ExtractMainContent()
	assert.NoError(t, err)
	assert.Contains(t, ci.Markdown, "annual_report_2024_en.pdf")
}
