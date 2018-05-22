package main

import (
	"cises/sitemap"
	"fmt"
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
