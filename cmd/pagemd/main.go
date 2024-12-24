package main

import (
	"flag"
	"fmt"
	"github.com/LubyRuffy/pagemd/pkg/maindiv"
	"log"
)

func main() {
	url := flag.String("url", "", "url to transform")
	html := flag.String("html", "", "html to transform")
	depth := flag.Bool("depth", false, "whether to care about depth")
	flag.Parse()

	if *html == "" && *url == "" {
		log.Fatal("url and html is empty")
	}

	htmlString, err := maindiv.FetchPageHTML(*url, true)
	if err != nil {
		log.Fatal(err)
	}

	//htmlString, err := data.Load("b.html")
	//htmlString, err := data.Load("a.html")
	//if err != nil {
	//	log.Fatal(err)
	//}

	node, err := maindiv.ExtractMainContent(htmlString, *depth)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Content found at: %s\n", node)
	fmt.Printf("Content : %s\n", node.HTML)
}
