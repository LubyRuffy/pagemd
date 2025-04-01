package pagecontent

import (
	"errors"
	"fmt"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// ErrNotHtml is an error returned when the content type of the response is not HTML.
var ErrNotHtml = errors.New("not html") // not html error

// fetchPageHTMLHeadless fetches the HTML content of a page using headless browser.
//
// Parameters:
//   - u: The URL of the page to be fetched.
//   - debug: A boolean flag indicating whether to run the browser in debug mode with DevTools and tracing enabled.
//
// Returns:
//   - string: The HTML content of the page.
//   - error: An error if the fetching process fails.
func fetchPageHTMLHeadless(u string, debug bool, onImg func(string, []byte)) (string, error) {
	var browser *rod.Browser
	if debug {
		// Headless runs the browser on foreground, you can also use flag "-rod=show"
		// Devtools opens the tab in each new tab opened automatically
		l := launcher.New().
			Headless(false).
			Devtools(true)
		defer l.Cleanup()

		// Trace shows verbose debug information for each action executed
		// SlowMotion is a debug related function that waits 2 seconds between
		// each action, making it easier to inspect what your code is doing.
		browser = rod.New().
			ControlURL(l.MustLaunch()).
			Trace(true).
			SlowMotion(2 * time.Second).
			MustConnect()

		// ServeMonitor plays screenshots of each tab. This feature is extremely
		// useful when debugging with headless mode.
		// You can also enable it with flag "-rod=monitor"
		//launcher.Open(browser.ServeMonitor(""))
	} else {
		// Launch a new browser with default options, and connect to it.
		browser = rod.New().MustConnect()
	}

	defer browser.MustClose()

	// Create a new page
	page := browser.MustPage()

	if onImg != nil {
		// 设置网络事件监听器
		router := page.HijackRequests()
		router.MustAdd("*", func(ctx *rod.Hijack) {
			// 我们只对图片请求感兴趣
			if ctx.Request.Type() == proto.NetworkResourceTypeImage {
				// 拦截请求
				go func() {
					// 获取响应
					err := ctx.LoadResponse(http.DefaultClient, true)
					if err != nil {
						log.Println("LoadResponse error:", err)
						return // Don't panic; it might just be a single image failure.
					}

					// Check if the response is valid.
					if ctx.Response == nil { // ctx.Response is available after LoadResponse
						log.Println("Response is nil for", ctx.Request.URL())
						return
					}

					// Check the status code.
					if ctx.Response.RawResponse.StatusCode != http.StatusOK {
						log.Printf("Image %s status code: %d\n", ctx.Request.URL(), ctx.Response.RawResponse.StatusCode)
						return // Only process 200 OK images.
					}

					if len(ctx.Response.Body()) == 0 {
						log.Printf("Image %s is empty\n", ctx.Request.URL())
						return // Skip empty images.
					}

					// todo：替换图片
					onImg(ctx.Request.URL().String(), []byte(ctx.Response.Body()))
				}()

				ctx.ContinueRequest(&proto.FetchContinueRequest{}) // 继续原始请求

			} else {
				ctx.ContinueRequest(&proto.FetchContinueRequest{}) //不拦截其他资源
			}
		})
		// 开始监听
		go router.Run()     //必须在单独的 goroutine 中运行，否则会阻塞
		defer router.Stop() // 确保停止劫持
	}

	page = page.Timeout(10 * time.Second).MustNavigate(u).MustWaitStable().CancelTimeout()

	removeInvisibleDiv := false
	if removeInvisibleDiv {
		// 使用JavaScript直接在页面上执行删除不可见div的操作
		// 这样避免了在Go代码中遍历元素可能导致的引用问题
		script := `() => 
    (function() {
        // 获取所有div元素
        const divs = Array.from(document.querySelectorAll('div'));
        let removedCount = 0;
        
        // 遍历并检查每个div
        for (const div of divs) {
            const style = window.getComputedStyle(div);
            
            // 检查元素是否不可见
            if (style.display === 'none' || 
                style.visibility === 'hidden' || 
                style.opacity === '0' || 
                (div.offsetWidth === 0 && div.offsetHeight === 0)) {
                
                // 如果元素有父节点，则移除它
                if (div.parentNode) {
                    div.parentNode.removeChild(div);
                    removedCount++;
                }
            }
        }
        
        return removedCount;
    })()
    `

		result := page.MustEval(script)
		fmt.Printf("已删除 %v 个不可见的div元素\n", result.Int())

		nodes := page.MustElements("div")
		for _, node := range nodes {
			if v, err := node.Visible(); err == nil && !v {
				node.MustRemove()
			}
		}
	}
	return page.HTML()
}

// fetchPageHTML fetches the HTML content of a page from the given URL.
//
// Parameters:
//   - u: The URL of the page to be fetched.
//   - headless: A boolean flag indicating whether to use headless browser mode.
//
// Returns:
//   - string: The HTML content of the page.
//   - error: An error if the fetching process fails.
func fetchPageHTML(u string, headless bool, debug bool) (string, error) {
	// Implement fetching the page content from the given URL
	if headless {
		return fetchPageHTMLHeadless(u, debug, nil)
	}

	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 判断header中是否指明了html
	if contentType := resp.Header.Get("Content-Type"); !strings.Contains(contentType, "text/html") {
		return "", ErrNotHtml
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
