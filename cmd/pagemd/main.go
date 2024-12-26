package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/LubyRuffy/pagemd/pkg/pagecontent"
	"log"
	"os"
)

func main() {
	output := flag.String("output", "", "output of markdown")

	ci, err := pagecontent.NewFromFlags(
		pagecontent.WithOnMainNodeFound(func(node *pagecontent.Node) {
			log.Println("found main node, size:", len(node.HTML))
		}),
		pagecontent.WithOnHtmlFetched(func(htmlContent string) {
			log.Println("fetched html, size:", len(htmlContent))
		}),
	).ExtractMainContent()
	if err != nil {
		log.Fatal(err)
	}

	d, _ := json.Marshal(ci)
	log.Println(string(d))

	// fmt.Printf("Content : %s\n%s", contentHtml, markdown)
	if *output == "" {
		fmt.Println(ci.TitleAuthorDate, ci.Markdown)
	} else {
		if err = os.WriteFile("out.md", []byte(ci.Markdown), 0644); err != nil {
			log.Fatal(err)
		}
	}
}
