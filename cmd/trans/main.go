package main

import (
	"context"
	"fmt"
	"github.com/LubyRuffy/pagemd/pkg/aitrans"
	"github.com/LubyRuffy/pagemd/pkg/pagecontent"
	"log"
)

func main() {
	_, md, err := pagecontent.NewFromFlags(
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

	log.Println("try to translate...")

	a := aitrans.New()
	a.TranslateToChinese(context.Background(),
		md,
		func(s string) {
			fmt.Printf("%s", s)
		})
}
