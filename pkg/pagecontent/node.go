package pagecontent

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"log"
	"math"
	"strings"
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
	densityScore := n.Density
	lengthScore := math.Log2(float64(n.TextLength))
	n.score = densityScore + lengthScore
	if depthCare {
		depthScore := math.Log(float64(n.Depth + 1))
		n.score = n.score * depthScore
	}
	return n.score
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
		if len(s.Text()) < MinContentText {
			return
		}

		node := NewNodeFromSelection(s)
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

	doc.Find("div").Each(extract)
	doc.Find("article").Each(extract)
	//doc.Find("p").Each(extract)
	//doc.Find("span").Each(extract)

	if bestNode != nil {
		return bestNode, nil
	}

	return nil, ErrHasNotContent
}
