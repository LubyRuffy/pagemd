package maindiv

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

// GetSelector 获取元素的选择器路径
func GetSelector(s *goquery.Selection) string {
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
			if len(classes) > 0 {
				path = append([]string{tag + "." + strings.Join(classes, ".")}, path...)
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
