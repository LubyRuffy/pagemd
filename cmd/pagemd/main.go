package main

import (
	"flag"
	"fmt"
	"github.com/LubyRuffy/pagemd/pkg/pagecontent"
	"log"
	"os"
)

func main() {
	url := flag.String("url", "", "url to transform")
	html := flag.String("html", "", "html to transform")
	output := flag.String("output", "", "output of markdown")
	depth := flag.Bool("depth", false, "whether to care about depth")
	headless := flag.Bool("headless", false, "whether headless when url fetch")
	debug := flag.Bool("debug", false, "whether debug")
	flag.Parse()

	if *html == "" && *url == "" {
		log.Fatal("url and html is empty")
	}

	_, markdown, err := pagecontent.NewAnalysis(
		pagecontent.WithDepthCare(*depth),
		pagecontent.WithHeadless(*headless),
		pagecontent.WithDebug(*debug),
		pagecontent.WithURL(*url),
		pagecontent.WithHTML(*html),
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

	// fmt.Printf("Content : %s\n%s", contentHtml, markdown)
	if *output == "" {
		fmt.Println(markdown)
	} else {
		if err = os.WriteFile("out.md", []byte(markdown), 0644); err != nil {
			log.Fatal(err)
		}
	}
}
