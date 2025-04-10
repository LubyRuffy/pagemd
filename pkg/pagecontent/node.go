package pagecontent

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const (
	MinContentText = 200
	MinDensity     = 10
)

var (
	ErrHasNotContent = errors.New("has not content")
)

// Node represents a div node with various attributes used for scoring and analysis.
type Node struct {
	score    float64 // score is the calculated score of the node based on text density, length, and depth.
	selector string  // selector is the CSS selector representing this node.

	TextLength int     // TextLength is the number of characters in the text content of the node.
	NodeCount  int     // NodeCount is the number of child nodes excluding common text nodes like <p>, <br>, etc.
	Depth      int     // Depth is the level of nesting of the node within the HTML structure.
	Density    float64 // Density is the ratio of text length to non-text child node count.
	Text       string  // Text is the trimmed text content of the node.
	HTML       string  // HTML is the original HTML code of the node.
	DirectText string  // DirectText is the text content of the node, excluding any child nodes.
}

// String returns a human-readable representation of the Node.
func (n *Node) String() string {
	return fmt.Sprintf("Node{TextLength: %d, NodeCount: %d, Depth: %d, Density: %.2f, Selector: %s, Score: %v, DirectText: %s}",
		n.TextLength, n.NodeCount, n.Depth, n.Density, n.selector, n.score, n.DirectText)
}

// CalculateScore computes the score of the node based on text density and length.
// If depthCare is true, the depth of the node is also considered in the scoring.
func (n *Node) CalculateScore(depthCare bool) float64 {
	// 基础分：密度分数
	densityScore := n.Density
	// 长度分数：文本长度的对数，避免过长文本占太大优势
	lengthScore := math.Log2(float64(n.TextLength))

	// 内容分析加分
	contentBonus := 0.0

	// 1. 考虑文本的链接密度
	// 有序列表或无序列表通常是页面的主要内容之一
	linkCount := strings.Count(n.HTML, "<a ")
	if linkCount > 3 {
		// 有一定数量的链接时，给予适当加分，但避免过度加分
		contentBonus += math.Min(float64(linkCount)*0.3, 15.0)
	}

	// 2. 检查结构特征
	// 重复模式通常是内容的一部分（如列表项、表格行等）
	repeatedPatterns := countRepeatedPatterns(n.HTML)
	if repeatedPatterns > 5 {
		contentBonus += math.Min(float64(repeatedPatterns)*0.5, 20.0)
	}

	// 3. 结构化元素加分
	// 文章的主要内容通常包含结构化元素
	if strings.Contains(n.HTML, "<table") ||
		strings.Contains(n.HTML, "<ul") ||
		strings.Contains(n.HTML, "<ol") {
		contentBonus += 10.0
	}

	// 4. 基于角色和语义的加分
	// 使用HTML5语义标签或ARIA角色的元素通常是主要内容
	if strings.Contains(n.selector, "[role='region']") ||
		strings.Contains(n.selector, "[role='main']") ||
		strings.Contains(n.selector, "[role='article']") ||
		strings.Contains(n.selector, "article") ||
		strings.Contains(n.selector, "main") ||
		strings.Contains(n.selector, "section") {
		contentBonus += 15.0
	}

	// 5. ID选择器检测 - 通常带有特定ID的元素是页面主要内容
	if strings.Contains(n.selector, "#content") ||
		strings.Contains(n.selector, "#main") ||
		strings.Contains(n.selector, "#article") ||
		strings.Contains(n.selector, "#report") {
		contentBonus += 25.0
	}

	// 6. 下载链接检测 - 包含下载链接的往往是重要内容
	if strings.Contains(n.HTML, ".pdf") ||
		strings.Contains(n.HTML, "download") ||
		strings.Contains(n.HTML, "Download") {
		contentBonus += 15.0
	}

	// 最终得分是密度分数、长度分数和内容分析加分的总和
	n.score = densityScore + lengthScore + contentBonus

	// 如果考虑深度，将得分乘以深度分数
	if depthCare {
		depthScore := math.Log(float64(n.Depth + 1))
		n.score = n.score * depthScore
	}

	return n.score
}

// countRepeatedPatterns counts the number of repeated HTML patterns in the content
// This helps identify lists, tables, and other structured content
func countRepeatedPatterns(html string) int {
	patterns := []string{
		"<li", "<tr", "<div class", "<p class", "<span class",
		"</li>", "</tr>", "href=", "class=\"", "id=\"",
	}

	count := 0
	for _, pattern := range patterns {
		count += strings.Count(html, pattern)
	}

	return count / len(patterns) // 返回平均值以平衡不同模式的影响
}

func valuableText(text string) string {
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\t", "")
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.ReplaceAll(text, " ", "")
	return text
}

// NewNodeFromSelection creates a new Node from a goquery.Selection.
// It extracts text, HTML content, and calculates various attributes like Depth, TextLength, NodeCount, and Density.
func NewNodeFromSelection(s *goquery.Selection) *Node {
	text := valuableText(strings.TrimSpace(s.Text()))
	htmlContent, err := s.Html()
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	n := s.Nodes[0].FirstChild
	for n != nil {
		if n.Type == html.TextNode {
			if v := strings.Trim(n.Data, "\n\t\r "); v != "" {
				buf.WriteString(n.Data)
			}
		}
		n = n.NextSibling
	}

	parents := s.ParentsFiltered("*").Nodes

	node := &Node{
		Depth:      len(parents),   // Depth is determined by the number of parent nodes.
		Text:       text,           // Text is the trimmed content of the node.
		HTML:       htmlContent,    // HTML stores the original HTML code of the node.
		selector:   getSelector(s), // Selector represents the CSS path to this node.
		TextLength: len(text),
		DirectText: buf.String(),
	}

	// Calculate NodeCount by subtracting common text nodes from total child nodes.
	node.NodeCount = s.Find("*").Length() -
		s.Find("p").Length() -
		s.Find("br").Length() -
		s.Find("code").Length() -
		s.Find("span").Length() -
		s.Find("pre").Length() -
		s.Find("article").Length() -
		s.Find("hr").Length() -
		s.Find("h1").Length() -
		s.Find("h2").Length() -
		s.Find("h3").Length() -
		s.Find("h4").Length() -
		s.Find("table").Length() -
		s.Find("tr").Length() -
		s.Find("td").Length() -
		s.Find("th").Length() -
		s.Find("ul").Length() -
		s.Find("li").Length() -
		s.Find("section").Length() +
		5 // 保证有值

	// Density is the ratio of text length to non-text child node count.
	node.Density = float64(node.TextLength-len(s.Find("style").Text())) / float64(node.NodeCount)

	return node
}

// getSelector 获取元素的选择器路径
func getSelector(s *goquery.Selection) string {
	var path []string
	cur := s

	for {
		if cur.Length() == 0 {
			break
		}

		// 获取当前元素的选择器
		tag := goquery.NodeName(cur)
		if id, exists := cur.Attr("id"); exists {
			path = append([]string{"#" + id}, path...)
			break
		}
		if class, exists := cur.Attr("class"); exists {
			classes := strings.Fields(class)
			// 处理冒号
			var okClasses []string
			for _, c := range classes {
				if strings.Contains(c, ":") {
					c = strings.ReplaceAll(c, ":", "\\:")
				}
				okClasses = append(okClasses, c)
			}

			if len(okClasses) > 0 {
				path = append([]string{tag + "." + strings.Join(okClasses, ".")}, path...)
			} else {
				path = append([]string{tag}, path...)
			}
		} else {
			path = append([]string{tag}, path...)
		}

		cur = cur.Parent()
	}

	return strings.Join(path, " > ")
}

// extractMainContent 从 HTML 内容中提取包含主要内容的节点
//
// 提取出 HTML 内容中的主要内容，并返回主要内容的节点和错误信息。
func extractMainContent(htmlContent string, depthCare bool, debug bool) (*Node, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htmlContent)))
	if err != nil {
		return nil, err
	}

	var bestNode *Node
	var maxScore float64

	extract := func(i int, s *goquery.Selection) {
		// 跳过文本内容太少的节点
		if len(s.Text()) < MinContentText {
			return
		}

		node := NewNodeFromSelection(s)
		// 跳过密度太低的节点
		if node.Density < MinDensity {
			return
		}

		score := node.CalculateScore(depthCare)

		if debug {
			log.Println(score, node)
		}

		if score > maxScore {
			maxScore = score
			bestNode = node

			if debug {
				log.Println("change to:", score, node)
			}
		}
	}

	// 1. 首先尝试通过ID选择器查找明确的内容区域
	idSelectors := []string{
		"#content_1_reportImage", "#content", "#main-content", "#article-content",
		"#report", "#main", "#primary", "#content-body", "#post-content",
	}

	for _, selector := range idSelectors {
		doc.Find(selector).Each(extract)
		if bestNode != nil {
			break
		}
	}

	// 2. 如果没找到，再尝试通过语义标签查找内容
	if bestNode == nil {
		semanticSelectors := []string{
			"article", "main", "section.content", "div.content", "div.main",
			"div.article", "div[role='main']", "div[role='article']",
			"div[role='region'][aria-label*='内容']", "div[role='region'][aria-label*='content']",
			"div.post-content", "div.entry-content",
		}

		for _, selector := range semanticSelectors {
			doc.Find(selector).Each(extract)
			if bestNode != nil {
				break
			}
		}
	}

	// 3. 如果仍没找到，查找具有特定特征的内容区域
	if bestNode == nil {
		// 先查找包含下载链接的区域
		doc.Find("div:has(a[href*='.pdf']), div:has(a:contains('Download'))").Each(extract)

		// 查找包含表格、表单等结构化元素的区域
		if bestNode == nil {
			doc.Find("div:has(table), div:has(form), div:has(ul), div:has(ol)").Each(extract)

			// 如果还没找到，查找所有可能的内容容器
			if bestNode == nil {
				doc.Find("div").Each(extract)
				doc.Find("td").Each(extract)
				doc.Find("section").Each(extract)
			}
		}
	}

	if bestNode != nil {
		return bestNode, nil
	}

	return nil, ErrHasNotContent
}
