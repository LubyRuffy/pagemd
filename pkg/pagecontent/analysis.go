package pagecontent

import (
	"errors"
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

type Analysis struct {
	cfg *Config
}

func (a *Analysis) ExtractMainContent() (string, string, error) {
	htmlContent := a.cfg.html

	if htmlContent == "" {
		if a.cfg.url == "" {
			return "", "", ErrNoURLProvided
		}

		var err error
		htmlContent, err = fetchPageHTML(a.cfg.url, a.cfg.headless, a.cfg.debug)
		if err != nil {
			return "", "", err
		}
	}

	node, err := extractMainContent(htmlContent, a.cfg.depthCare)
	if err != nil {
		return "", "", err
	}

	markdown, err := htmltomarkdown.ConvertString(node.HTML,
		converter.WithDomain(a.cfg.url))
	if err != nil {
		log.Fatal(err)
	}

	return node.HTML, markdown, err
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
