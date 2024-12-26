package pagecontent

import (
	"errors"
	"flag"
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"log"
)

var ErrNoURLProvided = errors.New("URL must be provided")

type Config struct {
	depthCare bool   // 算分是否关心深度
	headless  bool   // 是否headless模式请求网页
	debug     bool   // 是否调试模式，比如headless下显示浏览器
	html      string // HTML内容
	url       string // 网页URL

	onMainNodeFound func(node *Node)         // 提取了正文节点的回调
	onHtmlFetched   func(htmlContent string) // 请求url获取了html时的回调
}

type ConfigOpt func(*Config)

func WithDepthCare(depthCare bool) ConfigOpt {
	return func(cfg *Config) {
		cfg.depthCare = depthCare
	}
}

func WithHeadless(headless bool) ConfigOpt {
	return func(cfg *Config) {
		cfg.headless = headless
	}
}

func WithDebug(debug bool) ConfigOpt {
	return func(cfg *Config) {
		cfg.debug = debug
	}
}

func WithHTML(html string) ConfigOpt {
	return func(cfg *Config) {
		cfg.html = html
	}
}
func WithURL(url string) ConfigOpt {
	return func(cfg *Config) {
		cfg.url = url
	}
}
func WithOnHtmlFetched(onHtmlFetched func(htmlContent string)) ConfigOpt {
	return func(cfg *Config) {
		cfg.onHtmlFetched = onHtmlFetched
	}
}
func WithOnMainNodeFound(onMainNodeFound func(node *Node)) ConfigOpt {
	return func(cfg *Config) {
		cfg.onMainNodeFound = onMainNodeFound
	}
}

type Analysis struct {
	cfg *Config
}

type ContentInfo struct {
	*TitleAuthorDate
	URL             string `json:"url"`
	HTML            string `json:"html"`
	ContentHTML     string `json:"content_html"`
	ContentMarkdown string `json:"content_markdown"`
}

// ExtractMainContent 提取一个网页的正文内容，去除不相关的信息
func (a *Analysis) ExtractMainContent() (*ContentInfo, error) {
	htmlContent := a.cfg.html

	if htmlContent == "" {
		if a.cfg.url == "" {
			return nil, ErrNoURLProvided
		}

		var err error
		htmlContent, err = fetchPageHTML(a.cfg.url, a.cfg.headless, a.cfg.debug)
		if err != nil {
			return nil, err
		}
		a.onHtmlFetched(htmlContent)
	}

	tad, err := ExtractTitleAuthorDate(htmlContent)
	if err != nil {
		return nil, err
	}

	node, err := extractMainContent(htmlContent, a.cfg.depthCare, a.cfg.debug)
	if err != nil {
		return nil, err
	}
	a.onMainNodeFound(node)

	markdown, err := htmltomarkdown.ConvertString(node.HTML,
		converter.WithDomain(a.cfg.url))
	if err != nil {
		log.Fatal(err)
	}

	return &ContentInfo{
		TitleAuthorDate: tad,
		URL:             a.cfg.url,
		HTML:            htmlContent,
		ContentHTML:     node.HTML,
		ContentMarkdown: markdown,
	}, nil
}

func (a *Analysis) onMainNodeFound(node *Node) {
	if a.cfg.onMainNodeFound != nil {
		a.cfg.onMainNodeFound(node)
	}
}

func (a *Analysis) onHtmlFetched(htmlContent string) {
	if a.cfg.onHtmlFetched != nil {
		a.cfg.onHtmlFetched(htmlContent)
	}
}

func NewAnalysis(opts ...ConfigOpt) *Analysis {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return &Analysis{
		cfg: cfg,
	}
}

func NewFromFlags(opts ...ConfigOpt) *Analysis {
	url := flag.String("url", "", "url to transform")
	html := flag.String("html", "", "html to transform")
	depth := flag.Bool("depth", false, "whether to care about depth")
	headless := flag.Bool("headless", false, "whether headless when url fetch")
	debug := flag.Bool("debug", false, "whether debug")
	flag.Parse()

	if *html == "" && *url == "" {
		log.Fatal("url and html is empty")
	}

	opts = append(opts, WithDepthCare(*depth),
		WithHeadless(*headless),
		WithDebug(*debug),
		WithURL(*url),
		WithHTML(*html))

	return NewAnalysis(opts...)
}
