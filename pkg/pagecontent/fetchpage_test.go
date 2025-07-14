package pagecontent

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchPage(t *testing.T) {
	// 创建一个文件服务器
	fileServer := http.FileServer(http.Dir("../../testdata"))

	// 创建一个测试服务器
	server := httptest.NewServer(fileServer)
	defer server.Close()

	requested := false
	_, err := fetchPage(server.URL+"/img.html", true)
	assert.NoError(t, err)
	assert.True(t, requested)
}
