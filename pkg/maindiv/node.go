package maindiv

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math"
	"strings"
)

// Node div节点
type Node struct {
	score    float64 // 分数
	selector string  // css选择器，对应浏览器console中的 $('selector') 语法

	TextLength int     // 文本长度
	NodeCount  int     // 标签个数
	Depth      int     // 节点深度
	Density    float64 // 文本密度
	Text       string  // 文本
	HTML       string  // 原始html代码
}

func (n *Node) String() string {
	return fmt.Sprintf("Node{TextLength: %d, NodeCount: %d, Depth: %d, Density: %.2f, Selector: %s, Score: %v}",
		n.TextLength, n.NodeCount, n.Depth, n.Density, n.selector, n.score)
}

// CalculateScore 计算分数
func (n *Node) CalculateScore(depthCare bool) float64 {
	densityScore := n.Density
	lengthScore := math.Log(float64(n.TextLength)) / 5
	n.score = densityScore * lengthScore
	if depthCare {
		depthScore := math.Log(float64(n.Depth + 1))
		n.score = n.score * depthScore
	}
	return n.score
}

// NewNodeFromSelection 从goquery.Selection创建Node
func NewNodeFromSelection(s *goquery.Selection) *Node {
	text := strings.TrimSpace(s.Text())
	html, err := s.Html()
	if err != nil {
		panic(err)
	}

	node := &Node{
		Depth:      len(s.ParentsFiltered("*").Nodes), // 深度
		Text:       text,                              // 文本
		HTML:       html,                              // 文本
		selector:   GetSelector(s),                    // 计算选择器路径
		TextLength: len(text),
	}

	// 计算文本和节点个数
	node.NodeCount = s.Find("*").Length() -
		// 把常见的文本节点删除
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
		//s.Find("header").Length() -
		s.Find("section").Length()

	// 计算文本密度
	if node.NodeCount > 0 {
		node.Density = float64(node.TextLength) / float64(node.NodeCount)
	} else {
		node.Density = 1
	}

	return node
}
