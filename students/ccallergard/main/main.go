package main

import (
	"fmt"
	"github.com/gophercises/sitemap/students/ccallergard"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Provide http(s) URL as argument")
		return
	}

	site, err := sitemap.Map(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sitemap.XMLSiteMap(site))
}
