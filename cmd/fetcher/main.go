package main

import (
	"fmt"
	"log"

	"www.github.com/maxbrt/colibri/internal/rss"
)

func main() {
	sources, err := rss.ReadSources("./sources")
	if err != nil {
		log.Fatalf("%s", err)
	}
	for _, s := range sources {
		posts, err := rss.FetchAndParse(s)
		if err != nil {
			log.Fatalf("%s", err)
		}
		for _, p := range posts {
			fmt.Printf("%s, %s\n", p.Title, p.Description)
		}
	}
}
