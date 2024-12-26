package api

import (
	"fmt"
	"github.com/LubyRuffy/pagemd/pkg/aitrans"
	"github.com/LubyRuffy/pagemd/pkg/pagecontent"
	"github.com/gin-gonic/gin"
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

	e.GET(`/url/*any`, func(c *gin.Context) {
		if targetURL, exists := c.Get("targetURL"); exists {
			ci, err := pagecontent.NewAnalysis(pagecontent.WithURL(targetURL.(string))).ExtractMainContent()
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(200, ci)
			return
		}
		c.String(http.StatusNotFound, "Not Found")
	})
	e.GET(`/translate/*any`, func(c *gin.Context) {
		if targetURL, exists := c.Get("targetURL"); exists {
			ci, err := pagecontent.NewAnalysis(pagecontent.WithURL(targetURL.(string))).ExtractMainContent()
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}
			a := aitrans.New()
			v, err := a.TranslateToChinese(c,
				ci.Markdown,
				func(s string) {
					fmt.Printf("%s", s)
				})
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			cnci := struct {
				*pagecontent.ContentInfo
				MarkdownCN string `json:"markdown_cn"`
			}{
				ContentInfo: ci,
				MarkdownCN:  v,
			}
			c.JSON(200, cnci)
			return
		}
		c.String(http.StatusNotFound, "Not Found")
	})
	return e.Run(addr)
}
