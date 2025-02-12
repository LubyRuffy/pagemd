package pagecontent

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchPage(t *testing.T) {
	// 创建一个文件服务器
	fileServer := http.FileServer(http.Dir("../../testdata"))

	// 创建一个测试服务器
	server := httptest.NewServer(fileServer)
	defer server.Close()

	requested := false
	_, err := fetchPageHTMLHeadless(server.URL+"/img.html", true, func(url string, content []byte) {
		t.Logf("URL: %s", url)

		assert.True(t, len(content) > 0)
		// check content is png image
		assert.Equal(t, "image/png", http.DetectContentType(content))

		if strings.Contains(url, "/img.png") {
			requested = true
		}
	})
	assert.NoError(t, err)
	assert.True(t, requested)
}
