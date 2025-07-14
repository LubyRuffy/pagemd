package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LubyRuffy/pagemd/pkg/llm"
	"log"

	"github.com/LubyRuffy/pagemd/pkg/aitrans"
	"github.com/LubyRuffy/pagemd/pkg/pagecontent"
)

func main() {
	cfg, err := llm.Load("config.yaml")
	if err != nil {
		panic(err)
	}

	ci, err := pagecontent.NewFromFlags(
		pagecontent.WithOnMainContentFound(func(s string) {
			log.Println("found main content, size:", len(s))
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

	a := aitrans.New(cfg)
	a.TranslateToChinese(context.Background(),
		ci.Markdown,
		func(s string) {
			fmt.Printf("%s", s)
		})
}
