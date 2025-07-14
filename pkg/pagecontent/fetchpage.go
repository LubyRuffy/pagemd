package pagecontent

import (
	"github.com/LubyRuffy/pageviewer"
)

var (
	b *pageviewer.Browser
)

// fetchPage fetches the Article content of a page using headless browser.
//
// Parameters:
//   - u: The URL of the page to be fetched.
//   - debug: A boolean flag indicating whether to run the browser in debug mode with DevTools and tracing enabled.
//
// Returns:
//   - ReadabilityArticleWithMarkdown: The Article information of the page.
//   - error: An error if the fetching process fails.
func fetchPage(u string, debug bool) (*pageviewer.ReadabilityArticleWithMarkdown, error) {
	if b == nil {
		var err error
		b, err = pageviewer.NewBrowser(pageviewer.WithDebug(debug))
		if err != nil {
			return nil, err
		}
	}

	result, err := b.ReadabilityArticle(u)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
