package v1

import (
	"fmt"
	"github.com/LubyRuffy/pagemd/pkg/aitrans"
	"github.com/LubyRuffy/pagemd/pkg/pagecontent"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Bind(rg *gin.RouterGroup) {
	rg.POST("/fetch", func(c *gin.Context) {
		var req struct {
			URL string `json:"url"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		ci, err := pagecontent.NewAnalysis(pagecontent.WithURL(req.URL)).ExtractMainContent()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, ci)
		return
	})

	rg.GET(`/url/*any`, func(c *gin.Context) {
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
	rg.GET(`/translate/*any`, func(c *gin.Context) {
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
}
