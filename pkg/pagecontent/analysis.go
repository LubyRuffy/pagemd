package pagecontent

import (
	"errors"
	"flag"
	"log"
)

var ErrNoURLProvided = errors.New("URL must be provided")

type Config struct {
	depthCare bool   // 算分是否关心深度
	debug     bool   // 是否调试模式，比如headless下显示浏览器
	url       string // 网页URL

	onMainContentFound func(htmlContent string) // 提取了正文节点的回调
	onHtmlFetched      func(htmlContent string) // 请求url获取了html时的回调
}

type ConfigOpt func(*Config)

func WithDepthCare(depthCare bool) ConfigOpt {
	return func(cfg *Config) {
		cfg.depthCare = depthCare
	}
}

func WithDebug(debug bool) ConfigOpt {
	return func(cfg *Config) {
		cfg.debug = debug
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
func WithOnMainContentFound(onMainNodeFound func(string)) ConfigOpt {
	return func(cfg *Config) {
		cfg.onMainContentFound = onMainNodeFound
	}
}

type Analysis struct {
	cfg *Config
}

type ContentInfo struct {
	*TitleAuthorDate
	URL         string `json:"url"`
	RawHTML     string `json:"raw_html"`
	HTML        string `json:"html"`
	ContentHTML string `json:"content_html"`
	Markdown    string `json:"markdown"`
	Text        string `json:"text"`
}

// ExtractMainContent 提取一个网页的正文内容，去除不相关的信息
func (a *Analysis) ExtractMainContent() (*ContentInfo, error) {
	if a.cfg.url == "" {
		return nil, ErrNoURLProvided
	}

	result, err := fetchPage(a.cfg.url, a.cfg.debug)
	if err != nil {
		return nil, err
	}
	a.onHtmlFetched(result.HTML)
	a.onMainContentFound(result.Content)

	tad, err := ExtractTitleAuthorDate(result.HTML)
	if err != nil {
		return nil, err
	}

	return &ContentInfo{
		URL:             a.cfg.url,
		TitleAuthorDate: tad,
		RawHTML:         result.RawHTML,
		HTML:            result.HTML,
		ContentHTML:     result.Content,
		Markdown:        result.Markdown,
		Text:            result.TextContent,
	}, nil
}

func (a *Analysis) onMainContentFound(s string) {
	if a.cfg.onMainContentFound != nil {
		a.cfg.onMainContentFound(s)
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
	depth := flag.Bool("depth", false, "whether to care about depth")
	debug := flag.Bool("debug", false, "whether debug")
	flag.Parse()

	if *url == "" {
		log.Fatal("url and html is empty")
	}

	opts = append(opts, WithDepthCare(*depth),
		WithDebug(*debug),
		WithURL(*url))

	return NewAnalysis(opts...)
}
