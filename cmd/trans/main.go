package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LubyRuffy/pagemd/pkg/aitrans"
	"github.com/LubyRuffy/pagemd/pkg/pagecontent"
	"log"
)

func main() {
	ci, err := pagecontent.NewFromFlags(
		pagecontent.WithOnMainNodeFound(func(node *pagecontent.Node) {
			log.Println("found main node, size:", len(node.HTML))
		}),
		pagecontent.WithOnHtmlFetched(func(htmlContent string) {
			log.Println("fetched html, size:", len(htmlContent))
		}),
	).ExtractMainContent()
	if err != nil {
		panic(err)
	}

	d, _ := json.Marshal(ci)
	log.Println(string(d))

	log.Println("try to translate...")

	a := aitrans.New()
	a.TranslateToChinese(context.Background(),
		ci.Markdown,
		func(s string) {
			fmt.Printf("%s", s)
		})
}
