package pagecontent

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"unicode/utf8"
)

type TitleAuthorDate struct {
	Author string
	Title  string
	Date   string
}

// 查找两个字符串的最长公共子串
func longestCommonSubstring(s1, s2 string) string {
	r1 := []rune(s1)
	r2 := []rune(s2)
	m := len(r1)
	n := len(r2)
	dp := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		dp[i] = make([]int, n+1)
	}
	maxLength := 0
	endIndex := 0
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if r1[i-1] == r2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				if dp[i][j] > maxLength {
					maxLength = dp[i][j]
					endIndex = i
				}
			} else {
				dp[i][j] = 0
			}
		}
	}
	if maxLength > 0 {
		return string(r1[endIndex-maxLength : endIndex])
	}
	return ""
}

// 从多个候选标题中提取最准确的标题
func getAccurateText(texts []string) string {
	if len(texts) == 0 {
		return ""
	}
	if len(texts) == 1 {
		return texts[0]
	}
	texts = texts[:min(5, len(texts))]

	mostFrequentTexts := make(map[string]int, 0)
	for i := 0; i < len(texts); i++ {
		for j := i + 1; j < len(texts); j++ {
			text := longestCommonSubstring(texts[i], texts[j])
			if utf8.RuneCountInString(strings.Trim(text, " \r\t\n")) < 3 {
				continue
			}
			if v, ok := mostFrequentTexts[text]; !ok {
				mostFrequentTexts[text] = 1
			} else {
				mostFrequentTexts[text] = v + 1
			}

		}
	}

	var mostFrequentTextCount int
	var mostFrequentText string
	for k, v := range mostFrequentTexts {
		if v > mostFrequentTextCount {
			mostFrequentText = k
			mostFrequentTextCount = v
		} else if len(k) > len(mostFrequentText) {
			mostFrequentText = k
			mostFrequentTextCount = v
		}
	}
	return mostFrequentText
}

// ExtractTitleAuthorDate 提取title，author，date
func ExtractTitleAuthorDate(htmlContent string) (*TitleAuthorDate, error) {
	// title标签：
	// <title>Bluetooth Low Energy GATT Fuzzing - Quarkslab's blog</title>
	// <title>The Over-Engineering Pendulum | Three Dots Labs blog</title>
	// <title>Structured outputs · Ollama Blog</title>

	// meta property=og:title：
	// <meta property="og:title" content="Bluetooth Low Energy GATT Fuzzing"/>
	// <meta property="og:title" content="Structured outputs · Ollama Blog" />

	// meta name=twitter:title：
	// <meta name=twitter:title content="The Over-Engineering Pendulum">
	// <meta property="twitter:title" content="Structured outputs· Ollama Blog" />

	// article下面的第一个h1:
	// <article class=lg:col-7><h1 class="h2 mb-2">The Over-Engineering Pendulum</h1>
	// <article class="mx-auto flex flex-1 max-w-2xl w-full flex-col space-y-3 px-6 py-16 md:px-0"><h1 class="text-4xl font-semibold tracking-tight">Structured outputs</h1>

	// a rel="bookmark"：
	// <a href="./bluetooth-low-energy-gatt-fuzzing.html" rel="bookmark" title="Permalink to Bluetooth Low Energy GATT Fuzzing"> Bluetooth Low Energy  GATT Fuzzing</a>

	// div.class=title：
	// <div class="title" data-v-761cb514>HookCase：一款针对maxOS的逆向工程安全分析工具</div>

	// span.class=title：
	//<span class="title-span" data-v-761cb514>HookCase：一款针对maxOS的逆向工程安全分析工具</span>
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htmlContent)))
	if err != nil {
		return nil, err
	}

	trim := func(s string) string {
		s = strings.Trim(s, "\r\n\t ")
		for {
			// 替换多个空格
			if !strings.Contains(s, "  ") {
				break
			}
			s = strings.ReplaceAll(s, "  ", " ")
		}
		return s
	}

	var titles []string
	var authors []string
	doc.Find("*").Each(func(i int, selection *goquery.Selection) {
		switch strings.ToLower(goquery.NodeName(selection)) {
		case "title":
			if v := trim(selection.Text()); v != "" && !selection.Parent().Is("symbol") {
				titles = append(titles, v)
			}
		case "meta":
			// title
			if selection.AttrOr("property", "") == "og:title" ||
				selection.AttrOr("property", "") == "twitter:title" ||
				selection.AttrOr("name", "") == "og:title" ||
				selection.AttrOr("name", "") == "twitter:title" {
				if v := trim(selection.AttrOr("content", "")); v != "" {
					titles = append(titles, v)
				}
			}
			// author
			if selection.AttrOr("property", "") == "og:author" ||
				selection.AttrOr("property", "") == "author" ||
				selection.AttrOr("property", "") == "twitter:author" ||
				selection.AttrOr("property", "") == "article:author" ||
				selection.AttrOr("name", "") == "og:author" ||
				selection.AttrOr("name", "") == "author" ||
				selection.AttrOr("name", "") == "article:author" ||
				selection.AttrOr("name", "") == "twitter:author" {
				if v := trim(selection.AttrOr("content", "")); v != "" {
					authors = append(authors, v)
				}
			}
		case "article":
			if strings.ToLower(goquery.NodeName(selection.Children().First())) == "h1" {
				if v := trim(selection.Children().First().Text()); v != "" {
					titles = append(titles, v)
				}
			}
		case "div":
			if selection.HasClass("title") {
				if v := trim(selection.Text()); v != "" {
					titles = append(titles, v)
				}
			}
			if selection.HasClass("author") {
				if v := trim(selection.Text()); v != "" {
					authors = append(authors, v)
				}
			}
		case "span":
			if selection.HasClass("title") || selection.HasClass("title-span") {
				if v := trim(selection.Text()); v != "" {
					titles = append(titles, v)
				}
			}
			if selection.HasClass("author") || selection.HasClass("author-span") {
				if v := trim(selection.Text()); v != "" {
					authors = append(authors, v)
				}
			}
		case "a":
			if selection.AttrOr("rel", "") == "bookmark" {
				if v := trim(selection.Text()); v != "" {
					titles = append(titles, v)
				}
			}
			if selection.HasClass("author") {
				if v := trim(selection.Text()); v != "" {
					authors = append(authors, v)
				}
			}
		}
	})

	return &TitleAuthorDate{
		Title:  getAccurateText(titles),
		Author: getAccurateText(authors),
	}, nil
}
