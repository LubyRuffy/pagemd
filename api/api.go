package api

import (
	v1 "github.com/LubyRuffy/pagemd/api/v1"
	"github.com/LubyRuffy/pagemd/web"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func URLPassthroughMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, prefix := range []string{"/url/", "/translate/"} {
			if strings.Contains(path, prefix+"http://") || strings.Contains(path, prefix+"https://") {
				// 提取并处理 URL
				targetURL := path[len(prefix):] // 去掉开头的斜杠
				c.Set("targetURL", targetURL)

				// 可以选择跳过后续的 Gin 路由，或者修改请求
				c.Next() // 继续执行后续处理
				return
			}
		}
		c.Next()
	}
}

func Start(addr string) error {
	e := gin.Default()

	e.Use(URLPassthroughMiddleware())
	rg := e.Group("/api/v1")
	v1.Bind(rg)

	e.StaticFS("/web", http.FS(web.WebFs))
	e.GET(`/`, func(c *gin.Context) {
		c.Redirect(301, "/web/index.html")
	})
	log.Printf("Server listening on http://%s\n", addr)
	return e.Run(addr)
}
