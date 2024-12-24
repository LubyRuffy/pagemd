package pagecontent

import (
	"github.com/stretchr/testify/assert"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// TestNodeString tests the String method of Node.
func TestNodeString(t *testing.T) {
	node := &Node{
		score:      10.5,
		selector:   "#test",
		TextLength: 100,
		NodeCount:  10,
		Depth:      3,
		Density:    2.5,
		Text:       "Test text",
		HTML:       "<div id=\"test\">Test text</div>",
	}

	expected := "Node{TextLength: 100, NodeCount: 10, Depth: 3, Density: 2.50, Selector: #test, Score: 10.5}"
	if node.String() != expected {
		t.Errorf("Expected %q, got %q", expected, node.String())
	}
}

// TestNodeCalculateScore tests the CalculateScore method of Node.
func TestNodeCalculateScore(t *testing.T) {
	node := &Node{
		TextLength: 100,
		NodeCount:  10,
		Depth:      3,
		Density:    2.5,
	}

	score := node.CalculateScore(false)
	expected := 2.5 * (math.Log(float64(100)) / 5)
	if math.Abs(score-expected) > 0.0001 {
		t.Errorf("Expected %v, got %v", expected, score)
	}

	score = node.CalculateScore(true)
	expected *= math.Log(float64(3 + 1))
	if math.Abs(score-expected) > 0.0001 {
		t.Errorf("Expected %v, got %v", expected, score)
	}
}

// TestNewNodeFromSelection tests the NewNodeFromSelection function.
func TestNewNodeFromSelection(t *testing.T) {
	html := `<html><body><div id="test">
		<p>This is a test paragraph.</p>
		<span>Span text</span>
		<br/>
		<code>Code snippet</code>
		<pre>Preformatted text</pre>
		<article>Article content</article>
		<hr/>
		<h1>Main heading</h1>
		<h2>Subheading 1</h2>
		<h3>Subheading 2</h3>
		<h4>Subheading 3</h4>
		<section>Section content</section>
		Some text here.
		<img src="test.jpg" alt="Test image"></img>
	</div></body></html>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}

	sel := doc.Find("#test")
	node := NewNodeFromSelection(sel)

	expectedTextLength := 192
	expectedDepth := 2 // html body
	expectedNodeCount := sel.Find("*").Length() -
		sel.Find("p").Length() -
		sel.Find("br").Length() -
		sel.Find("code").Length() -
		sel.Find("span").Length() -
		sel.Find("pre").Length() -
		sel.Find("article").Length() -
		sel.Find("hr").Length() -
		sel.Find("h1").Length() -
		sel.Find("h2").Length() -
		sel.Find("h3").Length() -
		sel.Find("h4").Length() -
		sel.Find("section").Length()

	expectedDensity := float64(expectedTextLength) / float64(expectedNodeCount)

	if node.TextLength != expectedTextLength {
		t.Errorf("Expected TextLength %d, got %d", expectedTextLength, node.TextLength)
	}
	if node.Depth != expectedDepth {
		t.Errorf("Expected Depth %d, got %d", expectedDepth, node.Depth)
	}
	if node.NodeCount != expectedNodeCount {
		t.Errorf("Expected NodeCount %d, got %d", expectedNodeCount, node.NodeCount)
	}
	if math.Abs(node.Density-expectedDensity) > 0.0001 {
		t.Errorf("Expected Density %.2f, got %.2f", expectedDensity, node.Density)
	}

	expectedHTML := "\n\t\t<p>This is a test paragraph.</p>\n\t\t<span>Span text</span>\n\t\t<br/>\n\t\t<code>Code snippet</code>\n\t\t<pre>Preformatted text</pre>\n\t\t<article>Article content</article>\n\t\t<hr/>\n\t\t<h1>Main heading</h1>\n\t\t<h2>Subheading 1</h2>\n\t\t<h3>Subheading 2</h3>\n\t\t<h4>Subheading 3</h4>\n\t\t<section>Section content</section>\n\t\tSome text here.\n\t\t<img src=\"test.jpg\" alt=\"Test image\"/>\n\t"
	if node.HTML != expectedHTML {
		t.Errorf("Expected HTML %q, got %q", expectedHTML, node.HTML)
	}

	expectedText := "This is a test paragraph.\n\t\tSpan text\n\t\t\n\t\tCode snippet\n\t\tPreformatted text\n\t\tArticle content\n\t\t\n\t\tMain heading\n\t\tSubheading 1\n\t\tSubheading 2\n\t\tSubheading 3\n\t\tSection content\n\t\tSome text here."
	if node.Text != expectedText {
		t.Errorf("Expected Text %q, got %q", expectedText, node.Text)
	}
}

func Test_extractMainContent(t *testing.T) {
	testF := func(f string, s string) {
		h, err := os.ReadFile(filepath.Join("testdata", f))
		assert.NoError(t, err)
		node, err := extractMainContent(string(h), false)
		assert.NoError(t, err)
		assert.Equal(t, node.selector, s)
	}

	//testF("a.html", "#content > article > div.entry-content")
	//testF("b.html", "#tinymce-editor > div")
	//testF("c.html", "html > body > main > section.section.pt-7 > div.container > div.row.justify-center > article.lg\\:col-7 > div.content.mb-10")
	testF("d.html", "html.h-full.overflow-y-scroll > body.antialiased.min-h-screen.w-full.m-0.flex.flex-col > main.flex-grow > article.mx-auto.flex.flex-1.max-w-2xl.w-full.flex-col.space-y-3.px-6.py-16.md\\:px-0")
}

func TestGetSelector(t *testing.T) {
	tests := []struct {
		html     string
		selector string
		expected string
	}{
		{
			html:     `<div id="main"><p class="text">Hello</p></div>`,
			selector: "p",
			expected: "#main > p.text", // 优先ID
		},
		{
			html:     `<div id="main"><p id="unique">Hello</p></div>`,
			selector: "p",
			expected: "#unique",
		},
		{
			html:     `<div><p class="text1 text2">Hello</p></div>`,
			selector: "p",
			expected: "html > body > div > p.text1.text2", // 会补充html和body
		},
		{
			html:     `<div><p>Hello</p></div>`,
			selector: "p",
			expected: "html > body > div > p", // 会补充html和body
		},
		{
			html:     `<div class="row justify-center"><article class=lg:col-7><div1>aaa</div1></article></div>`,
			selector: "div1",
			expected: "html > body > div.row.justify-center > article.lg\\:col-7 > div1", // 会补充html和body
		},
	}

	for _, test := range tests {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(test.html))
		if err != nil {
			t.Fatalf("Failed to parse HTML: %v", err)
		}

		sel := doc.Find(test.selector)
		result := getSelector(sel)
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}
