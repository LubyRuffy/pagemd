package maindiv

import (
	"errors"
	"github.com/go-rod/rod/lib/launcher"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// ErrNotHtml is an error returned when the content type of the response is not HTML.
var ErrNotHtml = errors.New("not html") // not html error

// FetchPageHTMLHeadless fetches the HTML content of a page using headless browser.
//
// Parameters:
//   - u: The URL of the page to be fetched.
//   - debug: A boolean flag indicating whether to run the browser in debug mode with DevTools and tracing enabled.
//
// Returns:
//   - string: The HTML content of the page.
//   - error: An error if the fetching process fails.
func FetchPageHTMLHeadless(u string, debug bool) (string, error) {
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
	page := browser.MustPage(u).Timeout(10 * time.Second).MustWaitStable()
	return page.HTML()
}

// FetchPageHTML fetches the HTML content of a page from the given URL.
//
// Parameters:
//   - u: The URL of the page to be fetched.
//   - headless: A boolean flag indicating whether to use headless browser mode.
//
// Returns:
//   - string: The HTML content of the page.
//   - error: An error if the fetching process fails.
func FetchPageHTML(u string, headless bool) (string, error) {
	// Implement fetching the page content from the given URL
	if headless {
		return FetchPageHTMLHeadless(u, true)
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
