package maindiv

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
)

type Analysis struct {
}

// ExtractMainContent 从 HTML 内容中提取包含主要内容的节点
//
// 提取出 HTML 内容中的主要内容，并返回主要内容的节点和错误信息。
func ExtractMainContent(htmlContent string, depthCare bool) (*Node, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htmlContent)))
	if err != nil {
		return nil, err
	}

	var bestNode *Node
	var maxScore float64

	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		if len(s.Text()) < 100 {
			return
		}

		node := NewNodeFromSelection(s)
		score := node.CalculateScore(depthCare)

		if score > maxScore {
			maxScore = score
			bestNode = node
			//log.Println(score, node)
		}
	})

	if bestNode != nil {
		return bestNode, nil
	}
	return nil, nil
}
